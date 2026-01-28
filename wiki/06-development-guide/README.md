# 实战开发指南

> 掌握系统后，如何添加新功能和进行开发

## 📚 本指南内容

1. [如何添加新功能](add-new-feature.md)
2. [编码规范](coding-standards.md)
3. [测试指南](testing-guide.md)
4. [部署指南](deployment-guide.md)

---

## 开发流程总览

```
需求分析
  ↓
技术设计
  ↓
数据建模
  ↓
开发实现
  ├── 1. 定义Model
  ├── 2. 实现Repository
  ├── 3. 实现Service
  ├── 4. 实现Handler
  ├── 5. 注册路由
  └── 6. (可选) 实现Agent
  ↓
单元测试
  ↓
集成测试
  ↓
代码审查
  ↓
部署上线
  ↓
监控运维
```

---

## 快速开始

### 添加一个简单的API接口

假设我们要添加一个"获取用户统计数据"的接口：

#### 步骤1: 定义数据模型 (如果需要)

```go
// backend/internal/model/user_stats.go
package model

type UserStats struct {
    UserID           uint    `json:"user_id"`
    TotalInterviews  int     `json:"total_interviews"`
    AverageScore     float64 `json:"average_score"`
    // ...
}
```

#### 步骤2: 实现Repository

```go
// backend/internal/repository/user/stats.go
func (r *UserRepository) GetUserStats(userID uint) (*model.UserStats, error) {
    // 实现数据查询逻辑
}
```

#### 步骤3: 实现Service

```go
// backend/internal/service/user/stats_service.go
func (s *UserService) GetUserStatistics(userID uint) (*model.UserStats, error) {
    // 业务逻辑处理
    return s.repo.GetUserStats(userID)
}
```

#### 步骤4: 实现Handler

```go
// backend/api/handler/user/stats_handler.go
func GetUserStats(c context.Context, ctx *app.RequestContext) {
    userID := getUserIDFromContext(ctx)
    
    stats, err := userService.GetUserStatistics(userID)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    
    response.Success(ctx, stats)
}
```

#### 步骤5: 注册路由

```go
// backend/api/router/register.go
userGroup.GET("/stats", userHandler.GetUserStats)
```

#### 步骤6: 测试

```bash
curl http://localhost:8080/api/v1/user/stats \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 添加AI功能

### 添加一个新的智能体

假设我们要添加一个"模拟面试智能体"：

#### 步骤1: 创建智能体目录

```bash
mkdir -p backend/chatApp/agent/mock_interview
```

#### 步骤2: 定义智能体结构

```go
// backend/chatApp/agent/mock_interview/agent.go
package mock_interview

import (
    "github.com/cloudwego/eino/components/agent"
)

type MockInterviewAgent struct {
    // 配置和依赖
}

func NewMockInterviewAgent() *MockInterviewAgent {
    return &MockInterviewAgent{}
}

func (a *MockInterviewAgent) Run(ctx context.Context, input string) (string, error) {
    // 实现智能体逻辑
}
```

#### 步骤3: 设计Prompt

```go
const systemPrompt = `
你是一位资深的模拟面试教练...
`
```

#### 步骤4: 注册工具(如需要)

```go
tools := []agent.Tool{
    // 注册可用工具
}
```

#### 步骤5: 集成到Service

```go
// backend/chatApp/agent_service/mock_interview/service.go
```

---

## 开发检查清单

### 功能开发前

- [ ] 需求是否明确
- [ ] 技术方案是否确定
- [ ] 数据库设计是否完成
- [ ] API接口设计是否确定

### 开发过程中

- [ ] 遵循项目的编码规范
- [ ] 添加必要的注释
- [ ] 处理错误情况
- [ ] 考虑性能和安全

### 开发完成后

- [ ] 编写单元测试
- [ ] 手动测试功能
- [ ] 更新API文档
- [ ] 代码审查

---

## 常见开发任务

### 1. 添加新的数据表

1. 在`internal/model/`定义模型
2. 运行自动迁移或手动编写SQL
3. 在`internal/repository/`实现数据访问
4. 更新`db_schema.sql`

### 2. 修改现有API

1. 找到对应的Handler
2. 修改业务逻辑(Service层)
3. 如需修改数据结构，更新Model
4. 更新测试用例
5. 更新API文档

### 3. 优化数据库查询

1. 分析慢查询日志
2. 添加索引
3. 优化查询语句
4. 考虑添加缓存

### 4. 修改AI Prompt

1. 找到对应的智能体代码
2. 修改Prompt模板
3. 测试不同场景
4. 记录修改原因

---

## 调试技巧

### 1. 打印日志

```go
log.Printf("Debug: user_id=%d, data=%+v", userID, data)
```

### 2. 使用调试器

```bash
# 使用delve
dlv debug backend/main.go
```

### 3. 查看数据库查询

```go
// 启用GORM日志
db.Debug().Where("id = ?", id).First(&user)
```

### 4. 查看Redis数据

```bash
redis-cli
> KEYS pattern*
> GET key_name
```

---

## 性能优化

### 数据库优化

- [ ] 添加必要的索引
- [ ] 使用连接池
- [ ] 避免N+1查询
- [ ] 使用批量操作

### 缓存优化

- [ ] 缓存热点数据
- [ ] 设置合理的过期时间
- [ ] 使用缓存预热
- [ ] 处理缓存击穿

### API优化

- [ ] 使用分页
- [ ] 压缩响应数据
- [ ] 使用CDN
- [ ] 实现接口限流

---

## 安全注意事项

### 输入验证

```go
// 使用binding标签验证
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

### SQL注入防护

```go
// ✅ 使用参数化查询
db.Where("email = ?", email).First(&user)

// ❌ 避免字符串拼接
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

### XSS防护

```go
// 对输出进行转义
import "html"
safeOutput := html.EscapeString(userInput)
```

### CSRF防护

- 使用CSRF Token
- 验证Referer
- 使用SameSite Cookie

---

## ❓ 常见问题

### Q1: 如何处理并发请求？
**解决方案**: _____

### Q2: 如何实现事务？
```go
db.Transaction(func(tx *gorm.DB) error {
    // 事务操作
    return nil
})
```

### Q3: 如何处理大文件上传？
**解决方案**: _____

---

## 📚 参考资源

### 官方文档
- [Hertz文档](https://www.cloudwego.io/zh/docs/hertz/)
- [Eino文档](https://www.cloudwego.io/zh/docs/eino/)
- [GORM文档](https://gorm.io/zh_CN/docs/)

### 代码示例
- [Hertz Examples](https://github.com/cloudwego/hertz-examples)
- [Eino Examples](https://github.com/cloudwego/eino/tree/main/examples)

---

## ✅ 开发就绪检查

- [ ] 理解完整的开发流程
- [ ] 掌握各层的职责
- [ ] 了解编码规范
- [ ] 知道如何调试
- [ ] 了解性能优化方法
- [ ] 重视安全问题

---

**详细指南**: 
- [添加新功能详细步骤](add-new-feature.md)
- [编码规范](coding-standards.md)
- [测试指南](testing-guide.md)
