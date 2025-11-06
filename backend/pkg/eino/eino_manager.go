package eino

import (
	"context"
	"log"

	"ai-eino-interview-agent/internal/config"
)

// MockEinoAgent 模拟的Eino代理实例
type MockEinoAgent struct {
	Config config.EinoConfig
}

// EinoAgent 全局Eino代理实例
var EinoAgent *MockEinoAgent

// InitEino 初始化Eino框架
func InitEino(einoConfig config.EinoConfig) error {
	// 创建模拟的Eino代理
	EinoAgent = &MockEinoAgent{
		Config: einoConfig,
	}

	// 更新全局配置
	config.Global.Eino = einoConfig

	log.Println("Eino框架模拟实现初始化成功")
	return nil
}

// GetEinoAgent 获取Eino代理实例
func GetEinoAgent() *MockEinoAgent {
	return EinoAgent
}

// Chat 模拟聊天功能
func (a *MockEinoAgent) Chat(ctx context.Context, prompt string) (string, error) {
	// 这是一个模拟实现，返回一些默认响应
	// 实际项目中应该集成真实的AI服务
	log.Printf("Eino代理收到提示：%s\n", prompt)
	return "这是模拟的AI响应", nil
}

// GenerateInterviewQuestions 生成面试问题
func GenerateInterviewQuestions(ctx context.Context, jobTitle, difficulty, interviewType, resumeContent string) ([]string, error) {
	// 模拟生成面试问题
	questions := []string{
		"请介绍一下你在" + jobTitle + "方面的经验？",
		"你如何处理" + difficulty + "级别的技术挑战？",
		"在" + interviewType + "环境中，你最擅长的技能是什么？",
		"请分享一个你印象深刻的项目经验。",
	}
	return questions, nil
}

// AnalyzeResume 分析简历
func AnalyzeResume(ctx context.Context, resumeContent string) (string, error) {
	// 模拟简历分析
	return "这是模拟的简历分析结果。简历内容分析完毕。", nil
}

// EvaluateAnswer 评估面试回答
func EvaluateAnswer(ctx context.Context, question, answer, jobTitle string) (string, float64, error) {
	// 模拟回答评估
	return "回答评估结果", 85.5, nil
}
