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

// ComprehensiveInterviewAgentConfig 综合面试智能体配置
type ComprehensiveInterviewAgentConfig struct {
	UserID   uint
	ResumeID uint64
	// 是否是第一次提问（用于控制简历工具调用）
	IsFirstQuestion bool
}

// ComprehensiveInterviewAgent 综合面试智能体
// 适用于校招和社招的综合能力评估
type ComprehensiveInterviewAgent struct {
	agent  adk.Agent
	config *ComprehensiveInterviewAgentConfig
}

// NewComprehensiveInterviewAgent 创建综合面试智能体
func NewComprehensiveInterviewAgent(config *ComprehensiveInterviewAgentConfig) *ComprehensiveInterviewAgent {
	ctx := context.Background()

	// 根据是否是第一次提问来决定是否包含简历工具
	var tools []componenttool.BaseTool
	if config.IsFirstQuestion && config.ResumeID > 0 {
		tools = []componenttool.BaseTool{
			tool2.GetResumeInfoTool(),
		}
	}

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ComprehensiveInterviewAgent",
		Description: "综合面试智能体 - 适用于综合能力评估",
		Instruction: comprehensiveInstruction,
		Model:       chat.CreatOpenAiChatModel(ctx, config.UserID),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		MaxIterations: 5,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create comprehensive interview agent: %w", err))
	}

	return &ComprehensiveInterviewAgent{
		agent:  baseAgent,
		config: config,
	}
}

const comprehensiveInstruction = `你是专业的面试官，进行综合能力评估。

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
func (c *ComprehensiveInterviewAgent) GetAgent() adk.Agent {
	return c.agent
}
