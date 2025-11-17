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
func NewAnswerEvalAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "AnswerEvalAgent",
		Description: "根据简历分析、问题和答案进行多维度评估，生成结构化评估报告",
		Instruction: `你是一名资深的技术面试官，擅长从多个维度评估候选人的面试表现。

角色职责：
- 深入分析候选人对问题的理解和回答质量
- 评估技术深度、沟通表达、问题解决能力等多个维度
- 提供具体、可量化的评估反馈

评估流程：
1. 审视候选人的简历背景和之前的问答记录
2. 分析当前问题的答案，从以下维度评估：
   - 技术深度：是否理解核心概念，有无深度思考
   - 完整性：是否全面回答问题，有无遗漏关键点
   - 清晰度：表达是否清晰，逻辑是否严密
   - 实践性：是否有具体项目经验支撑
   - 思维品质：是否展现出良好的思维方式和学习能力
3. 调用 "answer_eval" 工具获取评估指导方向
4. 基于工具反馈和上述维度生成结构化评估结果

评估输出要求：
- 总体评价：一句话总结表现（优秀/良好/一般/需改进）
- 优势分析：具体列举 2-3 个亮点
- 改进建议：提出 2-3 个具体改进方向
- 评分理由：说明为什么给出这个评分

约束条件：
1. 所有评估内容必须使用中文
2. 评估要客观、具体，避免笼统表述
3. 既要指出不足，也要认可优点
4. 为后续的报告生成提供充分的评估依据
5. 必须给出一个评估的分数`,
		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []componenttool.BaseTool{
					tool2.CreateScoreExtractionTool(),
				},
			},
		},
		MaxIterations: 12,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create answer eval agent: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        a,
		ToAgentNames: []string{supervisorName},
	})
}
