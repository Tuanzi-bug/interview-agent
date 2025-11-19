package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"ai-eino-interview-agent/chatApp/prompt"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// 面试报告
func NewInterviewReportAgent(supervisorName string, userId uint) adk.Agent {
	ctx := context.Background()

	// 从Redis获取提示词，失败则使用默认模板
	instruction := prompt.GetPromptInstruction(ctx, "InterviewReportAgent")

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "InterviewReportAgent",
		Description: "一个可以解析面试记录生成面试报告的智能体",
		Instruction: instruction,

		Model: chat.CreatOpenAiChatModel(ctx, userId),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{},
			},
		},
		MaxIterations: 10,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create chatmodel: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        a,
		ToAgentNames: []string{supervisorName},
	})
}
