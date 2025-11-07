package builders

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

// OpenAIBuilder 是 OpenAI 协议的构建器实现。
type OpenAIBuilder struct{}

// Build 构建 OpenAI 聊天模型实例。
func (b *OpenAIBuilder) Build(ctx context.Context, config *factory.Config) (factory.ChatModel, error) {
	// 转换配置格式
	cfg := &openai.ChatModelConfig{
		APIKey:      config.APIKey,
		BaseURL:     config.BaseURL,
		Model:       config.Model,
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		TopP:        config.TopP,
	}

	// 创建 eino 的 OpenAI 聊天模型
	chatModel, err := openai.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// 使用共享的适配器包装
	return newEinoChatModelAdapter(chatModel), nil
}

// NewOpenAIBuilder 创建一个新的 OpenAI 构建器实例。
func NewOpenAIBuilder() factory.Builder {
	return &OpenAIBuilder{}
}
