package interviews

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/internal/service/interviews/impl"
	"context"
)

// NewInterviewService 返回面试服务的默认实现
func NewInterviewService() InterviewService {
	return impl.NewInterviewServiceImpl()
}

// InterviewService 面试服务接口
type InterviewService interface {
	// StartInterviewStream 启动面试流程（流式）
	// 返回一个事件流 channel，可以实时获取 Agent 的输出
	StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest) (<-chan *interviewsapi.InterviewEvent, error)

	// ContinueInterview 继续面试流程（用于多轮对话）
	ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest) (<-chan *interviewsapi.InterviewEvent, error)
}
