package model

import (
	"gorm.io/gorm"
)

// getDB 获取数据库实例的函数变量，由 repository 包在初始化时设置
var getDB func() *gorm.DB

// SetDBGetter 设置数据库获取函数，由 repository 包在初始化时调用
func SetDBGetter(fn func() *gorm.DB) {
	getDB = fn
}

var UserModelDao _UserModel

type (
	_UserModel struct {
	}
	UserModel struct {
		ID              uint64 `json:"id" gorm:"primaryKey;autoIncrement"`
		UserID          int64  `json:"user_id" gorm:"index;not null"`
		Name            string `json:"name" gorm:"size:128;not null;comment:模型显示名称（用户维度唯一）"`
		ModelKey        string `json:"model_key" gorm:"size:128;not null;comment:模型标识（doubao-1.5-vision-lite-250315）"`
		Protocol        string `json:"protocol" gorm:"size:64;not null;comment:协议类型（openai/ark/claude/gemini/deepseek/ollama/qwen/ernie）"`
		BaseURL         string `json:"base_url" gorm:"size:255;not null;comment:API 基础地址"`
		APIKeyEncrypted string `json:"api_key_encrypted" gorm:"type:text;not null;comment:加密后的 API 密钥"`
		ConfigJSON      string `json:"config_json" gorm:"type:json;comment:额外配置（如区域、访问密钥等）"`
		SecretHint      string `json:"secret_hint" gorm:"size:32;default:'';comment:密钥脱敏提示（如显示末尾4位）"`
		ProviderName    string `json:"provider_name" gorm:"size:64;not null;comment:提供商名称（如 OpenAI、Ark、DeepSeek）"`
		MetaID          uint64 `json:"meta_id" gorm:"comment:关联全局 model_meta.id（继承能力/图标）"`
		DefaultParams   string `json:"default_params" gorm:"type:json;comment:默认参数（如 temperature、max_tokens）"`
		Scope           int    `json:"scope" gorm:"not null;default:7;comment:使用范围（位掩码：1=智能体, 2=应用, 4=工作流）"`
		Status          int    `json:"status" gorm:"not null;default:1;comment:状态（0=禁用, 1=启用）"`
		CreatedAt       int64  `json:"created_at" gorm:"not null;default:0;comment:创建时间（毫秒时间戳）"`
		UpdatedAt       int64  `json:"updated_at" gorm:"not null;default:0;comment:更新时间（毫秒时间戳）"`
		Deleted         int    `json:"deleted" gorm:"not null;default:0;comment:删除状态（0=未删除, 1=已删除）"`
	}
)

// TableName 指定表名为 user_model（而不是默认的 usermodels）
func (u *UserModel) TableName() string {
	return "user_model"
}

// CreateUserModel 创建用户模型
func (u *_UserModel) CreateUserModel(data *UserModel) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&UserModel{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

// ListUserModels 查询用户模型列表
func (u *_UserModel) ListUserModels(userID int64, page, pageSize int) ([]*UserModel, int64, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var models []*UserModel
	var total int64

	query := getDB().Model(&UserModel{}).
		Where("user_id = ? AND deleted = ?", userID, 0)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	return models, total, nil
}

// GetUserModelByID 根据ID查询用户模型
func (u *_UserModel) GetUserModelByID(userID int64, modelID int64) (*UserModel, error) {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	var userModel UserModel
	err := getDB().Model(&UserModel{}).
		Where("id = ? AND user_id = ? AND deleted = ?", modelID, userID, 0).
		First(&userModel).Error
	if err != nil {
		return nil, err
	}
	return &userModel, nil
}

// UpdateUserModel 更新用户模型
func (u *_UserModel) UpdateUserModel(data *UserModel) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&UserModel{}).
		Where("id = ? AND user_id = ? AND deleted = ?", data.ID, data.UserID, 0).
		Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserModel 软删除用户模型
func (u *_UserModel) DeleteUserModel(userID int64, modelID int64) error {
	if getDB == nil {
		panic("getDB function not initialized, please call model.SetDBGetter first")
	}
	err := getDB().Model(&UserModel{}).
		Where("id = ? AND user_id = ? AND deleted = ?", modelID, userID, 0).
		Update("deleted", 1).Error
	if err != nil {
		return err
	}
	return nil
}
