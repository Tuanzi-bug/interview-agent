package service

import (
	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

type WechatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid,omitempty"`
}

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid,omitempty"`
}

// WechatLoginResponse 微信登录响应
type WechatLoginResponse struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

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

// GetWechatAccessToken 通过code获取微信access_token和openid
func (s *UserService) GetWechatAccessToken(code string) (*WechatTokenResponse, error) {
	// 构造请求URL
	reqURL := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?"+
		"appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		config.Global.Wechat.AppID,
		config.Global.Wechat.AppSecret,
		code)

	// 发送HTTP请求
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信响应失败: %v", err)
	}

	// 解析JSON响应
	var tokenResp WechatTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析微信响应失败: %v", err)
	}

	// 检查是否有错误
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("获取access_token失败: %s", string(body))
	}

	return &tokenResp, nil
}

// GetWechatUserInfo 通过access_token和openid获取微信用户信息
func (s *UserService) GetWechatUserInfo(accessToken, openID string) (*WechatUserInfo, error) {
	// 构造请求URL
	reqURL := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?"+
		"access_token=%s&openid=%s",
		accessToken,
		openID)

	// 发送HTTP请求
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信用户信息接口失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取微信用户信息响应失败: %v", err)
	}

	// 解析JSON响应
	var userInfo WechatUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析微信用户信息响应失败: %v", err)
	}

	// 检查是否有错误
	if userInfo.OpenID == "" {
		return nil, fmt.Errorf("获取用户信息失败: %s", string(body))
	}

	return &userInfo, nil
}

// WechatLoginOrRegister 微信登录或注册
func (s *UserService) WechatLoginOrRegister(wechatUser *WechatUserInfo) (*WechatLoginResponse, error) {
	db := repository.GetDB()

	// 查找用户
	var user model.User
	result := db.Where("wechat_open_id = ?", wechatUser.OpenID).First(&user)

	// 如果用户不存在，则创建新用户
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = model.User{
			Username:      fmt.Sprintf("wechat_%s", wechatUser.OpenID[:10]), // 生成默认用户名
			Email:         "",                                               // 微信登录用户可能没有邮箱
			PasswordHash:  "",                                               // 微信登录不需要密码
			Role:          "user",
			WechatOpenID:  wechatUser.OpenID,
			WechatUnionID: wechatUser.UnionID,
			Nickname:      wechatUser.Nickname,
			Avatar:        wechatUser.HeadImgURL,
		}

		if err := db.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("创建用户失败: %v", err)
		}
	} else if result.Error != nil {
		return nil, fmt.Errorf("查询用户失败: %v", result.Error)
	} else {
		// 如果用户存在，更新用户信息
		updates := map[string]interface{}{
			"nickname": wechatUser.Nickname,
			"avatar":   wechatUser.HeadImgURL,
		}

		// 如果之前没有union_id，现在更新
		if user.WechatUnionID == "" && wechatUser.UnionID != "" {
			updates["wechat_union_id"] = wechatUser.UnionID
		}

		if err := db.Model(&user).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新用户信息失败: %v", err)
		}
	}

	// 生成JWT token
	token, err := middleware.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	return &WechatLoginResponse{
		Token: token,
		User:  user,
	}, nil
}
