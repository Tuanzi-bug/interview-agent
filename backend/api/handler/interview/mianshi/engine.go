package mianshi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"ai-eino-interview-agent/chatApp/agent_service/interview"
	"ai-eino-interview-agent/internal/model"
	interviewservice "ai-eino-interview-agent/internal/service/interviews"
)

// InterviewDialogueData 对话数据结构
type InterviewDialogueData struct {
	Question string
	Answer   string
}

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

	questionIndex := 0
	const answerTimeout = 30 * time.Minute
	const heartbeatInterval = 15 * time.Second
	const maxQuestions = 6 // 最多问6个主问题

	// 用于存储主问题及其追问
	var mainQuestion *InterviewDialogueData
	var followUpQuestions []*InterviewDialogueData

	// 用于存储历史问题和回答
	type QuestionAnswer struct {
		Question string
		Answer   string
	}
	var questionHistory []QuestionAnswer

	// 用于追踪当前主问题是否已回答
	mainQuestionAnswered := false
	// 用于追踪已回答的追问数量
	answeredFollowUpCount := 0

	// 创建智能体服务（在循环外创建一次，避免重复实例化）
	agentSvc := interview.NewInterviewAgentService(session.UserID)

	// 确定智能体类型（根据面试类型和领域选择）
	agentType := e.selectAgentType(session)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[Interview Engine] Context cancelled, sessionID: %s", session.SessionID)
			return
		default:
		}

		// 等待用户答案
		if questionIndex > 0 {
			log.Printf("[Interview Engine] Waiting for answer, sessionID: %s, questionIndex: %d", session.SessionID, questionIndex)
			answer, received := WaitForAnswerWithHeartbeat(ctx, e.sessionManager, session.SessionID, answerTimeout, heartbeatInterval, e.writer)
			if !received {
				log.Printf("[Interview Engine] Answer timeout, sessionID: %s", session.SessionID)
				SendCompleteEvent(e.writer)
				break
			}

			// 保存用户回答到当前主问题或追问
			if mainQuestion != nil {
				if !mainQuestionAnswered {
					// 第一个回答：主问题
					mainQuestion.Answer = answer
					mainQuestionAnswered = true
				} else {
					// 后续回答：追问（按顺序赋值）
					if answeredFollowUpCount < len(followUpQuestions) {
						followUpQuestions[answeredFollowUpCount].Answer = answer
						answeredFollowUpCount++
					}
				}
			}

			// 检查是否还有未回答的追问
			hasUnansweredFollowUp := false
			for _, followUp := range followUpQuestions {
				if followUp.Answer == "" {
					hasUnansweredFollowUp = true
					break
				}
			}

			// 如果还有未回答的追问，发送下一个追问并继续等待
			if hasUnansweredFollowUp {
				// 发送下一个追问事件
				for i, followUp := range followUpQuestions {
					if followUp.Answer == "" {
						errSend := SendSSEEvent(e.writer, map[string]interface{}{
							"type":  "follow_up_question",
							"index": questionIndex,
							"order": i + 1,
							"data": map[string]interface{}{
								"question_text": followUp.Question,
							},
						})
						if errSend != nil {
							return
						}
						break
					}
				}
				// 继续循环，在下一次迭代中等待用户回答追问
				continue
			}

			// 所有追问都已回答，保存当前主问题及其追问，进入下一个主问题
			if mainQuestion != nil {
				e.saveDialogueData(ctx, session, mainQuestion, followUpQuestions)
				// 将主问题和回答添加到历史记录
				questionHistory = append(questionHistory, QuestionAnswer{
					Question: mainQuestion.Question,
					Answer:   mainQuestion.Answer,
				})
			}
			// 重置追问和状态
			followUpQuestions = nil
			mainQuestion = nil
			mainQuestionAnswered = false
			answeredFollowUpCount = 0

			// 检查是否已达到最大问题数，如果是则结束面试
			if questionIndex >= maxQuestions {
				log.Printf("[Interview Engine] Reached max questions (%d), ending interview, sessionID: %s", maxQuestions, session.SessionID)
				SendTopicCompleteEvent(e.writer)
				SendCompleteEvent(e.writer)
				break
			}
		}

		// 构建问题提示词
		var prompt string
		if questionIndex == 0 {
			// 第一个问题：只传递简历ID和难度
			prompt = fmt.Sprintf("请根据简历ID和难度等级生成一个面试问题。\n简历ID: %d\n难度等级: %s", session.ResumeId, session.Difficulty)
		} else {
			// 后续问题：包含历史回答，让智能体根据用户的回答调整后续问题
			historyText := ""
			for i, qa := range questionHistory {
				historyText += fmt.Sprintf("问题%d：%s\n回答%d：%s\n\n", i+1, qa.Question, i+1, qa.Answer)
			}
			prompt = fmt.Sprintf(`根据简历ID、难度等级和用户已回答的问题，生成下一个面试问题。

简历ID: %d
难度等级: %s

用户已回答的问题：
%s

请根据用户的回答情况和难度等级，生成下一个更有针对性的面试问题。`, session.ResumeId, session.Difficulty, historyText)
		}

		// 收集智能体的响应
		var questionResult map[string]interface{}
		err := agentSvc.RunInterviewWithCallback(ctx, agentType, session.HasResume, prompt, func(message string) error {
			// 解析响应
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(message), &result); err == nil {
				questionResult = result
			}
			return nil
		})

		if err != nil {
			log.Printf("[Interview Engine] Failed to generate question: %v, sessionID: %s", err, session.SessionID)
			SendErrorEvent(e.writer, "Failed to generate question: "+err.Error())
			SendCompleteEvent(e.writer)
			break
		}

		if len(questionResult) == 0 {
			log.Printf("[Interview Engine] Agent returned empty question result, ending interview, sessionID: %s", session.SessionID)
			SendTopicCompleteEvent(e.writer)
			SendCompleteEvent(e.writer)
			break
		}

		// 提取主问题信息（智能体返回格式：main_question.question_text）
		var questionText string
		var mainQ map[string]interface{}

		if mq, ok := questionResult["main_question"].(map[string]interface{}); ok {
			mainQ = mq
			if qt, ok := mainQ["question_text"].(string); ok {
				questionText = qt
			}
		}

		if questionText == "" {
			log.Printf("[Interview Engine] Failed to extract question text from agent response, sessionID: %s", session.SessionID)
			SendErrorEvent(e.writer, "Failed to extract question from agent response")
			SendCompleteEvent(e.writer)
			break
		}

		// 发送问题事件
		err = SendSSEEvent(e.writer, map[string]interface{}{
			"type":  "question",
			"index": questionIndex,
			"data": map[string]interface{}{
				"question_text": questionText,
			},
		})
		if err != nil {
			return
		}

		// 提取主问题的提问内容
		var questionContent string
		if mainQ != nil {
			if qt, ok := mainQ["question_text"].(string); ok {
				questionContent = qt
			}
		}

		// 创建主问题
		if mainQuestion == nil {
			// 这是一个新的主问题
			mainQuestion = &InterviewDialogueData{
				Question: questionContent,
				Answer:   "",
			}
			// 清空追问列表，准备接收新的追问
			followUpQuestions = nil
		}

		// 提取追问列表（智能体返回格式：follow_up_questions）
		if followUps, ok := questionResult["follow_up_questions"].([]interface{}); ok {
			for _, fu := range followUps {
				if followUpMap, ok := fu.(map[string]interface{}); ok {
					if questionText, ok := followUpMap["question_text"].(string); ok {
						followUpQuestions = append(followUpQuestions, &InterviewDialogueData{
							Question: questionText,
							Answer:   "",
						})
					}
				}
			}
		}

		log.Printf("[Interview Engine] Generated question, sessionID: %s, questionIndex: %d, follow-ups: %d", session.SessionID, questionIndex, len(followUpQuestions))

		// 发送就绪事件
		SendReadyEventWithSession(e.writer, questionIndex, session.SessionID)
		e.sessionManager.ClearAnswer(session.SessionID)

		// 更新会话中的问题计数
		session.QuestionCount = int32(questionIndex)

		// 问题已发送，准备等待用户回答
		// 在下一次循环中，questionIndex > 0 会进入等待答案的逻辑
		questionIndex++
	}

	// 保存最后一个主问题及其追问
	if mainQuestion != nil {
		e.saveDialogueData(ctx, session, mainQuestion, followUpQuestions)
		// 更新最后一个问题的计数
		session.QuestionCount = int32(questionIndex + 1)
	}
}

