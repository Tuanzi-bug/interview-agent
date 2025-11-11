package hertz

import (
	"context"
	"fmt"
	"log"
	"time"

	"ai-eino-interview-agent/internal/config"
	"ai-eino-interview-agent/internal/handler"
	"ai-eino-interview-agent/internal/middleware"
	"ai-eino-interview-agent/internal/model"
	"ai-eino-interview-agent/internal/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// HertzServer 全局Hertz服务器实例
var HertzServer *server.Hertz

// InitHertz 初始化Hertz框架
func InitHertz() error {
	// 直接使用Global配置的Host和Port
	HertzServer = server.Default(server.WithHostPorts(fmt.Sprintf("%s:%d", config.Global.Host, config.Global.Port)))

	// 配置中间件
	configureMiddleware()

	// 配置路由
	configureRoutes()

	// 配置Swagger文档
	SetupSwagger()

	log.Println("Hertz框架初始化成功")
	return nil
}

// GetHertzServer 获取Hertz服务器实例
func GetHertzServer() *server.Hertz {
	return HertzServer
}

// configureMiddleware 配置中间件
func configureMiddleware() {
	// 启用CORS
	if config.Global.Security.CORS.AllowOrigins != nil && len(config.Global.Security.CORS.AllowOrigins) > 0 {
		HertzServer.Use(CORSMiddleware())
	}

	// 直接启用日志中间件
	HertzServer.Use(AccessLogMiddleware())

	// 添加JWT认证中间件（某些路由需要）
	// HertzServer.Use(JWTMiddleware())
}

// configureRoutes 配置路由
func configureRoutes() {
	// 创建处理器实例
	userHandler := handler.NewUserHandler()
	resumeHandler := handler.NewResumeHandler()
	interviewHandler := handler.NewInterviewHandler()

	// 健康检查
	HertzServer.GET("/health", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(consts.StatusOK, utils.H{
			"status":  "ok",
			"message": "Interview Agent Service is running",
		})
	})

	// API路由组
	api := HertzServer.Group("/api/v1")
	{
		// 用户相关路由
		user := api.Group("/user")
		{
			user.POST("/register", userHandler.Register)
			user.POST("/login", userHandler.Login)
			user.GET("/profile", middleware.JWTMiddleware(), userHandler.GetProfile)
			user.PUT("/profile", middleware.JWTMiddleware(), userHandler.UpdateProfile)
			wechat := user.Group("/wechat")
			{
				// 微信登录相关路由
				wechat.GET("/login", userHandler.WechatLogin)
				wechat.GET("/callback", userHandler.WechatCallback)
			}
		}

		// 简历相关路由（需要认证）
		resume := api.Group("/resume", middleware.JWTMiddleware())
		{
			resume.POST("", resumeHandler.CreateResume)
			resume.GET("", resumeHandler.GetUserResumes)
			resume.GET("/:id", resumeHandler.GetResume)
			resume.PUT("/:id", resumeHandler.UpdateResume)
			resume.DELETE("/:id", resumeHandler.DeleteResume)
			resume.POST("/:id/analyze", resumeHandler.AnalyzeResume)
		}

		// 面试相关路由（需要认证）
		interview := api.Group("/interview", middleware.JWTMiddleware())
		{
			interview.POST("", interviewHandler.CreateInterview)
			interview.GET("", interviewHandler.GetUserInterviews)
			interview.GET("/:id", interviewHandler.GetInterview)
			interview.POST("/:id/start", interviewHandler.StartInterview)
			interview.POST("/:id/end", interviewHandler.EndInterview)
			interview.POST("/:id/answer", interviewHandler.SubmitAnswer)
			interview.GET("/:id/questions", interviewHandler.GetInterviewQuestions)
			interview.GET("/:id/result", interviewHandler.GetInterviewResult)
		}

		// 评估标准相关路由（需要认证）
		criteria := api.Group("/criteria", middleware.JWTMiddleware())
		{
			criteria.GET("", func(c context.Context, ctx *app.RequestContext) {
				// 从数据库获取评估标准
				var criteria []model.EvaluationCriteria
				db := repository.GetDB()
				err := db.Find(&criteria).Error
				if err != nil {
					ctx.JSON(consts.StatusInternalServerError, utils.H{
						"error": "Failed to fetch evaluation criteria",
					})
					return
				}
				ctx.JSON(consts.StatusOK, criteria)
			})
		}
	}
}

// 临时处理器函数（占位符）
func registerHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Register endpoint"})
}

func loginHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Login endpoint"})
}

func getUserProfileHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get user profile endpoint"})
}

func updateUserProfileHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Update user profile endpoint"})
}

func createResumeHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Create resume endpoint"})
}

func getUserResumesHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get user resumes endpoint"})
}

func getResumeHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get resume endpoint"})
}

func updateResumeHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Update resume endpoint"})
}

func deleteResumeHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Delete resume endpoint"})
}

func analyzeResumeHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Analyze resume endpoint"})
}

func createInterviewHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Create interview endpoint"})
}

func getUserInterviewsHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get user interviews endpoint"})
}

func getInterviewHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get interview endpoint"})
}

func startInterviewHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Start interview endpoint"})
}

func endInterviewHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "End interview endpoint"})
}

func submitAnswerHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Submit answer endpoint"})
}

func getInterviewQuestionsHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get interview questions endpoint"})
}

func getInterviewResultHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get interview result endpoint"})
}

func getEvaluationCriteriaHandler(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"message": "Get evaluation criteria endpoint"})
}

// CORSMiddleware CORS中间件
func CORSMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if string(ctx.Request.Method()) == "OPTIONS" {
			ctx.AbortWithStatus(consts.StatusNoContent)
			return
		}

		ctx.Next(c)
	}
}

// AccessLogMiddleware 访问日志中间件
func AccessLogMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 记录请求开始时间
		timestamp := time.Now()

		// 处理请求
		ctx.Next(c)

		// 记录日志
		latency := time.Since(timestamp)
		log.Printf("[%s] %s %s %d %v",
			timestamp.Format("2006-01-02 15:04:05"),
			string(ctx.Request.Method()),
			string(ctx.Request.URI().Path()),
			ctx.Response.StatusCode(),
			latency,
		)
	}
}
