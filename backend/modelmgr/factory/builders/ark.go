package builders

import (
	"ai-eino-interview-agent/modelmgr/factory"
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

// ArkBuilder 是火山引擎 Ark 协议的构建器实现。
// 使用的底层 SDK：
//   - github.com/cloudwego/eino-ext/components/model/ark
//
// 参考文档：
//   - https://www.volcengine.com/docs/82379
type ArkBuilder struct{}

// Build 构建 Ark 聊天模型实例。
func (b *ArkBuilder) Build(ctx context.Context, config *factory.Config) (factory.ChatModel, error) {
	// 验证必需参数
	if config.APIKey == "" {
		return nil, fmt.Errorf("ark builder: api_key is required")
	}
	if config.Model == "" {
		return nil, fmt.Errorf("ark builder: model (endpoint_id) is required")
	}

	// 转换配置格式
	// Ark SDK 的配置结构与 OpenAI 类似，但使用火山引擎的服务端点
	cfg := &ark.ChatModelConfig{
		APIKey:      config.APIKey,
		Model:       config.Model, // Ark 使用 Endpoint ID 作为 Model
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
		TopP:        config.TopP,
	}

	// 如果提供了自定义 BaseURL，则使用自定义端点
	if config.BaseURL != "" {
		cfg.BaseURL = config.BaseURL
	}

	// 创建 eino 的 Ark 聊天模型
	chatModel, err := ark.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("ark builder: failed to create chat model: %w", err)
	}

	// 使用共享的适配器包装
	return newEinoChatModelAdapter(chatModel), nil
}

// NewArkBuilder 创建一个新的 Ark 构建器实例。
func NewArkBuilder() factory.Builder {
	return &ArkBuilder{}
}
