package rht

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

// SpecializedInterviewAgentConfig 专项面试智能体配置
type SpecializedInterviewAgentConfig struct {
	UserID   uint
	ResumeID uint64
	// 是否是第一次提问（用于控制简历工具调用）
	IsFirstQuestion bool
}

// SpecializedInterviewAgent 专项面试智能体
// 针对特定技术栈的深度评估（如Java、Golang等）
type SpecializedInterviewAgent struct {
	agent  adk.Agent
	config *SpecializedInterviewAgentConfig
}

// NewSpecializedInterviewAgent 创建专项面试智能体
func NewSpecializedInterviewAgent(config *SpecializedInterviewAgentConfig) *SpecializedInterviewAgent {
	ctx := context.Background()

	// 根据是否是第一次提问来决定是否包含简历工具
	var tools []componenttool.BaseTool
	if config.IsFirstQuestion && config.ResumeID > 0 {
		tools = []componenttool.BaseTool{
			tool2.GetResumeInfoTool(),
		}
	}

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "SpecializedInterviewAgent",
		Description: "专项面试智能体 - 针对特定技术栈的深度评估",
		Instruction: specializedInstruction,
		Model:       chat.CreatOpenAiChatModel(ctx, config.UserID),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create specialized interview agent: %w", err))
	}

	return &SpecializedInterviewAgent{
		agent:  baseAgent,
		config: config,
	}
}

const specializedInstruction = `你是专业的技术面试官，进行深度的技术栈评估。

核心要求：
1. 根据提示词中的难度级别进行提问
2. 如果需要简历信息，调用get_resume_info工具（仅第一次）
3. 只返回JSON格式，不返回其他文本
4. 只生成面试官的提问，不生成用户回答
5. 自由选择合适的评估维度进行提问

必须返回的JSON格式：
{
  "questions": [{"question_text": "问题内容", "eval_dimension": "维度", "order": 1}],
  "dialogues": [{"speaker_type": "interviewer", "content": "提问", "display_order": 1}]
}`

// GetAgent 获取底层的Agent
func (s *SpecializedInterviewAgent) GetAgent() adk.Agent {
	return s.agent
}
