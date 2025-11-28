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
func NewQuestionAgent(userId uint, needResumeTool bool) adk.Agent {
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

	model, err := chat.CreatOpenAiChatModel(ctx, userId)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create OpenAI chat model: %w", err))
	}

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "QuestionAgent",
		Description: "一个专业面试提问的智能体",
		Instruction: `你是一个专业的面试官。根据提供的简历或信息，快速生成一个面试问题和对话。

支持两种面试类型：

【综合面试】
- 适用于校招和社招的综合能力评估
- 核心评估维度：
  * professional_field (专业领域)
  * project_experience (项目经历)
  * technical_depth (技术深度)
  * technical_foundation (技术基础)
  * team_collaboration (团队协作)
  * system_architecture_design (系统架构设计)

【专项面试】
- 针对特定技术栈的深度评估（如Java、Golang等）
- 核心评估维度：
  * basic_knowledge_mastery (基础知识掌握)
  * working_principle_practical_experience (工作原理与实践经验)
  * advanced_features_application (高级特性应用)
  * problem_troubleshooting_skills (问题排查能力)
  * architecture_design_thinking (架构设计思维)
  * performance_optimization_ability (性能优化能力)

难度级别：初级难度、中级难度、高级难度

任务：
1. 如果提供了resume_id字段的值，使用get_resume_info工具解析简历获取简历内容
2. 根据提示词中指定的面试类型、领域和难度级别进行提问
3. 面试提问阶段：
   - 根据候选人的背景和技能，围绕指定的核心评估维度进行提问
   - 每个维度的主题只出现一遍
   - 每个主要话题下提出2个深入追问问题

4. 提问策略：
   - 问题应该由浅入深，逐步深入
   - 结合候选人提供的背景信息
   - 确保问题与对应评估维度紧密相关
   - 根据难度级别调整问题深度
   - 保持专业性和针对性
5. 只返回JSON，不要返回其他文本

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

		Model:         model,
		ToolsConfig:   toolsConfig,
		MaxIterations: 20,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
	}
	return baseAgent
}
