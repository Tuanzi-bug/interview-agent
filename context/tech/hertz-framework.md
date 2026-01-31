# Hertz 框架集成指南

**最后更新**: 2026-02-01  
**维护者**: 后端团队

## 概述

Hertz 是字节跳动开源的高性能 HTTP 框架，为 Go 开发者提供高效、易用的 Web 服务开发能力。在面试吧平台中，Hertz 作为主 Web 框架处理所有 API 请求。

### 核心特性

- 🚀 **高性能**: 基准测试优于主流 Go 框架
- 📦 **即开即用**: 标准库风格 API，学习曲线低
- 🔧 **易于扩展**: 中间件、路由、错误处理都很灵活
- 🎯 **零拷贝**: 网络 I/O 优化，减少内存分配
- 🌐 **生产级验证**: 字节跳动内部大规模使用

## 项目中的 Hertz 应用

### 应用结构

```
backend/
├── main.go                 # 应用入口
├── api/
│   ├── handler/           # HTTP 处理器
│   │   ├── interview/     # 面试相关处理器
│   │   ├── resume/        # 简历相关处理器
│   │   └── user/          # 用户相关处理器
│   ├── model/             # 请求/响应模型
│   ├── response/          # 响应封装
│   └── router/            # 路由配置
│       ├── register.go    # 路由注册
│       └── middleware/    # 中间件
└── config.yaml            # 配置文件
```

### 初始化流程

```go
// main.go
func main() {
    // 1. 加载配置
    cfg := config.LoadConfig("config.yaml")
    
    // 2. 初始化服务
    svc := service.NewServices(cfg)
    
    // 3. 创建 Hertz 引擎
    h := hertz.Default()
    
    // 4. 注册中间件
    h.Use(middleware.Logger())
    h.Use(middleware.Recovery())
    h.Use(middleware.AuthMiddleware())
    
    // 5. 注册路由
    router.RegisterRoutes(h, svc)
    
    // 6. 启动服务
    h.Spin()
}
```

## 路由管理

### 基础路由

```go
// api/router/register.go
func RegisterRoutes(h *hertz.Engine, svc *service.Services) {
    // 用户路由组
    userGroup := h.Group("/api/v1/user")
    {
        userGroup.POST("/register", handler.Register)
        userGroup.POST("/login", handler.Login)
        userGroup.GET("/profile", handler.GetProfile)
        userGroup.PUT("/profile", handler.UpdateProfile)
    }
    
    // 面试路由组
    interviewGroup := h.Group("/api/v1/interview")
    {
        interviewGroup.POST("/create", handler.CreateInterview)
        interviewGroup.GET("/:id", handler.GetInterview)
        interviewGroup.GET("/:id/question", handler.GetQuestion)
        interviewGroup.POST("/:id/answer", handler.SubmitAnswer)
    }
    
    // 简历路由组
    resumeGroup := h.Group("/api/v1/resume")
    {
        resumeGroup.POST("/upload", handler.UploadResume)
        resumeGroup.GET("/list", handler.ListResumes)
        resumeGroup.POST("/:id/analyze", handler.AnalyzeResume)
    }
}
```

### 路由分组和中间件

```go
// 创建分组并应用中间件
api := h.Group("/api/v1")
api.Use(middleware.AuthRequired())

{
    // 仅需认证的路由
    api.GET("/user/profile", handler.GetProfile)
    api.POST("/interview/create", handler.CreateInterview)
    
    // 管理员路由组
    admin := api.Group("/admin")
    admin.Use(middleware.AdminRequired())
    {
        admin.GET("/users", handler.ListUsers)
        admin.DELETE("/user/:id", handler.DeleteUser)
    }
}

// 公开路由（无认证）
h.POST("/auth/register", handler.Register)
h.POST("/auth/login", handler.Login)
```

### 参数传递

