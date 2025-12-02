package mianshi

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"ai-eino-interview-agent/chatApp/agent/service"
	"ai-eino-interview-agent/internal/alert"
	"ai-eino-interview-agent/internal/mq"
	interviewservice "ai-eino-interview-agent/internal/service/interviews"
)

// InterviewEngine 面试引擎 - 处理核心面试逻辑
type InterviewEngine struct {
	sessionManager *SessionManager
	interviewSvc   interviewservice.InterviewManager
	writer         io.Writer
}

// NewInterviewEngine 创建面试引擎
func NewInterviewEngine(sessionManager *SessionManager, interviewSvc interviewservice.InterviewManager, writer io.Writer) *InterviewEngine {
	return &InterviewEngine{
		sessionManager: sessionManager,
		interviewSvc:   interviewSvc,
		writer:         writer,
	}
}

// RunInterviewLoop 运行面试循环
func (e *InterviewEngine) RunInterviewLoop(ctx context.Context, session *InterviewSession) {
	defer e.cleanup(session)

	questionIndex := 0
	var resumeContent string
	const maxFollowUps = 2
	const answerTimeout = 30 * time.Minute
	const heartbeatInterval = 15 * time.Second

	// 根据面试类型选择维度
	dimensions := e.getDimensions(session.Type)
	dimensionIndex := 0
	followUpCount := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Interview Engine] Context cancelled, sessionID: %s", session.SessionID)
			return
		default:
		}

		// 等待用户答案（除了第一个问题）
		if questionIndex > 0 {
			log.Printf("[Interview Engine] Waiting for answer, sessionID: %s, questionIndex: %d", session.SessionID, questionIndex)
			answer, received := WaitForAnswerWithHeartbeat(ctx, e.sessionManager, session.SessionID, answerTimeout, heartbeatInterval, e.writer)
			if !received {
				log.Printf("[Interview Engine] Answer timeout, sessionID: %s", session.SessionID)
				SendErrorEvent(e.writer, "等待答案超时，面试已结束")
				SendCompleteEvent(e.writer)
				break
			}

			if answer == "quit" {
				SendCompleteEvent(e.writer)
				break
			}

			// 保存用户回答
			session.AllDialogues = append(session.AllDialogues, map[string]interface{}{
				"speaker_type":  "candidate",
				"content":       answer,
				"display_order": uint32(questionIndex)*100 + uint32(followUpCount),
			})

			// 决定是否继续追问或进入下一个维度
			if followUpCount < maxFollowUps {
				followUpCount++
			} else {
				dimensionIndex++
				followUpCount = 0
				questionIndex++
				if dimensionIndex >= len(dimensions) {
					SendTopicCompleteEvent(e.writer)
					SendCompleteEvent(e.writer)
					break
				}
			}
		}

		if questionIndex == 0 {
			questionIndex++
		}

		// 生成问题提示词
		prompt := e.buildPrompt(questionIndex, session, resumeContent, dimensions[dimensionIndex], followUpCount)

		// 调用智能体生成问题
		isFirstQuestion := questionIndex == 1 && followUpCount == 0
		result, err := service.GenerateInterviewQuestions(ctx, prompt, session.UserID, isFirstQuestion)
		if err != nil {
			SendErrorEvent(e.writer, "Failed to generate question: "+err.Error())
			SendCompleteEvent(e.writer)
			break
		}

		if len(result.Questions) == 0 {
			SendTopicCompleteEvent(e.writer)
			SendCompleteEvent(e.writer)
			break
		}

		// 缓存简历内容（仅第一次）
		if questionIndex == 1 {
			resumeContent = session.Query
		}

		// 发送问题事件
		q := result.Questions[0]
		SendQuestionEvent(e.writer, questionIndex, q)

		// 保存问题
		session.AllQuestions = append(session.AllQuestions, map[string]interface{}{
			"question_text":  q.QuestionText,
			"eval_dimension": q.EvalDimension,
			"order":          q.Order,
		})

		// 保存提问对话
		if len(result.Dialogues) > 0 {
			for _, d := range result.Dialogues {
				if d.SpeakerType == "interviewer" {
					session.AllDialogues = append(session.AllDialogues, map[string]interface{}{
						"speaker_type":  "interviewer",
						"content":       d.Content,
						"display_order": uint32(questionIndex)*100 + uint32(followUpCount),
					})
					break
				}
			}
		}

		// 发送就绪事件
		SendReadyEventWithSession(e.writer, questionIndex, session.SessionID)
		e.sessionManager.ClearAnswer(session.SessionID)
	}

	// 保存所有面试数据
	e.saveInterviewData(ctx, session)
}

