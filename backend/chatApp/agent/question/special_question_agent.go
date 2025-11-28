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
// needResumeTool: 是否需要简历解析工具（仅第一个问题需要）
// 专项面试智能体 不需要传简历的
func NewSpecialQuestionAgent(userId uint, needResumeTool bool) adk.Agent {
	ctx := context.Background()

	// 根据是否需要简历工具构建不同的配置
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

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionAgent",
		Description: "一个专业面试提问的智能体",
		Instruction: `你是一个专业的面试官。根据针对特定技术栈的深度评估，快速生成一个面试问题和对话。

【专项面试】
- 针对特定技术栈的深度评估

难度级别：初级难度、中级难度、高级难度

任务：
1. 根据提示词中指定的技术栈、领域和难度级别进行提问
2. 面试提问阶段：
   - 每个维度的主题只出现一遍
   - 每个主要话题下提出2个深入追问问题

3. 提问策略：
   - 问题应该由浅入深，逐步深入
   - 确保问题与对应评估维度紧密相关
   - 根据难度级别调整问题深度
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
    {"speaker_type": "interviewer", "content": "提问", "display_order": 1}
  ]
}`,

		Model:         chat.CreatOpenAiChatModel(ctx, userId),
		ToolsConfig:   toolsConfig,
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}
	return baseAgent
}
