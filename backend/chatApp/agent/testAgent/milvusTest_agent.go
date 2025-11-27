package testAgent

import (
	"ai-eino-interview-agent/chatApp/chat"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	componenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"golang.org/x/net/context"
)

// NewEvaluationAgent 用于生成评估报告的智能体
func NewMilvuesAgent(userId uint) adk.Agent {
	ctx := context.Background()

	// 构建系统指令
	instruction := buildEvaluationInstruction()

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "EvaluationAgent",
		Description: "一个专业问答检索的智能体",
		Instruction: instruction,

		Model: chat.CreatOpenAiChatModel(ctx, userId),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					//tool2.GetInterviewsDataTool(),
					tool2.GetMilvusRetrieverTool(),
				},
			},
		},
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create evaluation agent: %w", err))
	}
	return baseAgent
}

// buildEvaluationInstruction 构建评估智能体的系统指令
func buildEvaluationInstruction() string {
	return `
你是一个有帮助的问答助手。通过使用检索milvus数据库然后进行格式优化重写，回答用户提出的问题。`
}