// buildPrompt 构建提示词
func (e *InterviewEngine) buildPrompt(questionIndex int, session *InterviewSession, resumeContent string, dimension string, followUpCount int) string {
	return service.BuildInterviewPrompt(questionIndex, session.Query, session.ResumeId, session.HasResume, dimension, followUpCount, session.Type, session.Domain, session.Difficulty)
}

// ConvertToInterfaceSlice 将 []map[string]interface{} 转换为 []interface{}
func ConvertToInterfaceSlice(data []map[string]interface{}) []interface{} {
	result := make([]interface{}, len(data))
	for i, item := range data {
		result[i] = item
	}
	return result
}

// getDimensions 根据面试类型获取维度
func (e *InterviewEngine) getDimensions(interviewType string) []string {
	if interviewType == "综合面试" {
		return []string{
			"professional_field",
			"project_experience",
			"technical_depth",
			"technical_foundation",
			"team_collaboration",
			"system_architecture_design",
		}
	}
	// 专项面试
	return []string{
		"basic_knowledge_mastery",
		"working_principle_practical_experience",
		"advanced_features_application",
		"problem_troubleshooting_skills",
		"architecture_design_thinking",
		"performance_optimization_ability",
	}
}

// saveInterviewData 保存面试数据
func (e *InterviewEngine) saveInterviewData(ctx context.Context, session *InterviewSession) {
	log.Printf("[Interview Engine] Saving dialogues, sessionID: %s", session.SessionID)

	// 带重试机制的保存
	const maxRetries = 3
	var saveErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		saveErr = e.interviewSvc.SaveInterviewDialogues(ctx, session.UserID, session.RecordID, ConvertToInterfaceSlice(session.AllQuestions), ConvertToInterfaceSlice(session.AllDialogues))
		if saveErr == nil {
			break
		}

		// 判断是否为可重试的错误
		if !IsRetryableError(saveErr) {
			if saveErr != context.Canceled && saveErr != context.DeadlineExceeded {
				alert.SendDatabaseErrorAlert(
					fmt.Sprintf("SaveInterviewDialogues (不可重试) - UserID: %d, RecordID: %d", session.UserID, session.RecordID),
					saveErr,
					attempt+1,
				)
			}
			break
		}

		// 最后一次尝试失败
		if attempt == maxRetries-1 {
			alert.SendDatabaseErrorAlert(
				fmt.Sprintf("SaveInterviewDialogues (重试耗尽) - UserID: %d, RecordID: %d", session.UserID, session.RecordID),
				saveErr,
				maxRetries,
			)
			break
		}

		// 指数退避
		backoffDuration := time.Duration(100*(1<<uint(attempt))) * time.Millisecond
		time.Sleep(backoffDuration)
	}

	// 更新面试记录状态
	duration := int64(time.Since(session.StartTime).Seconds())
	e.updateInterviewRecord(ctx, session, duration)

	// 发布评估消息
	e.publishEvaluationMessages(ctx, session)

	log.Printf("[Interview Engine] Interview completed, sessionID: %s", session.SessionID)
}

// updateInterviewRecord 更新面试记录
func (e *InterviewEngine) updateInterviewRecord(ctx context.Context, session *InterviewSession, duration int64) {
	log.Printf("[Interview Engine] Updating interview record, sessionID: %s, duration: %d seconds", session.SessionID, duration)
	// 实现更新逻辑（可根据需要调用服务）
	_ = duration
}

// publishEvaluationMessages 发布评估消息
func (e *InterviewEngine) publishEvaluationMessages(ctx context.Context, session *InterviewSession) {
	log.Printf("[Interview Engine] Publishing evaluation messages, sessionID: %s, userID: %d, recordID: %d", session.SessionID, session.UserID, session.RecordID)

	// 发布评估报告生成消息
	if err := mq.PublishEvaluationReport(ctx, session.UserID, session.RecordID); err != nil {
		log.Printf("[Interview Engine] Failed to publish evaluation report message: %v, sessionID: %s", err, session.SessionID)
	}

	// 发布主题评估消息
	if err := mq.PublishTopicEvaluation(ctx, session.UserID, session.RecordID); err != nil {
		log.Printf("[Interview Engine] Failed to publish topic evaluation message: %v, sessionID: %s", err, session.SessionID)
	}
}

// cleanup 清理资源
func (e *InterviewEngine) cleanup(session *InterviewSession) {
	if r := recover(); r != nil {
		SendErrorEvent(e.writer, fmt.Sprintf("面试异常: %v", r))
	}
	// 延迟删除会话
	go func() {
		time.Sleep(10 * time.Second)
		e.sessionManager.DeleteSession(session.SessionID)
	}()
}
