package ext

import (
	"ai-eino-interview-agent/chatApp/agent"
	"context"
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
	Status string `json:"status,omitempty"`
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
// 返回一个事件流 channel，可以实时获取 Agent 的输出
// 注意：调用者需要负责关闭 channel（当事件流结束时，channel 会自动关闭）
func StartInterviewStream(ctx context.Context, query string) (<-chan *InterviewEvent, error) {
	runner := getRunner(ctx)

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
				case "ResumAnalysisAgent":
					status = "resume_analysis"
				case "QuestionGeneratorAgent":
					status = "question_generation"
				case "AnswerEvalAgent":
					status = "answer_evaluation"
				case "InterviewReportAgent":
					status = "report_generation"
				}

				eventChan <- &InterviewEvent{
					Type:       "transfer",
					AgentName:  event.AgentName,
					TransferTo: event.Action.TransferToAgent.DestAgentName,
					Status:     status,
				}
				continue
			}

			// 处理Agent输出
			if event.Output != nil && event.Output.MessageOutput != nil {
				eventChan <- &InterviewEvent{
					Type:      "message",
					AgentName: event.AgentName,
					Message:   event.Output.MessageOutput.Message.Content,
				}

				// 如果是报告Agent的输出，发送完成事件
				if event.AgentName == "InterviewReportAgent" {
					eventChan <- &InterviewEvent{
						Type:   "done",
						Status: "completed",
					}
				}
			}
		}
	}()

	return eventChan, nil
}

// ContinueInterview 继续面试流程（用于多轮对话）
func ContinueInterview(ctx context.Context, query string) (<-chan *InterviewEvent, error) {

	return StartInterviewStream(ctx, query)
}
