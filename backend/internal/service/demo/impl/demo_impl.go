package impl

import (
	demoapi "ai-eino-interview-agent/api/model/demo"
	"ai-eino-interview-agent/internal/model"
	"context"
	"net/http"
	"time"
)

type DemoServer struct {
	httpClient *http.Client
}

func NewDemoService() *DemoServer {
	return &DemoServer{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DemoServer) Create(_ context.Context, req demoapi.CreateDemoRequest) (*demoapi.CreateDemoResponse, error) {
	userRecord := &model.Demo{
		Name: req.Name,
	}

	if err := model.DemoDao.Create(userRecord); err != nil {
		return nil, err
	}

	return demoapi.NewCreateDemoResponse(), nil
}
