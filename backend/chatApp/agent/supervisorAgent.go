package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"ai-eino-interview-agent/chatApp/prompt"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
)

// 面试调度Supervisor
func NewInterviewSupervisorAgent(userId uint) adk.Agent {
	ctx := context.Background()
	supervisorName := "NewInterviewSupervisorAgent" // 固定名称，子Agent回调用

	// 步骤1：创建所有增强子Agent（传入Supervisor名称）
	resumeAgent := NewResumAnalysisAgent(supervisorName, userId)
	questionAgent := NewQuestionGeneratorAgent(supervisorName, userId)
	evalAgent := NewAnswerEvalAgent(supervisorName, userId)
	reportAgent := NewInterviewReportAgent(supervisorName, userId)

	// 步骤2：配置Supervisor核心逻辑
	// 从Redis获取提示词，失败则使用默认模板
	instruction := prompt.GetPromptInstruction(ctx, "InterviewSupervisorAgent")

	supervisorConfig := &adk.ChatModelAgentConfig{
		Name:        supervisorName,
		Description: "面试调度中心，分配简历分析、问题生成、回答评估、报告生成任务",
		// 关键：Supervisor的Instruction定义任务分配规则
		Instruction: instruction,
		Model:       chat.CreatOpenAiChatModel(ctx, userId),
	}

	// 步骤3：创建Supervisor Agent
	supervisor, err := adk.NewChatModelAgent(context.Background(), supervisorConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 步骤4：注册子Agent到Supervisor（关键：让Supervisor能找到子Agent）
	registeredSupervisor, err := adk.SetSubAgents(context.Background(), supervisor, []adk.Agent{
		resumeAgent,
		questionAgent,
		evalAgent,
		reportAgent,
	})
	if err != nil {
		log.Fatalf("子Agent注册失败：%v", err)
	}

	return registeredSupervisor
}
