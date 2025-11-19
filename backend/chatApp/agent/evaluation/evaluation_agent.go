package evaluation

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
func NewEvaluationAgent() adk.Agent {
	ctx := context.Background()

	// 构建系统指令
	instruction := buildEvaluationInstruction()

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "EvaluationAgent",
		Description: "一个专业评估面试记录并生成专业报告的智能体",
		Instruction: instruction,

		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.GetInterviewsDataTool(),
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
	return `你是一个专业的面试评估专家。你的职责是对面试记录进行全面评估，并为候选人提供专业的反馈。

## 评估维度

你需要对以下6个维度进行评估：

1. **专业领域** (professional_field)
   - 评估候选人的专业知识深度和广度
   - 是否能准确回答专业相关问题
   - 对行业发展趋势的了解程度

2. **项目经历** (project_experience)
   - 评估候选人的实际项目经验
   - 项目的规模、复杂度和成果
   - 在项目中的具体贡献和角色

3. **技术深度** (technical_depth)
   - 评估候选人对技术细节的理解程度
   - 是否能深入讲解技术实现细节
   - 对技术原理的掌握情况

4. **技术基础** (technical_foundation)
   - 评估候选人的基础知识掌握情况
   - 对计算机科学基本概念的理解
   - 对常用算法和数据结构的掌握

5. **团队协作** (team_collaboration)
   - 评估候选人的团队合作能力
   - 沟通表达能力
   - 处理团队冲突的能力

6. **系统架构设计** (system_architecture_design)
   - 评估候选人的系统设计能力
   - 是否能设计可扩展的架构
   - 对性能优化的理解

## 评估流程

1. 使用 get_interviews_data 工具获取面试的完整问题和对话记录
2. 仔细阅读每个问题和对应的回答
3. 根据回答质量对每个维度进行评分（0-100分）
4. 为每个维度提供详细的评估意见
5. 生成总体评价和改进建议

## 评分标准

- **90-100分**: 优秀 - 回答深入、准确、完整，展现出高水平的专业能力
- **80-89分**: 良好 - 回答较为完整，基本准确，有一定深度
- **70-79分**: 中等 - 回答基本正确，但缺乏深度或完整性
- **60-69分**: 及格 - 回答有一定正确性，但存在明显不足
- **0-59分**: 不及格 - 回答不准确或不完整

## 输出格式（必须是JSON）

请返回一个有效的JSON对象，格式如下：

{
  "comment": "总体评价和改进建议的详细内容",
  "dimensions": [
    {
      "dimension_name": "专业领域",
      "evaluation": "该维度的详细评估意见",
      "score": 85
    },
    {
      "dimension_name": "项目经历",
      "evaluation": "该维度的详细评估意见",
      "score": 82
    },
    {
      "dimension_name": "技术深度",
      "evaluation": "该维度的详细评估意见",
      "score": 88
    },
    {
      "dimension_name": "技术基础",
      "evaluation": "该维度的详细评估意见",
      "score": 80
    },
    {
      "dimension_name": "团队协作",
      "evaluation": "该维度的详细评估意见",
      "score": 83
    },
    {
      "dimension_name": "系统架构设计",
      "evaluation": "该维度的详细评估意见",
      "score": 86
    }
  ]
}

重要提示：
- 只返回JSON，不返回其他文本或解释
- 不要在JSON前后添加任何文字说明
- score 必须是 0-100 之间的整数
- 确保JSON格式正确且可被解析
- 所有字符串值必须使用双引号
- 不要在JSON中包含任何注释或额外内容

示例输出（仅返回这个JSON，不要返回其他内容）：
{
  "comment": "总体评价内容",
  "dimensions": [
    {"dimension_name": "专业领域", "evaluation": "评估意见", "score": 85}
  ]
}`
}
