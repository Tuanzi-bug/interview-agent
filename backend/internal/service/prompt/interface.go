package prompt

import (
	"ai-eino-interview-agent/internal/service/prompt/impl"
	"context"
)

func NewPromptManager() PromptService {
	return impl.NewPromptService()
}

// PromptService 提示词服务接口
type PromptService interface {
	// GetPromptTemplate 获取提示词模板
	GetPromptTemplate(ctx context.Context, agentName string) (string, error)
}
