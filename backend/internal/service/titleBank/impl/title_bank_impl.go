package impl

import (
	titleBankApi "ai-eino-interview-agent/api/model/titleBank"
	"ai-eino-interview-agent/internal/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type TitleBankServer struct {
}

func NewTitleBankServer() *TitleBankServer {
	return &TitleBankServer{}
}

func (t *TitleBankServer) CreateInterviewTitle(ctx context.Context, req titleBankApi.CreateInterviewTitleRequest) (*titleBankApi.CreateInterviewTitleResponse, error) {
	_, err := model.InterviewTitleDao.FindInterviewTitleByTitle(req.GetTitle())
	if err == nil {
		return nil, errors.New("面试题已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	interviewTitle := &model.InterviewTitle{
		Title:  req.GetTitle(),
		Type:   req.GetType(),
		Domain: req.GetDomain(),
		Level:  req.GetLevel(),
	}

	if err := model.InterviewTitleDao.CreateInterviewTitle(interviewTitle); err != nil {
		return nil, err
	}

	return &titleBankApi.CreateInterviewTitleResponse{
		"success",
	}, nil
}
func (t *TitleBankServer) CreateInterviewLabel(ctx context.Context, req titleBankApi.CreateInterviewLabelRequest) (*titleBankApi.CreateInterviewLabelResponse, error) {
	interviewTitle, err := model.InterviewTitleDao.FindInterviewTitleById(uint64(req.GetTitleID()))
	if interviewTitle == nil {
		return nil, errors.New("面试题不存在")
	}
	if err != nil {
		return nil, err
	}
	interviewLabel := &model.InterviewLabel{
		Label:   req.GetLabel(),
		TitleId: uint64(req.GetTitleID()),
	}

	if err := model.InterviewLabelDao.CreateInterviewLabel(interviewLabel); err != nil {
		return nil, err
	}

	return &titleBankApi.CreateInterviewLabelResponse{
		"success",
	}, nil
}
func (t *TitleBankServer) CreateInterviewParse(ctx context.Context, req titleBankApi.CreateInterviewParseRequest) (*titleBankApi.CreateInterviewParseResponse, error) {
	interviewTitle, err := model.InterviewTitleDao.FindInterviewTitleById(uint64(req.GetTitleID()))
	if interviewTitle == nil {
		return nil, errors.New("面试题不存在")
	}
	if err != nil {
		return nil, err
	}
	interviewParse := &model.InterviewParse{
		TitleId: uint64(req.GetTitleID()),
		Parse:   req.GetParse(),
		Url:     req.GetURL(),
	}

	if err := model.InterviewParseDao.CreateInterviewParse(interviewParse); err != nil {
		return nil, err
	}

	return &titleBankApi.CreateInterviewParseResponse{
		"success",
	}, nil
}
func (t *TitleBankServer) CreateUserTitleInteract(ctx context.Context, req titleBankApi.CreateUserTitleInteractRequest) (*titleBankApi.CreateUserTitleInteractResponse, error) {
	interviewTitle, err := model.InterviewTitleDao.FindInterviewTitleById(uint64(req.GetTitleID()))
	if interviewTitle == nil {
		return nil, errors.New("面试题不存在")
	}
	if err != nil {
		return nil, err
	}

	user, err := model.UserDao.FindByID(uint(req.GetUserID()))
	if user == nil {
		return nil, errors.New("用户不存在")
	}
	if err != nil {
		return nil, err
	}
	userTitleInteract := &model.UserTitleInteract{
		UserId:   uint64(req.GetUserID()),
		TitleId:  uint64(req.GetTitleID()),
		Interact: int(req.GetInteract()),
	}

	if err := model.UserTitleInteractDao.CreateUserTitleInteract(userTitleInteract); err != nil {
		return nil, err
	}

	return &titleBankApi.CreateUserTitleInteractResponse{
		"success",
	}, nil
}
