package user

import (
	"ai-eino-interview-agent/api/model/user"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/service/user/impl"
	"context"
)

// NewModelManager 返回 ModelManager 接口的默认实现
func NewModelManager() ModelManager {
	return impl.NewUserModelServer()
}

type ModelManager interface {
	CreateUserModel(
		ctx context.Context,
		userID int64,
		req user.CreateUserModelRequest,
	) (string, error)
	ListUserModels(
		ctx context.Context,
		userID int64,
		req user.ListUserModelsRequest,
	) ([]*model.UserModel, int64, error)
	UserModelDetail(
		ctx context.Context,
		userID int64,
		modelID int64,
	) (*model.UserModel, error)
	UpdateUserModel(
		ctx context.Context,
		userID int64,
		req user.UpdateUserModelRequest,
	) error
	DeleteUserModel(
		ctx context.Context,
		userID int64,
		modelID int64,
	) error
}
