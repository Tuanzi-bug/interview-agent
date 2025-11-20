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

	// SaveInterviewDialogues 保存面试对话和问题主题
	SaveInterviewDialogues(ctx context.Context, userID uint, recordID uint64, questions []interface{}, dialogues []interface{}) error

	ListInterviewRecords(ctx context.Context, userID uint, page, pageSize *int32) ([]*interviewsapi.InterviewRecordDTO, int64, error)

	GetInterviewRecord(ctx context.Context, userID uint, recordID uint64) (*interviewsapi.InterviewRecordDTO, error)

	// GetInterviewEvaluation 根据用户ID和报告ID获取面试评估报告
	GetInterviewEvaluation(ctx context.Context, userID uint, reportID uint64) (interface{}, error)

	// GetAnswerReport 根据用户ID和报告ID获取答题报告
	GetAnswerReport(ctx context.Context, userID uint, reportID uint64) (interface{}, error)
}
