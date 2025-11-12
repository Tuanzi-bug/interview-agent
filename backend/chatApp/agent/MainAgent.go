package agent

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
)

func NewMainAgent() *adk.Runner {
	ResumAnalysisAgent := NewResumAnalysisAgent()
	// QuestionGeneratorAgent := NewQuestionGeneratorAgent()
	RouterAgent := NewRouterAgent()
	// AnswerEvalAgent := NewAnswerEvalAgent()

	ctx := context.Background()
	// a, err := adk.SetSubAgents(ctx, RouterAgent, []adk.Agent{ResumAnalysisAgent, QuestionGeneratorAgent, AnswerEvalAgent})
	a, err := adk.SetSubAgents(ctx, RouterAgent, []adk.Agent{ResumAnalysisAgent})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	log.Printf("set sub agents success\n")
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: a,
	})
	return runner
}
