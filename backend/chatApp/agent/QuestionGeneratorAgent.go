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
func NewQuestionGeneratorAgent() adk.Agent {
	ctx := context.Background()

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionGeneratorAgent",
		Description: "根据前面的简历去分析结果去提问问题",
		Instruction: `你是一名资深的技术面试官，负责围绕候选人的背景开展模拟面试，每轮仅生成一个问题并等待回答。

工作流程：
	1. 调用工具 "gen_question" 获取问题的增强方案，调用时必须提供字段：
    2. 生成对应的问题后调用 "ask_for_clarification" 工具获取用户回答。
	3. 进入下一轮循环。


关键约束：
- 保持专业、鼓励且以候选人为中心的语气，明确当前面试阶段（如“技术深挖”、“总结反馈”等）。`,
		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreateGenQuestionTool(),
					tool2.NewAskForInputTool(),
				},
			},
		},
		MaxIterations: 12,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}

	return a
}
