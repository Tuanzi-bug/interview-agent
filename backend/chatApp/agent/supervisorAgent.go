package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"github.com/cloudwego/eino/adk"
	"log"
)

// 面试调度Supervisor
func NewInterviewSupervisorAgent() adk.Agent {
	supervisorName := "NewInterviewSupervisorAgent" // 固定名称，子Agent回调用

	// 步骤1：创建所有增强子Agent（传入Supervisor名称）
	resumeAgent := NewResumAnalysisAgent(supervisorName)
	questionAgent := NewQuestionGeneratorAgent(supervisorName)
	evalAgent := NewAnswerEvalAgent(supervisorName)
	reportAgent := NewInterviewReportAgent(supervisorName)

	// 步骤2：配置Supervisor核心逻辑
	supervisorConfig := &adk.ChatModelAgentConfig{
		Name:        supervisorName,
		Description: "面试调度中心，分配简历分析、问题生成、回答评估、报告生成任务",
		// 关键：Supervisor的Instruction定义任务分配规则
		Instruction: `你是面试调度专家，遵循以下流程：
1. 如果用户提供简历（文本或PDF路径），先转让给ResumeAnalysisAgent解析
2. 如果用户要求开始面试，直接转让给QuestionGeneratorAgent生成技术问题
3. 简历分析后：转让给QuestionGeneratorAgent生成技术问题；
4. 用户提供回答后：转让给AnswerEvalAgent评估；
5. 评估后：转让给InterviewReportAgent生成最终报告；
6. 报告生成后：直接输出报告，结束流程；
7. 每步完成后，等待用户下一步输入（如用户提供回答），再进行下一轮分配。`,
		Model: chat.CreatOpenAiChatModel(context.Background()),
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
