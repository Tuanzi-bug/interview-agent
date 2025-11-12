package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"fmt"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"log"
)

// 面试报告
func NewInterviewReportAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "InterviewReportAgent",
		Description: "一个可以解析面试记录生成面试报告的智能体",
		Instruction: `你是一名资深的面试报告专家，负责对用户的面试记录进行分析,并输出对应的面试报告。`,

		Model: chat.CreatOpenAiChatModel(ctx),
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
