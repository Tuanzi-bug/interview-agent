package interviews

import (
	interviewsapi "ai-eino-interview-agent/api/model/interviews"
	"ai-eino-interview-agent/internal/service/interviews/impl"
	"context"
)

// NewInterviewService 初始化面试服务的实例
func NewInterviewService() InterviewService {
	return impl.NewInterviewServiceImpl()
}

// InterviewService 面试服务接口
type InterviewService interface {
	// StartInterviewStream 启动面试流程（流式）
	// 返回一个事件流 channel，可以实时获取 Agent 的输出
	// maxQuestions: 最大答题次数，0 表示无限制
	StartInterviewStream(ctx context.Context, req *interviewsapi.StartInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error)

	// ContinueInterview 继续面试流程（用于多轮对话）
	// maxQuestions: 最大答题次数，0 表示无限制
	ContinueInterview(ctx context.Context, req *interviewsapi.ContinueInterviewRequest, maxQuestions int) (<-chan *interviewsapi.InterviewEvent, error)

	// SaveInterviewRecord 保存面试记录
	SaveInterviewRecord(ctx context.Context, userID uint, title, query string) (uint64, error)

	// UpdateInterviewRecord 更新面试记录（用于保存对话历史和状态）
	// lastModifiedAt: 上次修改的时间戳，用于并发控制检查
	UpdateInterviewRecord(ctx context.Context, recordID uint64, messages string, status, currentAgent string) error

	// CompleteInterviewRecord 完成面试记录（保存最终报告和评分）
	CompleteInterviewRecord(ctx context.Context, recordID uint64, report string, duration int64, score *float64) error

	ListInterviewRecords(ctx context.Context, userID uint, page, pageSize int) ([]*interviewsapi.InterviewRecordDTO, int64, error)

	GetInterviewRecord(ctx context.Context, userID uint, recordID uint64) (*interviewsapi.InterviewRecordDTO, error)
}
