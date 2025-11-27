package demo

import (
	demoapi "ai-eino-interview-agent/api/model/demo"
	"ai-eino-interview-agent/internal/service/demo/impl"
	"golang.org/x/net/context"
)

func NewDemoService() DemoService {
	return impl.NewDemoService()
}

type DemoService interface {
	Create(ctx context.Context, req demoapi.CreateDemoRequest) (*demoapi.CreateDemoResponse, error)
}
