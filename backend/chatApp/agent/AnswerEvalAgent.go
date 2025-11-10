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

// NewAnswerEvalAgent 基于现有模板构建的模拟面试智能体
func NewAnswerEvalAgent() adk.Agent {
	ctx := context.Background()

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "AnswerEvalAgent",
		Description: "根据前面的简历分析和问题,对用户的答案进行评估",
		Instruction: `你是一名资深的技术面试官 你要去调用 "answer_eval" 工具去获取评估的方向生成评估结果
			流程：
			1. 调用 "answer_eval" 工具获取评估的方向。
			2. 根据评估的方向生成评估结果。
			3. 返回评估结果。
			约束：
			1. 评估结果必须是中文。
		`,
		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreateAnswerEvalTool(),
				},
			},
		},
		MaxIterations: 12,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create answer eval agent: %w", err))
	}

	return a
}
