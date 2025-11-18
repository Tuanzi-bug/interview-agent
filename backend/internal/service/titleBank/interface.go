package titlebank

import (
	titleBankApi "ai-eino-interview-agent/api/model/titleBank"
	"ai-eino-interview-agent/internal/service/titleBank/impl"
	"context"
)

func NewTitleBankManager() TitleBankManager {
	return impl.NewTitleBankServer()
}

type TitleBankManager interface {
	CreateInterviewTitle(ctx context.Context, req titleBankApi.CreateInterviewTitleRequest) (*titleBankApi.CreateInterviewTitleResponse, error)
	CreateInterviewLabel(ctx context.Context, req titleBankApi.CreateInterviewLabelRequest) (*titleBankApi.CreateInterviewLabelResponse, error)
	CreateInterviewParse(ctx context.Context, req titleBankApi.CreateInterviewParseRequest) (*titleBankApi.CreateInterviewParseResponse, error)
	CreateUserTitleInteract(ctx context.Context, req titleBankApi.CreateUserTitleInteractRequest) (*titleBankApi.CreateUserTitleInteractResponse, error)
}
