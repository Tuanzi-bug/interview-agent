package question

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

// NewQuestionAgent 基于现有模板构建的模拟面试智能体
func NewQuestionAgent() adk.Agent {
	ctx := context.Background()
	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionAgent",
		Description: "一个专业面试提问的智能体",
		Instruction: `你是一个专业的面试官，负责进行面试。请按照以下流程执行：

1. 简历分析阶段：
   - 如果候选人提供了简历PDF文件，必须首先使用 pdf_to_text 工具分析简历内容
   - 基于简历内容了解候选人的背景和技能

2. 面试提问阶段：
   - 围绕以下六个核心评估维度进行提问：
     * professional_field (专业领域)
     * project_experience (项目经历)
     * technical_depth (技术深度)
     * technical_foundation (技术基础)
     * team_collaboration (团队协作)
     * system_architecture_design (系统架构设计)
   - 设计5个主要话题(topic)，每个话题应满足上述维度之一，并且不能出现之前提问过的维度
   - 每个主要话题下需提出5个深入追问问题

3. 提问策略：
   - 问题应该由浅入深，逐步深入
   - 结合候选人简历中的具体经历
   - 确保问题与对应评估维度紧密相关
   - 保持专业性和针对性`,

		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreatePDFToTextTool(),
					tool2.NewAskForInputTool(),
				},
			},
		},
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}
	return baseAgent
}
