package ext

import (
	"ai-eino-interview-agent/chatApp/agent"
	"ai-eino-interview-agent/chatApp/tool"
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"sync"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// InterviewResult 面试结果
type InterviewResult struct {
	// 完整的对话历史
	Messages []*Message `json:"messages"`
	// 最终报告
	Report string `json:"report,omitempty"`
	// 状态：resume_analysis, question_generation, answer_evaluation, report_generation, completed
	Status string `json:"status"`
	// 当前活跃的 Agent 名称
	CurrentAgent string `json:"current_agent,omitempty"`
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`            // user, assistant, system
	Content string `json:"content"`         // 消息内容
	Agent   string `json:"agent,omitempty"` // 发送消息的 Agent 名称
}

// InterviewEvent 面试事件（用于流式响应）
type InterviewEvent struct {
	// 事件类型：message, transfer, error, done
	Type string `json:"type"`
	// Agent 名称
	AgentName string `json:"agent_name,omitempty"`
	// 消息内容（当 Type 为 message 时）
	Message string `json:"message,omitempty"`
	// 转让目标（当 Type 为 transfer 时）
	TransferTo string `json:"transfer_to,omitempty"`
	// 错误信息（当 Type 为 error 时）
	Error string `json:"error,omitempty"`
	// 状态更新
	Status *string `json:"status,omitempty"`
	// 最终报告（当 Type 为 done 时）
	Report string `json:"report,omitempty"`
	// 面试评分（当 Type 为 done 时）
	Score *float64 `json:"score,omitempty"`
	// 面试时长（秒）（当 Type 为 done 时）
	Duration int64 `json:"duration,omitempty"`
	// 反馈信息（当 Type 为 done 时）
	Feedback string `json:"feedback,omitempty"`
	// 对话历史（JSON 格式）（当 Type 为 done 时）
	Messages string `json:"messages,omitempty"`
}

var (
	// 全局 runner 实例（单例模式，避免重复创建）
	globalRunner     *adk.Runner
	globalRunnerOnce sync.Once
)

// getRunner 获取全局 runner 实例（懒加载单例）
func getRunner(ctx context.Context) *adk.Runner {
	globalRunnerOnce.Do(func() {
		interviewSupervisor := agent.NewInterviewSupervisorAgent()
		globalRunner = adk.NewRunner(ctx, adk.RunnerConfig{
			Agent: interviewSupervisor,
		})
	})
	return globalRunner
}

// StartInterviewStream 启动面试流程（流式）
// query: 用户输入的查询
// maxQuestions: 最大提问数量，0 表示无限制
// 返回一个事件流 channel，可以实时获取 Agent 的输出
// 注意：调用者需要负责关闭 channel（当事件流结束时，channel 会自动关闭）
func StartInterviewStream(ctx context.Context, query string, maxQuestions int) (<-chan *InterviewEvent, error) {
	runner := getRunner(ctx)

	// 将 maxQuestions 添加到 context 中，供 agents 使用
	if maxQuestions > 0 {
		// 检查是否已有 session ID（用于继续面试时复用）
		sessionID, ok := ctx.Value(tool.SessionIDKey).(string)
		if !ok || sessionID == "" {
			// 生成唯一的 session ID
			sessionID = tool.GenerateSessionID()
		}
		ctx = context.WithValue(ctx, tool.MaxQuestionsKey, maxQuestions)
		ctx = context.WithValue(ctx, tool.QuestionCountKey, 0)
		ctx = context.WithValue(ctx, tool.SessionIDKey, sessionID)
	}

	// 构建消息
	messages := []adk.Message{
		schema.UserMessage(query),
	}

	// 启动面试流程
	iter := runner.Run(ctx, messages)

	// 创建事件 channel
	eventChan := make(chan *InterviewEvent, 10)

	// 在 goroutine 中处理事件流
	go func() {
		defer close(eventChan)

		// 用于收集对话历史
		var messageHistory []Message

		for {
			event, ok := iter.Next()
			if !ok {
				// 发送完成事件
				eventChan <- &InterviewEvent{
					Type: "done",
				}
				break
			}

			if event.Err != nil {
				// 发送错误事件
				eventChan <- &InterviewEvent{
					Type:  "error",
					Error: event.Err.Error(),
				}
				return
			}

			// 处理转让事件
			if event.Action != nil && event.Action.TransferToAgent != nil {
				status := ""
				switch event.Action.TransferToAgent.DestAgentName {
				case "ResumeAnalysisAgent":
					status = "resume_analysis"
				case "QuestionGeneratorAgent":
					status = "question_generation"
				case "AnswerEvalAgent":
					status = "answer_evaluation"
				case "InterviewReportAgent":
					status = "report_generation"
				}

				// 将对话历史转换为 JSON
				filtered := filterInterviewMessages(messageHistory)
				var messagesJSON []byte
				if len(filtered) > 0 {
					messagesJSON, _ = json.Marshal(filtered)
				}

				eventChan <- &InterviewEvent{
					Type:       "transfer",
					AgentName:  event.AgentName,
					TransferTo: event.Action.TransferToAgent.DestAgentName,
					Status:     &status,
					Messages:   string(messagesJSON),
				}
				continue
			}

			// 处理Agent输出
			if event.Output != nil && event.Output.MessageOutput != nil {
				messageContent := event.Output.MessageOutput.Message.Content
				role := event.Output.MessageOutput.Message.Role

				// 收集对话历史（包含 Agent 信息，方便后续过滤）
				messageHistory = append(messageHistory, Message{
					Role:    string(role),
					Content: messageContent,
					Agent:   event.AgentName,
				})

				// 将对话历史转换为 JSON（只保留面试官提问）
				filtered := filterInterviewMessages(messageHistory)
				var messagesJSON []byte
				if len(filtered) > 0 {
					messagesJSON, _ = json.Marshal(filtered)
				}

				// 如果是报告Agent的输出，发送完成事件并返回
				if event.AgentName == "InterviewReportAgent" {
					completed := "completed"
					var score *float64
					score = extractScoreFromReport(messageContent)
					eventChan <- &InterviewEvent{
						Type:     "done",
						Status:   &completed,
						Score:    score,
						Report:   messageContent,
						Messages: string(messagesJSON),
					}
					return
				}

				// 其他Agent的输出，发送消息事件
				eventChan <- &InterviewEvent{
					Type:      "message",
					AgentName: event.AgentName,
					Message:   messageContent,
					Messages:  string(messagesJSON),
				}
			}
		}
	}()

	return eventChan, nil
}

// ContinueInterview 继续面试流程（用于多轮对话）
func ContinueInterview(ctx context.Context, query string, maxQuestions int) (<-chan *InterviewEvent, error) {

	return StartInterviewStream(ctx, query, maxQuestions)
}

// filterInterviewMessages 只保留面试官提问
// 约定：
// - 面试官提问：role == "assistant" 且 Agent == "QuestionGeneratorAgent"
func filterInterviewMessages(history []Message) []Message {
	var result []Message
	for _, m := range history {
		if m.Role == "assistant" && m.Agent == "QuestionGeneratorAgent" {
			result = append(result, m)
		}
	}
	return result
}

// extractScoreFromReport 从面试报告中提取分数
// 适用场景：报告格式相对固定，分数表现形式一致

func extractScoreFromReport(report string) *float64 {
	// 匹配模式：
	// - "总分：85" 或 "总分: 85"
	// - "综合评分：82.5" 或 "综合评分: 82.5"
	// - "面试评分：90分" 或 "面试评分: 90分"

	patterns := []string{
		`总(?:体)?(?:评)?分[：:]\s*(\d+(?:\.\d+)?)`,  // 总分：85 或 总体评分：82.5
		`综合(?:评)?分[：:]\s*(\d+(?:\.\d+)?)`,       // 综合评分：85
		`面试(?:评)?分[：:]\s*(\d+(?:\.\d+)?)(?:分)?`, // 面试评分：85分
		`(?:最终)?评分[：:]\s*(\d+(?:\.\d+)?)`,       // 评分：85
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(report)
		if len(matches) > 1 {
			score, err := strconv.ParseFloat(matches[1], 64)
			if err == nil && score >= 0 && score <= 100 {
				return &score
			}
		}
	}

	return nil
}
