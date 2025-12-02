package rht

import (
	"fmt"
)

// InterviewAgentType 面试智能体类型
type InterviewAgentType string

const (
	// SpecializedType 专项面试
	SpecializedType InterviewAgentType = "specialized"
	// ComprehensiveType 综合面试
	ComprehensiveType InterviewAgentType = "comprehensive"
)

// InterviewAgentFactory 面试智能体工厂
type InterviewAgentFactory struct{}

// NewInterviewAgentFactory 创建工厂实例
func NewInterviewAgentFactory() *InterviewAgentFactory {
	return &InterviewAgentFactory{}
}

// CreateAgent 根据类型创建对应的智能体
// agentType: "specialized" 或 "comprehensive"
// userID: 用户ID
// resumeID: 简历ID
// isFirstQuestion: 是否是第一次提问
func (f *InterviewAgentFactory) CreateAgent(agentType InterviewAgentType, userID uint, resumeID uint64, isFirstQuestion bool) (interface{}, error) {
	switch agentType {
	case SpecializedType:
		config := &SpecializedInterviewAgentConfig{
			UserID:          userID,
			ResumeID:        resumeID,
			IsFirstQuestion: isFirstQuestion,
		}
		return NewSpecializedInterviewAgent(config), nil

	case ComprehensiveType:
		config := &ComprehensiveInterviewAgentConfig{
			UserID:          userID,
			ResumeID:        resumeID,
			IsFirstQuestion: isFirstQuestion,
		}
		return NewComprehensiveInterviewAgent(config), nil

	default:
		return nil, fmt.Errorf("unknown interview agent type: %s", agentType)
	}
}

// CreateSpecializedAgent 创建专项面试智能体
func (f *InterviewAgentFactory) CreateSpecializedAgent(userID uint, resumeID uint64, isFirstQuestion bool) *SpecializedInterviewAgent {
	config := &SpecializedInterviewAgentConfig{
		UserID:          userID,
		ResumeID:        resumeID,
		IsFirstQuestion: isFirstQuestion,
	}
	return NewSpecializedInterviewAgent(config)
}

// CreateComprehensiveAgent 创建综合面试智能体
func (f *InterviewAgentFactory) CreateComprehensiveAgent(userID uint, resumeID uint64, isFirstQuestion bool) *ComprehensiveInterviewAgent {
	config := &ComprehensiveInterviewAgentConfig{
		UserID:          userID,
		ResumeID:        resumeID,
		IsFirstQuestion: isFirstQuestion,
	}
	return NewComprehensiveInterviewAgent(config)
}
