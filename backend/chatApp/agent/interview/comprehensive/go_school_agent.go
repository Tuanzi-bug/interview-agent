package comprehensive

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

// NewGoSchoolAgent 创建 Golang 校招面试官智能体
// 专注于评估应届毕业生的 Golang 基础知识、学习能力和潜力
func NewGoSchoolAgent(userId uint, needResumeTool bool) adk.Agent {
	ctx := context.Background()

	var toolsConfig adk.ToolsConfig
	if needResumeTool {
		toolsConfig = adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.GetResumeInfoTool(),
				},
			},
		}
	}

	model, err := chat.CreatOpenAiChatModel(ctx, userId)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create OpenAI chat model: %w", err))
	}

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "GoSchoolAgent",
		Description:   "Golang 校招面试官智能体，专注于评估应届毕业生的 Golang 基础和学习潜力",
		Instruction:   GoSchoolAgentInstruction,
		Model:         model,
		ToolsConfig:   toolsConfig,
		MaxIterations: 15,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create Go school agent: %w", err))
	}
	return baseAgent
}
