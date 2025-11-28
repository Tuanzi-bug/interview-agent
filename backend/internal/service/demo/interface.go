package demo

import (
	demoapi "ai-eino-interview-agent/api/model/demo"
	"ai-eino-interview-agent/internal/service/demo/impl"
	"context"
)

// NewDemoManager 返回用户管理接口的默认实现
func NewDemoManager() DemoManager {
	return impl.NewDemoServer()
}

type DemoManager interface {
	Create(ctx context.Context, req demoapi.CreateDemoModelRequest) (*demoapi.CreateDemoModelResponse, error)
}
