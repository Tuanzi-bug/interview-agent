package service

import (
	"errors"

	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*LoginResponse, error) {
	db := repository.GetDB()

	// 检查用户名是否已存在
	var existingUser model.User
	result := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser)
	if result.Error == nil {
		return nil, errors.New("用户名或邮箱已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 加密密码（实际项目中应该使用bcrypt等安全的加密方式）
	passwordHash := req.Password // 简化实现，实际需要加密

	// 创建用户
	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         "user",
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	db := repository.GetDB()

	// 查找用户
	var user model.User
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}

	// 验证密码（简化实现）
	if user.PasswordHash != req.Password {
		return nil, errors.New("密码错误")
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	db := repository.GetDB()

	var user model.User
	result := db.First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, username string, email string) (*model.User, error) {
	db := repository.GetDB()

	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	// 更新字段
	updates := map[string]interface{}{}
	if username != "" {
		updates["username"] = username
	}
	if email != "" {
		updates["email"] = email
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新加载用户信息
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
