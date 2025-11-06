package handler

import (
	"context"
	"strconv"

	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// ResumeHandler 简历处理器
type ResumeHandler struct {
	resumeService *service.ResumeService
}

// NewResumeHandler 创建简历处理器实例
func NewResumeHandler() *ResumeHandler {
	return &ResumeHandler{
		resumeService: service.NewResumeService(),
	}
}

// CreateResume 创建简历接口
// @Summary 创建简历
// @Description 为当前登录用户创建一个新的简历
// @Tags 简历管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateResumeRequest true "创建简历请求参数"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 500 {object} map[string]interface{} "创建失败"
// @Router /api/v1/resume [post]
func (h *ResumeHandler) CreateResume(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	var req service.CreateResumeRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resume, err := h.resumeService.CreateResume(userID, req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "创建失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "创建成功",
		"data":    resume,
	})
}

// GetUserResumes 获取用户的所有简历接口
// @Summary 获取用户简历列表
// @Description 获取当前登录用户的所有简历
// @Tags 简历管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 500 {object} map[string]interface{} "获取失败"
// @Router /api/v1/resume [get]
func (h *ResumeHandler) GetUserResumes(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	resumes, err := h.resumeService.GetUserResumes(userID)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "获取失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取成功",
		"data":    resumes,
	})
}

// GetResume 获取单个简历接口
// @Summary 获取简历详情
// @Description 根据简历ID获取指定简历的详细信息
// @Tags 简历管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "简历ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的简历ID"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 404 {object} map[string]interface{} "简历不存在"
// @Router /api/v1/resume/{id} [get]
func (h *ResumeHandler) GetResume(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取简历ID
	resumeIDStr := ctx.Param("id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	resume, err := h.resumeService.GetResumeByID(uint(resumeID), userID)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取成功",
		"data":    resume,
	})
}

// UpdateResume 更新简历接口
// @Summary 更新简历
// @Description 更新指定ID的简历信息
// @Tags 简历管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "简历ID"
// @Param request body service.CreateResumeRequest true "更新简历请求参数"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 404 {object} map[string]interface{} "简历不存在"
// @Failure 500 {object} map[string]interface{} "更新失败"
// @Router /api/v1/resume/{id} [put]
func (h *ResumeHandler) UpdateResume(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取简历ID
	resumeIDStr := ctx.Param("id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	var req service.CreateResumeRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	resume, err := h.resumeService.UpdateResume(uint(resumeID), userID, req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "更新成功",
		"data":    resume,
	})
}

// DeleteResume 删除简历接口
// @Summary 删除简历
// @Description 删除指定ID的简历
// @Tags 简历管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "简历ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "无效的简历ID或删除失败"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/resume/{id} [delete]
func (h *ResumeHandler) DeleteResume(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取简历ID
	resumeIDStr := ctx.Param("id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	err = h.resumeService.DeleteResume(uint(resumeID), userID)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "删除成功",
	})
}

// AnalyzeResume 分析简历接口
// @Summary 分析简历
// @Description 使用AI分析指定ID的简历内容
// @Tags 简历管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "简历ID"
// @Success 200 {object} map[string]interface{} "分析成功"
// @Failure 400 {object} map[string]interface{} "无效的简历ID"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 500 {object} map[string]interface{} "分析失败"
// @Router /api/v1/resume/{id}/analyze [post]
func (h *ResumeHandler) AnalyzeResume(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取简历ID
	resumeIDStr := ctx.Param("id")
	resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的简历ID",
		})
		return
	}

	resume, err := h.resumeService.AnalyzeResume(uint(resumeID), userID)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "分析失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "分析成功",
		"data":    resume,
	})
}