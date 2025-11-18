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

// NewMockInterviewAgent 基于现有模板构建的模拟面试智能体
func NewQuestionGeneratorAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	// 从Redis获取提示词，失败则使用默认模板
	instruction := prompt.GetPromptInstruction(ctx, "QuestionGeneratorAgent")

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