```go
// 路径参数
func GetInterview(ctx context.Context, c *app.RequestContext) {
    id := c.Param("id")
    interview := svc.GetInterview(id)
    c.JSON(200, interview)
}

// 查询参数
func ListInterviews(ctx context.Context, c *app.RequestContext) {
    page := c.Query("page", "1")
    size := c.Query("size", "10")
    
    interviews := svc.ListInterviews(page, size)
    c.JSON(200, interviews)
}

// 请求体
func CreateInterview(ctx context.Context, c *app.RequestContext) {
    var req CreateInterviewRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, "Invalid request")
        return
    }
    
    interview := svc.CreateInterview(req)
    c.JSON(201, interview)
}
```

## 中间件开发

### 认证中间件

```go
// api/router/middleware/auth.go
func AuthRequired() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 获取 Token
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, "Missing authorization header")
            c.Abort()
            return
        }
        
        // 验证 Token
        claims, err := jwt.ValidateToken(token)
        if err != nil {
            c.JSON(401, "Invalid token")
            c.Abort()
            return
        }
        
        // 存储用户信息到上下文
        c.Set("user_id", claims.UserID)
        c.Set("user_role", claims.Role)
        
        // 继续处理请求
        c.Next(ctx)
    }
}

// 使用认证中间件
h.Use(middleware.AuthRequired())
```

### 日志中间件

```go
// api/router/middleware/logger.go
func Logger() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        startTime := time.Now()
        
        // 记录请求信息
        logger.Info("incoming request",
            "method", c.Request.Method,
            "path", c.Request.URI,
            "ip", c.ClientIP(),
        )
        
        // 继续处理
        c.Next(ctx)
        
        // 记录响应信息
        duration := time.Since(startTime)
        logger.Info("request completed",
            "status", c.Response.StatusCode,
            "duration", duration,
        )
    }
}
```

### 错误恢复中间件

```go
// api/router/middleware/recovery.go
func Recovery() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("panic recovered", "error", err)
                c.JSON(500, "Internal server error")
            }
        }()
        
        c.Next(ctx)
    }
}
```

## 请求和响应处理

### 统一响应格式

```go
// api/response/response.go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

// 成功响应
func Success(c *app.RequestContext, data interface{}) {
    c.JSON(200, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

// 错误响应
func Error(c *app.RequestContext, code int, message string) {
    c.JSON(code, Response{
        Code:    -1,
        Message: message,
    })
}

// 使用示例
func GetProfile(ctx context.Context, c *app.RequestContext) {
    userID := c.GetString("user_id")
    profile := svc.GetProfile(userID)
    response.Success(c, profile)
}
```

### 请求数据绑定

```go
// api/model/interview.go
type CreateInterviewRequest struct {
    InterviewType string `json:"interview_type" binding:"required"` // 面试类型
    Duration      int    `json:"duration" binding:"required"`       // 面试时长(分钟)
    Difficulty    string `json:"difficulty" binding:"required"`     // 难度: easy/medium/hard
}

// 绑定和验证
func CreateInterview(ctx context.Context, c *app.RequestContext) {
    var req CreateInterviewRequest
    if err := c.BindJSON(&req); err != nil {
        response.Error(c, 400, "Invalid request parameters")
        return
    }
    
    // 调用服务
    interview := svc.CreateInterview(req)
    response.Success(c, interview)
}
```

## 上传文件处理

```go
// 处理简历上传
func UploadResume(ctx context.Context, c *app.RequestContext) {
    // 获取上传的文件
    file, err := c.FormFile("resume")
    if err != nil {
        response.Error(c, 400, "Failed to get file")
        return
    }
    
    // 保存文件
    savePath := fmt.Sprintf("uploads/%s_%d", userID, time.Now().Unix())
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        response.Error(c, 500, "Failed to save file")
        return
    }
    
    // 异步处理文件
    go func() {
        svc.ProcessResumeFile(savePath)
    }()
    
    response.Success(c, map[string]string{"path": savePath})
}
```

## WebSocket 支持

