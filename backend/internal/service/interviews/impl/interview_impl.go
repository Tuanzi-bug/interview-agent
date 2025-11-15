package impl

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/chatApp/agent/ext"
	"ai-eino-interview-agent/internal/model"
	"context"
	"time"
)

// InterviewServiceImpl 面试服务实现
type InterviewServiceImpl struct{}

// NewInterviewServiceImpl 创建面试服务实例
func NewInterviewServiceImpl() *InterviewServiceImpl {
	return &InterviewServiceImpl{}
}

// StartInterviewStream 启动面试流程（流式）
func (s *InterviewServiceImpl) StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的流式接口
	eventChan, err := ext.StartInterviewStream(ctx, req.Query, maxQuestions)
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
func (s *InterviewServiceImpl) ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error) {
	// 调用 chatApp 中的继续面试接口
	eventChan, err := ext.ContinueInterview(ctx, req.Query, maxQuestions)
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

// SaveInterviewRecord 保存面试记录
func (s *InterviewServiceImpl) SaveInterviewRecord(ctx context.Context, userID uint, title, query string) (uint64, error) {
	record := &model.InterviewRecord{
		UserID:   userID,
		Title:    title,
		Query:    query,
		Status:   "pending",
		Messages: "[]",
	}
	err := model.InterviewRecordDao.CreateInterviewRecord(record)
	if err != nil {
		return 0, err
	}
	return record.ID, nil
}

// UpdateInterviewRecord 更新面试记录（用于保存对话历史和状态）
func (s *InterviewServiceImpl) UpdateInterviewRecord(ctx context.Context, recordID uint64, messages string, status, currentAgent string, lastModifiedAt time.Time) error {
	record := &model.InterviewRecord{
		ID:             recordID,
		Messages:       messages,
		Status:         status,
		CurrentAgent:   currentAgent,
		LastModifiedAt: lastModifiedAt,
	}
	return model.InterviewRecordDao.UpdateInterviewRecord(record)
}

// CompleteInterviewRecord 完成面试记录（保存最终报告和评分）
func (s *InterviewServiceImpl) CompleteInterviewRecord(ctx context.Context, recordID uint64, report string, duration int64, score *float64) error {
	return model.InterviewRecordDao.CompleteInterviewRecord(recordID, report, duration, score)
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
	if event.Status != nil && *event.Status != "" {
		apiEvent.Status = event.Status
	}
	if event.Report != "" {
		apiEvent.Report = &event.Report
	}
	if event.Score != nil {
		apiEvent.Score = event.Score
	}
	if event.Duration > 0 {
		apiEvent.Duration = &event.Duration
	}
	if event.Feedback != "" {
		apiEvent.Feedback = &event.Feedback
	}
	if event.Messages != "" {
		apiEvent.Messages = &event.Messages
	}

	return apiEvent
}
