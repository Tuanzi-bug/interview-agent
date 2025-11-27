package agent

//
//import (
//	"ai-eino-interview-agent/chatApp/chat"
//	"ai-eino-interview-agent/chatApp/prompt"
//	tool2 "ai-eino-interview-agent/chatApp/tool"
//	"context"
//	"fmt"
//	"log"
//
//	"github.com/cloudwego/eino/adk"
//	"github.com/cloudwego/eino/components/tool"
//	"github.com/cloudwego/eino/compose"
//)
//
//func NewResumAnalysisAgent(supervisorName string, userId uint) adk.Agent {
//	ctx := context.Background()
//
//	// 从Redis获取提示词，失败则使用默认模板
//	instruction := prompt.GetPromptInstruction(ctx, "ResumeAnalysisAgent")
//
//	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
//		Name:        "ResumeAnalysisAgent",
//		Description: "一个可以解析简历pdf分析简历的智能体",
//		Instruction: instruction,
//		Model:       chat.CreatOpenAiChatModel(ctx, userId),
//		ToolsConfig: adk.ToolsConfig{
//			ToolsNodeConfig: compose.ToolsNodeConfig{
//				Tools: []tool.BaseTool{tool2.CreatePDFToTextTool()},
//			},
//		},
//		MaxIterations: 10,
//	})
//	if err != nil {
//		log.Fatal(fmt.Errorf("failed to create chatmodel: %w", err))
//	}
//
//	// 增强：完成后自动回调Supervisor
//	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
//		Agent:        baseAgent,
//		ToAgentNames: []string{supervisorName},
//	})
//}