```go
// WebSocket 端点 - 实时面试
func InterviewWebSocket(ctx context.Context, c *app.RequestContext) {
    err := c.ProtocolUpgrade(
        func(ctx *websocket.Conn) {
            // WebSocket 连接建立
            userID := c.GetString("user_id")
            
            for {
                // 接收消息
                var msg InterviewMessage
                err := ctx.ReadJSON(&msg)
                if err != nil {
                    break
                }
                
                // 处理消息
                response := svc.ProcessInterviewMessage(msg)
                
                // 发送响应
                if err := ctx.WriteJSON(response); err != nil {
                    break
                }
            }
        },
    )
    
    if err != nil {
        response.Error(c, 500, "WebSocket upgrade failed")
    }
}

// 注册路由
h.GET("/ws/interview/:id", handler.InterviewWebSocket)
```

## 错误处理

### 定义错误类型

```go
// internal/errors/errors.go
type ErrorCode int

const (
    ErrInvalidRequest ErrorCode = 400
    ErrUnauthorized   ErrorCode = 401
    ErrNotFound       ErrorCode = 404
    ErrInternalServer ErrorCode = 500
)

type APIError struct {
    Code    ErrorCode
    Message string
}

// 自定义错误
func NewError(code ErrorCode, message string) *APIError {
    return &APIError{Code: code, Message: message}
}
```

### 错误处理中间件

```go
func ErrorHandler() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("handler error", "error", err)
                c.JSON(500, Response{
                    Code:    -1,
                    Message: "Internal server error",
                })
            }
        }()
        
        // 检查上下文中的错误
        c.Next(ctx)
        
        // 处理返回的错误
        // ...
    }
}
```

## 性能优化

### 连接池配置

```yaml
# config.yaml
server:
  port: 8080
  read_timeout: 10s
  write_timeout: 10s
  max_request_body_size: 10485760  # 10MB
  
http_client:
  max_idle_conns: 100
  max_idle_conns_per_host: 2
  idle_conn_timeout: 90s
```

### 响应压缩

```go
// 启用 Gzip 压缩
import "github.com/hertz-contrib/gzip"

h := hertz.Default()
h.Use(gzip.Gzip(
    gzip.DefaultCompression,
    gzip.WithExcludedPathsRegexs([]string{".png$", ".jpg$"}),
))
```

### 缓存控制

```go
// 设置缓存头
func CacheControl(duration time.Duration) app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(duration.Seconds())))
        c.Next(ctx)
    }
}

// 使用
h.GET("/api/public-data", handler.GetPublicData)
h.Use(middleware.CacheControl(1 * time.Hour))
```

## 测试

### 单元测试示例

```go
// api/handler/user_test.go
func TestGetProfile(t *testing.T) {
    // 创建测试请求
    req := httptest.NewRequest("GET", "/api/v1/user/profile", nil)
    
    // 创建测试上下文
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 执行请求
    w := httptest.NewRecorder()
    handler.GetProfile(ctx, w)
    
    // 验证结果
    assert.Equal(t, 200, w.Code)
}
```

## 常见问题

**Q: 如何在 Hertz 中使用数据库?**
→ 使用 GORM 作为 ORM，在 handler 中注入 service 层

**Q: 如何处理 CORS?**
→ 使用 Hertz 的 CORS 中间件或在路由中设置响应头

**Q: 如何实现请求日志?**
→ 创建日志中间件，在每个请求前后记录信息

**Q: 如何优化响应时间?**
→ 使用缓存、数据库查询优化、异步处理等

## 相关资源

### 官方文档
- [Hertz 官方文档](https://www.hertzframework.com/)
- [Hertz GitHub](https://github.com/cloudwego/hertz)

### 中间件库
- [Hertz 社区中间件](https://github.com/hertz-contrib)

### 项目代码
- [Handler 实现](../../backend/api/handler/)
- [Router 配置](../../backend/api/router/)
- [响应处理](../../backend/api/response/)

### 相关文档
- [架构设计](architecture.md)
- [Eino 框架](eino-framework.md)

---

**维护者**: 后端团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
