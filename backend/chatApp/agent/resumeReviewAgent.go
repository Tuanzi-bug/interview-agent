package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"ai-eino-interview-agent/chatApp/prompt"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewResumeReviewAgent 创建简历评审Agent
// 专门负责评审简历、给出分数和优化建议
func NewResumeReviewAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	// 从Redis获取提示词，失败则使用默认模板
	instruction := prompt.GetPromptInstruction(ctx, "ResumeReviewAgent")

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ResumeReviewAgent",
		Description: "一个可以评审简历、给出分数和优化建议的智能体",
		Instruction: instruction,
		Model:       chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{
					tool2.CreatePDFToTextTool(),   // PDF解析工具
					tool2.CreateResumeScoreTool(), // 简历评分工具
				},
			},
		},
		MaxIterations: 10,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create resume review agent: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}
