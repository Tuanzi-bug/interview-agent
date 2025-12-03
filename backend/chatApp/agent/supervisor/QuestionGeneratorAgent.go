package supervisor

//
//import (
//	"ai-eino-interview-agent/chatApp/agent/session"
//	"ai-eino-interview-agent/chatApp/chat"
//	"ai-eino-interview-agent/chatApp/prompt"
//	tool2 "ai-eino-interview-agent/chatApp/tool"
//	"context"
//	"fmt"
//	"log"
//
//	"github.com/cloudwego/eino/adk"
//	componenttool "github.com/cloudwego/eino/components/tool"
//	"github.com/cloudwego/eino/compose"
//)
//
//// NewQuestionGeneratorAgent 创建问题生成智能体
//// 支持通过会话值传递上下文信息，优化性能
//func NewQuestionGeneratorAgent(supervisorName string, userId uint) adk.Agent {
//	ctx := context.Background()
//
//	// 从Redis获取提示词，失败则使用默认模板
//	instruction := prompt.GetPromptInstruction(ctx, "QuestionGeneratorAgent")
//
//	baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
//		Name:        "QuestionGeneratorAgent",
//		Description: "根据前面的简历去分析结果去提问问题",
//		Instruction: instruction,
//		Model:       chat.CreatOpenAiChatModel(ctx, userId),
//		ToolsConfig: adk.ToolsConfig{
//			ToolsNodeConfig: compose.ToolsNodeConfig{
//				Tools: []componenttool.BaseTool{
//					//tool2.CreateGenQuestionTool(),
//					tool2.NewAskForInputTool(),
//					//tool2.NewExitTool(),
//				},
//			},
//		},
//		MaxIterations: 20,
//	})
//	if err != nil {
//		log.Fatal(fmt.Errorf("failed to create question generator agent: %w", err))
//	}
//
//	// 增强：完成后自动回调Supervisor
//	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
//		Agent:        baseAgent,
//		ToAgentNames: []string{supervisorName},
//	})
//}
//
//// InitializeQuestionGeneratorContext 初始化问题生成器的会话值
//// 在 Agent 执行前调用，用于存储简历、配置等信息
//func InitializeQuestionGeneratorContext(ctx context.Context, resumeContent string,
//	interviewType, domain, difficulty string) error {
//	scm := session.NewSessionContextManager(ctx)
//
//	// 存储简历内容
//	if err := scm.SetResumeContent(resumeContent); err != nil {
//		return fmt.Errorf("failed to set resume content: %w", err)
//	}
//
//	// 存储面试配置
//	if err := scm.SetInterviewConfig(interviewType, domain, difficulty); err != nil {
//		return fmt.Errorf("failed to set interview config: %w", err)
//	}
//
//	// 初始化空的对话历史
//	if err := scm.ClearConversationHistory(); err != nil {
//		return fmt.Errorf("failed to clear conversation history: %w", err)
//	}
//
//	return nil
//}
