package main

import (
	"ai-eino-interview-agent/chatApp/agent"
	"fmt"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"golang.org/x/net/context"
	"log"
)

func main() {
	ctx := context.Background()

	// 1. 创建Supervisor Agent（包含所有子Agent）
	interviewSupervisor := agent.NewInterviewSupervisorAgent()

	// 2. 创建Runner
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: interviewSupervisor,
	})

	query := buildResumeAnalysisQuery()

	// 3. 启动面试流程（用户输入简历）
	log.Println("====== 面试流程启动 ======")
	iter := runner.Run(ctx, query)

	// 4. 处理事件流
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatalf("面试流程出错：%v", event.Err)
		}

		// 打印转让事件
		if event.Action != nil && event.Action.TransferToAgent != nil {
			log.Printf("\n🔄 调度中心转让任务给：%s", event.Action.TransferToAgent.DestAgentName)
			continue
		}

		// 打印Agent输出
		if event.Output != nil && event.Output.MessageOutput != nil {
			log.Printf("\n📢 %s 输出：\n%s", event.AgentName, event.Output.MessageOutput.Message.Content)
			if event.AgentName == "InterviewReportAgent" {
				log.Println("\n====== 面试流程结束 ======")
			}
		}
	}
	//ctx := context.Background()
	//runner := agent.NewInterviewProcessAgent()
	//
	//mockMessages := buildResumeAnalysisQuery()
	//
	//iter := runner.Run(ctx, mockMessages)
	//for {
	//	event, ok := iter.Next()
	//	if !ok {
	//		break
	//	}
	//	if event.Err != nil {
	//		log.Fatal(event.Err)
	//	}
	//	if event.Action != nil {
	//		log.Printf("\nAgent[%s]: transfer to %+v\n\n======\n", event.AgentName, event.Action.TransferToAgent.DestAgentName)
	//	} else {
	//		log.Printf("\nAgent[%s]:\n%+v\n\n======\n", event.AgentName, event.Output.MessageOutput.Message)
	//	}
	//
	//	if event.Output != nil && event.Output.MessageOutput.Message.Content != "" {
	//		lastMessage, _, err := adk.GetMessage(event)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		fmt.Printf("lastMessage: content=%+v role=%s\n", lastMessage.Content, lastMessage.Role)
	//	}
	//}
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
	//query := "帮我解析 C:\\Users\\LittleBear\\Desktop\\GoTest.pdf 这个PDF文件 开始模拟面试,进行5轮问题回答,"

	query := "我是3年go开发经验，直接开始模拟面试"

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