// saveDialogueData 保存单个主问题及其追问到数据库
func (e *InterviewEngine) saveDialogueData(ctx context.Context, session *InterviewSession, mainQuestion *InterviewDialogueData, followUpQuestions []*InterviewDialogueData) {
	// 构建主问题对象
	// ParentID 会在 SaveInterviewDialogueWithParent 中设置为 0（表示主问题）
	mainQ := &model.InterviewDialogue{
		UserID:   session.UserID,
		ReportID: session.RecordID,
		Question: mainQuestion.Question,
		Answer:   mainQuestion.Answer,
	}

	// 构建追问对象列表
	// ParentID 会在 SaveInterviewDialogueWithParent 中设置为主问题的 ID
	var followUps []*model.InterviewDialogue
	for _, followUp := range followUpQuestions {
		followUps = append(followUps, &model.InterviewDialogue{
			UserID:   session.UserID,
			ReportID: session.RecordID,
			Question: followUp.Question,
			Answer:   followUp.Answer,
		})
	}

	// 直接保存到数据库
	// SaveInterviewDialogueWithParent 会自动处理 ParentID：
	// - 主问题的 ParentID = 0
	// - 追问的 ParentID = 主问题的 ID
	if err := e.interviewSvc.SaveInterviewDialogueWithParent(ctx, session.UserID, session.RecordID, mainQ, followUps); err != nil {
		log.Printf("[Interview Engine] Failed to save dialogue data to database: %v, sessionID: %s", err, session.SessionID)
		SendErrorEvent(e.writer, "Failed to save interview data: "+err.Error())
		return
	}

	log.Printf("[Interview Engine] Saved dialogue data to database, sessionID: %s, main question: %s, follow-ups: %d", session.SessionID, mainQuestion.Question, len(followUpQuestions))
}

// selectAgentType 根据面试类型和领域选择智能体类型
func (e *InterviewEngine) selectAgentType(session *InterviewSession) interview.InterviewAgentType {
	// 综合面试
	if session.Type == "综合面试" {
		switch session.Domain {
		case "校招简历面试":
			return interview.ComprehensiveSchool
		case "社招简历面试":
			return interview.ComprehensiveSocial
		default:
			// 社招简历面试为默认选项
			return interview.ComprehensiveSocial
		}
	}

	// 专项面试
	switch session.Domain {
	case "Java":
		return interview.SpecializedJava
	case "MQ":
		return interview.SpecializedMQ
	case "MySQL":
		return interview.SpecializedMySQL
	case "Redis":
		return interview.SpecializedRedis
	case "Go":
		fallthrough
	default:
		return interview.SpecializedGo
	}
}
