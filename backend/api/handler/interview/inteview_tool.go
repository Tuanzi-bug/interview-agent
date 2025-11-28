package interview

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/service"
	"ai-eino-interview-agent/internal/alert"
	"ai-eino-interview-agent/internal/mq"
	interviewservice "ai-eino-interview-agent/internal/service/interviews"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"golang.org/x/net/context"
)

//面试工具类 复用到这里

// runInterviewLoopAsync 异步运行面试循环
func runInterviewLoopAsyncTool(ctx context.Context, userId uint, writer io.Writer, session *InterviewSession, interviewService interviewservice.InterviewManager) {
	defer func() {
		if r := recover(); r != nil {
			sendErrorEvent(writer, fmt.Sprintf("面试异常: %v", r))
		}
		// 延迟删除会话，给前端充足时间来获取最后的数据
		go func() {
			time.Sleep(10 * time.Second)
			GetSessionManager().DeleteSession(session.SessionID)
		}()
	}()

	sm := GetSessionManager()
	var resumeContent string
	const answerTimeout = 30 * time.Minute
	const heartbeatInterval = 15 * time.Second

	// 根据面试类型选择维度
	var dimensions []string
	// 专项面试
	//dimensions = []string{
	//	"basic_knowledge_mastery",
	//	"working_principle_practical_experience",
	//	"advanced_features_application",
	//	"problem_troubleshooting_skills",
	//	"architecture_design_thinking",
	//	"performance_optimization_ability",
	//}
	dimensionIndex := 0
	questionNumber := 1
	awaitingAnswer := false
	lastQuestionIndex := 0

	for dimensionIndex < len(dimensions) || awaitingAnswer {
		select {
		case <-ctx.Done():
			log.Printf("[Interview Loop] Context cancelled, sessionID: %s", session.SessionID)
			return
		default:
		}

		if awaitingAnswer {
			log.Printf("[Interview Loop] Waiting for answer, sessionID: %s, questionIndex: %d", session.SessionID, lastQuestionIndex)
			answer, received := waitForAnswerWithHeartbeat(sm, session.SessionID, answerTimeout, heartbeatInterval, writer)
			log.Printf("[Interview Loop] Answer received: %v, sessionID: %s", received, session.SessionID)
			if !received {
				log.Printf("[Interview Loop] Answer timeout, sessionID: %s", session.SessionID)
				sendErrorEvent(writer, "等待答案超时，面试已结束")
				sendCompleteEvent(writer)
				break
			}

			session.AllDialogues = append(session.AllDialogues, map[string]interface{}{
				"speaker_type":  "candidate",
				"content":       answer,
				"display_order": uint32(lastQuestionIndex) * 100,
			})

			awaitingAnswer = false
			if dimensionIndex >= len(dimensions) {
				sendTopicCompleteEvent(writer)
				break
			}
			continue
		}

		var prompt string
		currentQuestionIndex := questionNumber
		if dimensionIndex >= len(dimensions) {
			break
		}

		if currentQuestionIndex == 1 {
			prompt = buildPrompt(currentQuestionIndex, session.Query, session.ResumeId, session.HasResume, dimensions[dimensionIndex], 0, session.Type, session.Domain, session.Difficulty)
		} else {
			userAnswers := ""
			for i := 1; i < currentQuestionIndex; i++ {
				displayOrder := uint32(i) * 100
				for _, dialogue := range session.AllDialogues {
					d := dialogue.(map[string]interface{})
					if d["speaker_type"] == "candidate" && d["display_order"] == displayOrder {
						userAnswers += fmt.Sprintf("问题 %d 的回答：%s\n", i, d["content"])
						break
					}
				}
			}
			prompt = buildPrompt(currentQuestionIndex, resumeContent+"\n\n用户已回答的问题：\n"+userAnswers, 0, false, dimensions[dimensionIndex], 0, session.Type, session.Domain, session.Difficulty)
		}

		// 只有第一个问题才需要简历解析工具
		isFirstQuestion := currentQuestionIndex == 1
		result, err := service.GenerateInterviewQuestions(ctx, prompt, userId, isFirstQuestion, session.Type)
		if err != nil {
			sendErrorEvent(writer, "Failed to generate question: "+err.Error())
			sendCompleteEvent(writer)
			break
		}

		//if len(result.Questions) == 0 {
		//	sendTopicCompleteEvent(writer)
		//	sendCompleteEvent(writer)
		//	break
		//}

		if currentQuestionIndex == 1 {
			resumeContent = session.Query
		}

		q := result.Questions[0]
		sendQuestionEvent(writer, currentQuestionIndex, q)

		session.AllQuestions = append(session.AllQuestions, map[string]interface{}{
			"question_text":  q.QuestionText,
			"eval_dimension": q.EvalDimension,
			"order":          q.Order,
		})

		if len(result.Dialogues) > 0 {
			for _, d := range result.Dialogues {
				if d.SpeakerType == "interviewer" {
					session.AllDialogues = append(session.AllDialogues, map[string]interface{}{
						"speaker_type":  "interviewer",
						"content":       d.Content,
						"display_order": uint32(currentQuestionIndex) * 100,
					})
					break
				}
			}
		}

		sendReadyEventWithSession(writer, currentQuestionIndex, session.SessionID)
		sm.ClearAnswer(session.SessionID)
		awaitingAnswer = true
		lastQuestionIndex = currentQuestionIndex
		questionNumber++
		dimensionIndex++
	}

	log.Printf("[Interview Loop] Saving dialogues, sessionID: %s", session.SessionID)
	// 保存面试对话，带重试机制
	const maxRetries = 3
	var saveErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		saveErr = interviewService.SaveInterviewDialogues(ctx, session.UserID, session.RecordID, session.AllQuestions, session.AllDialogues)
		if saveErr == nil {
			break
		}

		// 判断是否为可重试的错误
		if !IsRetryableError(saveErr) {
			// 上下文取消/超时通常是用户主动中断或请求生命周期结束，不发送告警
			if !errors.Is(saveErr, context.Canceled) && !errors.Is(saveErr, context.DeadlineExceeded) {
				alert.SendDatabaseErrorAlert(
					fmt.Sprintf("SaveInterviewDialogues (不可重试) - UserID: %d, RecordID: %d", session.UserID, session.RecordID),
					saveErr,
					attempt+1,
				)
			}
			break
		}

		// 最后一次尝试失败（所有重试机会耗尽）
		if attempt == maxRetries-1 {
			alert.SendDatabaseErrorAlert(
				fmt.Sprintf("SaveInterviewDialogues (重试耗尽) - UserID: %d, RecordID: %d", session.UserID, session.RecordID),
				saveErr,
				maxRetries,
			)
			break
		}

		// 指数退避：等待 100ms * 2^attempt
		backoffDuration := time.Duration(100*(1<<uint(attempt))) * time.Millisecond
		time.Sleep(backoffDuration)
	}

	duration := int64(time.Since(session.StartTime).Seconds())

	updateDTO := &interviewsapi.InterviewRecordDTO{
		ID:       int64(session.RecordID),
		UserID:   int32(session.UserID),
		Status:   "completed",
		Duration: &duration,
	}

	log.Printf("[Interview Loop] Updating interview record, sessionID: %s, duration: %d seconds", session.SessionID, duration)
	if err := interviewService.UpdateInterviewRecord(ctx, updateDTO); err != nil {
		// 记录更新失败，但不中断流程
		log.Printf("[Interview Loop] Failed to update interview record: %v, sessionID: %s", err, session.SessionID)
		_ = err
	}

	// 面试完成后，发送 MQ 消息触发评估报告生成
	log.Printf("[Interview Loop] Publishing evaluation messages, sessionID: %s, userID: %d, recordID: %d", session.SessionID, session.UserID, session.RecordID)

	// 发布评估报告生成消息
	if err := mq.PublishEvaluationReport(ctx, session.UserID, session.RecordID); err != nil {
		log.Printf("[Interview Loop] Failed to publish evaluation report message: %v, sessionID: %s", err, session.SessionID)
	}

	// 发布主题评估消息
	if err := mq.PublishTopicEvaluation(ctx, session.UserID, session.RecordID); err != nil {
		log.Printf("[Interview Loop] Failed to publish topic evaluation message: %v, sessionID: %s", err, session.SessionID)
	}

	log.Printf("[Interview Loop] Interview completed, sessionID: %s", session.SessionID)
}
