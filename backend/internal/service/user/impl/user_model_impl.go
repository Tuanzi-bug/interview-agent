package impl

import (
	"ai-eino-interview-agent/api/model/user"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/service/common"
	"context"
	"errors"
	"time"
)

type UserModelServer struct {
}

func NewUserModelServer() *UserModelServer {
	return &UserModelServer{}
}

// CreateUserModel 创建用户模型
func (s *UserModelServer) CreateUserModel(ctx context.Context,
	userID int64,
	req user.CreateUserModelRequest) (string, error) {
	apiKey, err := common.EncryptAPIKey(req.APIKey)
	if err != nil {
		return "fail", errors.New("密钥加密失败")
	}

	// 处理可选字段
	var metaID uint64
	if req.IsSetMetaID() {
		metaID = uint64(*req.MetaID)
	}

	var configJSON string
	if req.IsSetConfigJSON() {
		configJSON = *req.ConfigJSON
	}

	var defaultParams string
	if req.IsSetDefaultParams() {
		defaultParams = *req.DefaultParams
	}

	scope := 7 // 默认值
	if req.IsSetScope() {
		scope = int(*req.Scope)
	}

	status := 1 // 默认值
	if req.IsSetStatus() {
		status = int(*req.Status)
	}

	err = model.UserModelDao.CreateUserModel(&model.UserModel{
		UserID:          userID,
		Name:            req.GetName(),
		ModelKey:        req.GetModelKey(),
		Protocol:        req.GetProtocol(),
		BaseURL:         req.GetBaseURL(),
		APIKeyEncrypted: apiKey,
		ConfigJSON:      configJSON,
		MetaID:          metaID,
		DefaultParams:   defaultParams,
		Scope:           scope,
		Status:          status,
		ProviderName:    req.GetProviderName(),
	})
	if err != nil {
		return "fail", err
	}
	return "success", nil
}

// ListUserModels 列出用户模型
func (s *UserModelServer) ListUserModels(ctx context.Context,
	userID int64,
	req user.ListUserModelsRequest) ([]*model.UserModel, int64, error) {
	page := int(req.GetPage())
	if page <= 0 {
		page = 1
	}
	pageSize := int(req.GetSize())
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return model.UserModelDao.ListUserModels(userID, page, pageSize)
}

// UserModelDetail 获取用户模型详情
func (s *UserModelServer) UserModelDetail(ctx context.Context,
	userID int64,
	modelID int64) (*model.UserModel, error) {
	return model.UserModelDao.GetUserModelByID(userID, modelID)
}

// UpdateUserModel 更新用户模型
func (s *UserModelServer) UpdateUserModel(ctx context.Context,
	userID int64,
	req user.UpdateUserModelRequest) error {
	// 先获取现有模型
	existingModel, err := model.UserModelDao.GetUserModelByID(userID, req.GetID())
	if err != nil {
		return err
	}

	// 更新字段
	existingModel.Name = req.GetName()
	existingModel.ModelKey = req.GetModelKey()
	existingModel.Protocol = req.GetProtocol()
	existingModel.BaseURL = req.GetBaseURL()
	existingModel.ProviderName = req.GetProviderName()

	// 处理可选字段 - APIKey
	if req.IsSetAPIKey() {
		apiKey, err := common.EncryptAPIKey(req.GetAPIKey())
		if err != nil {
			return errors.New("密钥加密失败")
		}
		existingModel.APIKeyEncrypted = apiKey
	}

	// 处理可选字段 - MetaID
	if req.IsSetMetaID() {
		existingModel.MetaID = uint64(req.GetMetaID())
	}

	// 处理可选字段 - DefaultParams
	if req.IsSetDefaultParams() {
		existingModel.DefaultParams = req.GetDefaultParams()
	}

	// 处理可选字段 - ConfigJSON
	if req.IsSetConfigJSON() {
		existingModel.ConfigJSON = req.GetConfigJSON()
	}

	// 处理可选字段 - Scope
	if req.IsSetScope() {
		existingModel.Scope = int(req.GetScope())
	}

	// 处理可选字段 - Status
	if req.IsSetStatus() {
		existingModel.Status = int(req.GetStatus())
	}

	existingModel.UpdatedAt = time.Now().UnixMilli()
	return model.UserModelDao.UpdateUserModel(existingModel)
}

// DeleteUserModel 删除用户模型
func (s *UserModelServer) DeleteUserModel(ctx context.Context,
	userID int64,
	modelID int64) error {
	return model.UserModelDao.DeleteUserModel(userID, modelID)
}
