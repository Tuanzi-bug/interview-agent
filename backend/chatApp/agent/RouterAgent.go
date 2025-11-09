package agent

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"log"

	"github.com/cloudwego/eino/adk"
)

func NewRouterAgent() adk.Agent {
	ctx := context.Background()
	a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "RouterAgent",
		Description: "智能任务分发器，将用户请求转交给最合适的专家助手。",
		Instruction: `你是一个智能任务路由器。请分析用户请求，并将其委派给最擅长的专家助手处理；
		1. 如果用户请求是关于简历分析的，则委派给ResumAnalysisAgent处理；
		2. 如果用户请求根据用户的简历去提问问题的，则委派给QuestionGeneratorAgent处理；
		3. 如果用户的请求是根据用户的简历和所有的问题和答案去评估用户的面试表现，则委派给InterviewAgent处理;
		4. 如果没有合适的助手，则直接告知用户无法处理'

		约束: 将子Agent的返回结果原封不动的返回给用户
		`,
		Model: chat.CreatOpenAiChatModel(ctx),
	})
	if err != nil {
		log.Fatal(err)
	}
	return a
}
