package questionbank

import (
	questionbankapi "ai-eino-interview-agent/api/model/questionbank"
	"ai-eino-interview-agent/internal/service/questionbank/impl"
	"context"
)

// NewQuestionBankService 初始化试题库服务的实例
func NewQuestionBankService() QuestionBankService {
	return impl.NewQuestionBankServiceImpl()
}

// QuestionBankService 试题库服务接口
type QuestionBankService interface {
	// GetMainCategories 获取所有大类（1类）标签
	GetMainCategories(ctx context.Context) (*questionbankapi.GetMainCategoriesResponse, error)

	// GetSubCategories 根据大类ID获取所有小分支（2类）
	GetSubCategories(ctx context.Context, parentID int64) (*questionbankapi.GetSubCategoriesResponse, error)
}
