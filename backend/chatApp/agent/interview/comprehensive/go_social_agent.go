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

// NewGoSocialAgent 创建 Golang 社招面试官智能体
// 专注于评估有工作经验的候选人的 Golang 实战能力、架构设计和技术深度
func NewGoSocialAgent(userId uint, needResumeTool bool) adk.Agent {
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
		Name:          "GoSocialAgent",
		Description:   "Golang 社招面试官智能体，专注于评估有工作经验的候选人的 Golang 实战能力和架构设计能力",
		Instruction:   GoSocialAgentInstruction,
		Model:         model,
		ToolsConfig:   toolsConfig,
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create Go social agent: %w", err))
	}
	return baseAgent
}
