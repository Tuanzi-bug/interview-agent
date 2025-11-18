package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"ai-eino-interview-agent/chatApp/prompt"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewAnswerEvalAgent 基于现有模板构建的模拟面试智能体
func NewAnswerEvalAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	// 从Redis获取提示词，失败则使用默认模板
	instruction := prompt.GetPromptInstruction(ctx, "AnswerEvalAgent")

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "AnswerEvalAgent",
		Description: "根据简历分析、问题和答案进行多维度评估，生成结构化评估报告",
		Instruction: instruction,
		Model:       chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreateScoreExtractionTool(),
				},
			},
		},
		MaxIterations: 12,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create answer eval agent: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        a,
		ToAgentNames: []string{supervisorName},
	})
}
