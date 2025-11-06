package handler

import (
	"context"
	"net/http"

	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Register 用户注册
func (h *UserHandler) Register(c context.Context, ctx *app.RequestContext) {
	var req service.RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.userService.Register(req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "注册成功",
		"data":    response,
	})
}

// Login 用户登录
func (h *UserHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req service.LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	response, err := h.userService.Login(req)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// GetProfile 获取用户资料
func (h *UserHandler) GetProfile(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取用户资料成功",
		"data":    user,
	})
}

// UpdateProfile 更新用户资料
func (h *UserHandler) UpdateProfile(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 解析请求参数
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(userID, req.Username, req.Email)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "更新失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新成功",
		"data":    user,
	})
}