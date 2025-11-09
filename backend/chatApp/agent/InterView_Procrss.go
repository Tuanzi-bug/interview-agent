package agent

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino/adk"
	"log"
)

func NewInterviewProcessAgent() *adk.Runner {
	ctx := context.Background()

	ResumeAnalysisAgent := NewResumAnalysisAgent()
	QuestionGeneratorAgent := NewQuestionGeneratorAgent()
	AnswerEvalAgent := NewAnswerEvalAgent()
	InterviewReportAgent := NewInterviewReportAgent()

	loopAgent, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
		Name:        "InterviewProcessAgent",
		Description: "根据最大迭代次数 生成问题 → 向用户获取回答 → 评估回答 ",

		SubAgents:     []adk.Agent{QuestionGeneratorAgent, AnswerEvalAgent},
		MaxIterations: 10,
	})

	sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
		Name:        "InterviewProcessAgent",
		Description: " 模拟面试流程: 分析简历 -> 生成问题 -> 向用户获取回答 -> 评估回答 -> 生成面试报告",
		SubAgents:   []adk.Agent{ResumeAnalysisAgent, loopAgent, InterviewReportAgent},
	})

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	log.Printf("set sub agents success\n")
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: sequentialAgent,
	})

	return runner
}
