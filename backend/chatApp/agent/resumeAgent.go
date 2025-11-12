package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	tool2 "ai-eino-interview-agent/chatApp/tool"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewResumAnalysisAgent(supervisorName string) adk.Agent {
	ctx := context.Background()

	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ResumAnalysisAgent",
		Description: "一个可以解析简历pdf分析简历的智能体",
		Instruction: `你是一名资深的简历分析专家，负责对用户的简历进行分析,并输出对应的分析结果。

工作流程：
1. 当用户提供简历或相关背景信息时，先使用 "pdf_to_text" 工具提取文本（如需要），进行结构化评估。
2. 从模块完整度、技能匹配度、量化成果、语言表达等角度给出详细反馈，并提供可执行的改进建议；并进行简历评分（0-100分）。`,
		Model: chat.CreatOpenAiChatModel(ctx),
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{tool2.CreatePDFToTextTool()},
			},
		},
		MaxIterations: 10,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create chatmodel: %w", err))
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}
