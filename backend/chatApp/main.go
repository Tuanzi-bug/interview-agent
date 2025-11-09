package main

import (
	"ai-eino-interview-agent/chatApp/agent"
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func main1() {
	ctx := context.Background()

	//使用message模版
	fmt.Printf("===create messages===\n")
	message := chat.MessagesTemplate()
	//fmt.Printf("messages: %+v\n\n", chat.MessagesTemplate())

	//创建llm
	fmt.Printf("===create llm===\n")
	model := chat.CreatOpenAiChatModel(ctx)
	//log.Printf("create llm success\n\n")

	//使用llm生成Steam流式回复
	fmt.Printf("===llm stream ===\n")
	streamResult := chat.Stream(ctx, model, message)
	chat.ReportSteam(streamResult)
}

func mockResumeAnalysisJSON() string {
	return `{
  "candidate_name": "王小明",
  "summary": "后端开发工程师，擅长 Go/Kratos 微服务与分布式系统架构。",
  "highlights": [
    "在深圳愿景互娱负责抽奖系统重构，利用中间表归并方案与容灾限流保障高并发稳定。",
    "主导智能烹饪知识检索系统，基于 Eino + Milvus + Neo4j 实现多模态 RAG，召回率提升 40%。"
  ],
  "risk_points": [
    "教育背景缺少量化指标（GPA、核心课程等）。",
    "部分项目成果缺乏可验证的指标描述。"
  ],
  "skills": {
    "proficient": ["Go", "Kratos", "gRPC", "MySQL"],
    "familiar": ["Milvus", "Neo4j", "Docker Compose", "分布式缓存"],
    "learning": ["全链路追踪优化", "图数据库检索策略"]
  }
}`
}

func mockHistoryQAJSON() string {
	return `[
  {
    "question": "请做一个 1 分钟的自我介绍，突出你的核心技术优势。",
    "answer": "我是王小明，目前在深圳愿景互娱担任后端开发实习生，主要负责 Go/Kratos 微服务、抽奖系统重构和智能 Agent 项目。"
  },
  {
    "question": "介绍你在深圳愿景互娱参与的抽奖系统重构项目。",
    "answer": "我通过中间表方案将独立抽奖系统接入活动系统，实现卡池热更，并配合压力测试将 RPC p99 延迟从 600ms 降到 400ms。"
  },
  {
    "question": "在智能烹饪知识检索系统中，你是如何设计检索策略的？",
    "answer": "我结合向量检索、图结构检索和 BM25，采用 RRF 融合以及 BERT 重排序，使 Top-3 准确率提升 25%。"
  }
]
`
}

func buildQuestionGeneratorQuery() []adk.Message {
	questionIndex := 3
	mockResumeAnalysis := mockResumeAnalysisJSON()
	mockHistoryQA := mockHistoryQAJSON()

	mockPayload := fmt.Sprintf(`{
  "intent": "mock_interview",
  "question_index": %d,
  "resume_analysis": %s,
  "history_qa": %s,
  "meta": {"stage": "技术深挖"}
}`, questionIndex, mockResumeAnalysis, mockHistoryQA)

	mockQuery := "我想继续进行程序员模拟面试，请根据以下上下文生成一个问题：\n" + mockPayload

	mockMessages := []adk.Message{
		schema.UserMessage(mockQuery),
	}
	return mockMessages
}

func buildResumeAnalysisQuery() []adk.Message {
	//这个需要你去提供你的简历
	query := "帮我解析C:\\Users\\akf\\Desktop\\akf.pdf这个PDF文件 开始模拟面试,进行5轮问题回答,"

	mockMessages := []adk.Message{
		schema.UserMessage(query),
	}
	return mockMessages
}

func buildAnswerEvalQuery() []adk.Message {
	mockResumeAnalysis := mockResumeAnalysisJSON()
	mockHistoryQA := mockHistoryQAJSON()

	mockPayload := fmt.Sprintf(`{
  "intent": "answer_eval",
  "resume_analysis": %s,
  "history_qa": %s,
  "meta": {"stage": "技术深挖"}
}`, mockResumeAnalysis, mockHistoryQA)

	query := "请根据以下的简历分析与历史问答记录，综合评估候选人在本次模拟面试中的表现，给出 0-100 的分数，并写出评语：\n" + mockPayload

	mockMessages := []adk.Message{
		schema.UserMessage(query),
	}
	return mockMessages
}

func main() {

	ctx := context.Background()
	runner := agent.NewInterviewProcessAgent()

	mockMessages := buildResumeAnalysisQuery()

	iter := runner.Run(ctx, mockMessages)
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(event.Err)
		}
		if event.Action != nil {
			log.Printf("\nAgent[%s]: transfer to %+v\n\n======\n", event.AgentName, event.Action.TransferToAgent.DestAgentName)
		} else {
			log.Printf("\nAgent[%s]:\n%+v\n\n======\n", event.AgentName, event.Output.MessageOutput.Message)
		}

		if event.Output != nil && event.Output.MessageOutput.Message.Content != "" {
			lastMessage, _, err := adk.GetMessage(event)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("lastMessage: content=%+v role=%s\n", lastMessage.Content, lastMessage.Role)
		}
	}
}
