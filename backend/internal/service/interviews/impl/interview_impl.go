package impl

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/ext"
	"context"
)

// InterviewServiceImpl 面试服务实现
type InterviewServiceImpl struct{}

// NewInterviewServiceImpl 创建面试服务实例
func NewInterviewServiceImpl() *InterviewServiceImpl {
	return &InterviewServiceImpl{}
}

// StartInterviewStream 启动面试流程（流式）
func (s *InterviewServiceImpl) StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的流式接口
	eventChan, err := ext.StartInterviewStream(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
	go func() {
		defer close(apiEventChan)
		for event := range eventChan {
			apiEvent := convertToAPIEvent(event)
			select {
			case apiEventChan <- apiEvent:
			case <-ctx.Done():
				return
			}
		}
	}()

	return apiEventChan, nil
}

// ContinueInterview 继续面试流程（用于多轮对话）
func (s *InterviewServiceImpl) ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的继续面试接口
	eventChan, err := ext.ContinueInterview(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	// 转换事件类型：从 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
	apiEventChan := make(chan *interviewsapi.InterviewEvent, 10)
	go func() {
		defer close(apiEventChan)
		for event := range eventChan {
			apiEvent := convertToAPIEvent(event)
			select {
			case apiEventChan <- apiEvent:
			case <-ctx.Done():
				return
			}
		}
	}()

	return apiEventChan, nil
}

// convertToAPIEvent 将 ext.InterviewEvent 转换为 interviewsapi.InterviewEvent
func convertToAPIEvent(event *ext.InterviewEvent) *interviewsapi.InterviewEvent {
	apiEvent := interviewsapi.NewInterviewEvent()
	apiEvent.Type = event.Type

	if event.AgentName != "" {
		apiEvent.AgentName = &event.AgentName
	}
	if event.Message != "" {
		apiEvent.Message = &event.Message
	}
	if event.TransferTo != "" {
		apiEvent.TransferTo = &event.TransferTo
	}
	if event.Error != "" {
		apiEvent.Error = &event.Error
	}
	if event.Status != "" {
		apiEvent.Status = &event.Status
	}

	return apiEvent
}
