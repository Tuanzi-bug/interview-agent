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
func NewQuestionAgent(userId uint) adk.Agent {
	ctx := context.Background()
	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionAgent",
		Description: "一个专业面试提问的智能体",
		Instruction: `你是一个专业的面试官。根据提供的简历或信息，快速生成一个面试问题和对话。

任务：
1. 如果提供了PDF路径，使用pdf_to_text工具解析简历，只有在面试刚开始的时候调用一次即可
2. 面试提问阶段：
   - 根据候选人的背景和技能，围绕以下核心评估维度进行提问：
     * professional_field (专业领域)
     * project_experience (项目经历)
     * technical_depth (技术深度)
     * technical_foundation (技术基础)
     * team_collaboration (团队协作)
     * system_architecture_design (系统架构设计)
   - 围绕这6个主要话题，每个话题应满足上述维度之一，进行提问，每个纬度的主题只出现一遍，都问完了就退出面试
   - 每个主要话题下提出2个深入追问问题

3. 提问策略：
   - 问题应该由浅入深，逐步深入
   - 结合候选人提供的背景信息
   - 确保问题与对应评估维度紧密相关
   - 保持专业性和针对性
4. 只返回JSON，不要返回其他文本

必须返回的JSON格式：
{
  "questions": [{
    "question_text": "问题内容",
    "eval_dimension": "纬度",
    "order": 1
  }],
  "dialogues": [
    {"speaker_type": "interviewer", "content": "提问", "display_order": 1},
    {"speaker_type": "candidate", "content": "回答", "display_order": 2},
    {"speaker_type": "interviewer", "content": "追问", "display_order": 3}
  ]
}`,

		Model: chat.CreatOpenAiChatModel(ctx, userId),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreatePDFToTextTool(),
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
