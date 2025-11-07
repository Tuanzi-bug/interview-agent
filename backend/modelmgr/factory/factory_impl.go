package factory

import "context"

// defaultFactory 是 Factory 接口的默认实现
type defaultFactory struct {
	builders map[Protocol]Builder
}

// NewFactory 创建一个新的工厂实例
func NewFactory(builders map[Protocol]Builder) Factory {
	return &defaultFactory{
		builders: builders,
	}
}

// Create 根据协议和配置创建聊天模型实例
func (f *defaultFactory) Create(ctx context.Context, protocol Protocol, config *Config) (ChatModel, error) {
	builder, ok := f.builders[protocol]
	if !ok {
		return nil, ErrProtocolNotSupported
	}
	return builder.Build(ctx, config)
}

// Support 检查是否支持指定的协议
func (f *defaultFactory) Support(protocol Protocol) bool {
	_, ok := f.builders[protocol]
	return ok
}
