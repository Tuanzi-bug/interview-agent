package handler

import (
	"context"
	"strconv"

	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// InterviewHandler 面试处理器
type InterviewHandler struct {
	interviewService *service.InterviewService
}

// NewInterviewHandler 创建面试处理器实例
func NewInterviewHandler() *InterviewHandler {
	return &InterviewHandler{
		interviewService: service.NewInterviewService(),
	}
}

// CreateInterview 创建面试接口
// @Summary 创建面试
// @Description 根据用户选择的岗位和简历创建新的面试
// @Tags 面试管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body service.CreateInterviewRequest true "创建面试请求参数"
// @Success 200 {object} map[string]interface{} "创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误或创建失败"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/interview [post]
func (h *InterviewHandler) CreateInterview(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	var req service.CreateInterviewRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	interview, err := h.interviewService.CreateInterview(userID, req)
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
		"data":    interview,
	})
}

// GetUserInterviews 获取用户面试列表接口
// @Summary 获取用户面试列表
// @Description 获取当前用户的所有面试记录
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/interview [get]
func (h *InterviewHandler) GetUserInterviews(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取状态参数（可选）
	status := ctx.Query("status")

	interviews, err := h.interviewService.GetUserInterviews(userID, status)
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
		"data":    interviews,
	})
}

// GetInterview 获取面试详情接口
// @Summary 获取面试详情
// @Description 根据面试ID获取面试详细信息
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 404 {object} map[string]interface{} "面试不存在"
// @Router /api/v1/interview/{id} [get]
func (h *InterviewHandler) GetInterview(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	interview, err := h.interviewService.GetInterviewByID(uint(interviewID), userID)
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
		"data":    interview,
	})
}

// StartInterview 开始面试接口
// @Summary 开始面试
// @Description 开始指定ID的面试，生成面试问题
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Success 200 {object} map[string]interface{} "开始成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID或开始失败"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 404 {object} map[string]interface{} "面试不存在"
// @Router /api/v1/interview/{id}/start [post]
func (h *InterviewHandler) StartInterview(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	interview, err := h.interviewService.StartInterview(uint(interviewID), userID)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "面试已开始",
		"data":    interview,
	})
}

// EndInterview 结束面试接口
// @Summary 结束面试
// @Description 结束指定ID的面试，生成面试结果
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Success 200 {object} map[string]interface{} "结束成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID或结束失败"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/interview/{id}/end [post]
func (h *InterviewHandler) EndInterview(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	interview, err := h.interviewService.EndInterview(uint(interviewID), userID)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "面试已结束",
		"data":    interview,
	})
}

// SubmitAnswer 提交回答接口
// @Summary 提交回答
// @Description 提交面试问题的回答
// @Tags 面试管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Param request body service.SubmitAnswerRequest true "提交回答请求参数"
// @Success 200 {object} map[string]interface{} "提交成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID或请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/interview/{id}/answer [post]
func (h *InterviewHandler) SubmitAnswer(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	var req service.SubmitAnswerRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	question, err := h.interviewService.SubmitAnswer(uint(interviewID), userID, req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "回答已提交",
		"data":    question,
	})
}

// GetInterviewQuestions 获取面试问题列表接口
// @Summary 获取面试问题列表
// @Description 获取指定面试ID的所有问题和回答
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Router /api/v1/interview/{id}/questions [get]
func (h *InterviewHandler) GetInterviewQuestions(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	interview, err := h.interviewService.GetInterviewByID(uint(interviewID), userID)
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
		"data":    interview.Questions,
	})
}

// GetInterviewResult 获取面试结果接口
// @Summary 获取面试结果
// @Description 获取指定ID面试的评价结果
// @Tags 面试管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "面试ID"
// @Success 200 {object} map[string]interface{} "获取成功"
// @Failure 400 {object} map[string]interface{} "无效的面试ID或面试未完成"
// @Failure 401 {object} map[string]interface{} "未授权访问"
// @Failure 404 {object} map[string]interface{} "面试不存在"
// @Router /api/v1/interview/{id}/result [get]
func (h *InterviewHandler) GetInterviewResult(c context.Context, ctx *app.RequestContext) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "未授权访问",
		})
		return
	}

	// 获取面试ID
	interviewIDStr := ctx.Param("id")
	interviewID, err := strconv.ParseUint(interviewIDStr, 10, 32)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "无效的面试ID",
		})
		return
	}

	interview, err := h.interviewService.GetInterviewByID(uint(interviewID), userID)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    404,
			"message": err.Error(),
		})
		return
	}

	// 检查面试是否已完成
	if interview.Status != "completed" {
		ctx.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "面试尚未完成，无法查看结果",
		})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "获取成功",
		"data": map[string]interface{}{
			"score":      interview.Score,
			"evaluation": interview.Evaluation,
			"duration":   interview.Duration,
			"questions":  interview.Questions,
		},
	})
}
