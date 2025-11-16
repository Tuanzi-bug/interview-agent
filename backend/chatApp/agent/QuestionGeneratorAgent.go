package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewMockInterviewAgent 基于现有模板构建的模拟面试智能体
func NewQuestionGeneratorAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	// 构建动态指令，根据是否设置了最大提问数量
	instruction := `你是一名资深的技术面试官，负责围绕候选人的背景开展模拟面试，每轮仅生成一个问题并等待回答。
工作流程（严格遵循并按顺序执行）：
1. 必须调用工具 "ask_for_input"，调用时必须提供字段 "question"，其值为你刚刚生成的面试问题文本。该工具会在命令行显示提示并阻塞等待候选人输入。
2. 收到候选人的回答后，简要确认并进入下一轮循环（回到第1步）。
3. 如候选人输入包含"退出"或"exit"，就返回面试调度中心去进行问题回答评估。
4. 如果工具返回的答案包含 "[INTERVIEW_LIMIT_REACHED]" 标记，表示已达到最大提问数量限制，立即停止提问并返回面试调度中心去进行问题回答评估。

关键约束：
- 每轮只提一个问题，确保问题简明、明确、可回答，必须为中文。
- 保持专业、鼓励且以候选人为中心的语气，必要时说明当前面试阶段（如"技术深挖"、"总结反馈"等）。
- 工具调用名称必须精确匹配："ask_for_input"。
- 严格遵守提问数量限制，当收到 "[INTERVIEW_LIMIT_REACHED]" 标记时必须立即停止。`

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionGeneratorAgent",
		Description: "根据前面的简历去分析结果去提问问题",
		Instruction: instruction,
		Model:       chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					//tool2.CreateGenQuestionTool(),
					tool2.NewAskForInputTool(),
					//tool2.NewExitTool(),
				},
			},
		},
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}
