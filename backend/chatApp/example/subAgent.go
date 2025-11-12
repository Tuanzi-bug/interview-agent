package example

import (
	"ai-eino-interview-agent/chatApp/chat"
	"context"
	"github.com/cloudwego/eino/adk"
	"log"
)

// 1. 简历分析Agent（增强版：完成后回调Supervisor）
func NewResumeAnalysisAgent(supervisorName string) adk.Agent {
	// 基础Agent配置
	baseAgent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "ResumeAnalysisAgent",
		Description: "解析简历，提取核心技能、项目经验、评分",
		Instruction: `你是简历分析专家，解析用户提供的简历，输出核心技能（3-5个）、项目经验摘要、0-100分评分，格式简洁。`,
		Model:       chat.CreatOpenAiChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 增强：完成后自动回调Supervisor
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}

// 2. 技术问题生成Agent（增强版）
func NewTechQuestionAgent(supervisorName string) adk.Agent {
	baseAgent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "TechQuestionAgent",
		Description: "基于简历核心技能，生成3个技术面试问题",
		Instruction: `你是Go后端技术面试官，基于简历中的核心技能，生成3个技术问题（含基础、进阶、项目相关），格式清晰。`,
		Model:       chat.CreatOpenAiChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}

// 3. 回答评估Agent（增强版）
func NewAnswerEvalAgent(supervisorName string) adk.Agent {
	baseAgent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "AnswerEvalAgent",
		Description: "评估用户的面试回答，给出评分和改进建议",
		Instruction: `你是技术面试评估专家，基于问题和用户回答，给出0-10分评分和1-2条改进建议，格式简洁。`,
		Model:       chat.CreatOpenAiChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}

// 4. 面试报告Agent（增强版）
func NewInterviewReportAgent(supervisorName string) adk.Agent {
	baseAgent, err := adk.NewChatModelAgent(context.Background(), &adk.ChatModelAgentConfig{
		Name:        "InterviewReportAgent",
		Description: "汇总面试全程结果，生成结构化面试报告",
		Instruction: `你是面试报告专家，汇总简历分析、问题、回答评估，生成包含候选人基本信息、技能匹配度、回答评分、录用建议的报告。`,
		Model:       chat.CreatOpenAiChatModel(context.Background()),
	})
	if err != nil {
		log.Fatal(err)
	}
	return adk.AgentWithDeterministicTransferTo(context.Background(), &adk.DeterministicTransferConfig{
		Agent:        baseAgent,
		ToAgentNames: []string{supervisorName},
	})
}
