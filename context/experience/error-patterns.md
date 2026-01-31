# 错误模式和解决方案

**最后更新**: 2026-02-01  
**维护者**: 技术团队

## 概述

记录项目开发中遇到的常见错误、根本原因和解决方案。

## Go 相关错误

### 错误 1: goroutine 泄漏

**症状**: 应用内存占用不断增加，无法释放

**原因**: 
- goroutine 阻塞在 channel 接收/发送
- 未关闭 channel
- 循环中创建的 goroutine 未能正常退出

**解决方案**:
```go
// ❌ 错误：goroutine 永不退出
go func() {
    for {
        data := <-ch  // 如果没有发送方，永久阻塞
        process(data)
    }
}()

// ✅ 正确：使用 context 控制
go func(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case data := <-ch:
            process(data)
        }
    }
}(ctx)

// 在适当时机取消
cancel()
```

**预防**:
- 使用 `context` 进行超时/取消控制
- 定期检查 goroutine 数量: `runtime.NumGoroutine()`
- 使用工具检测: `go test -race`

---

## Hertz 框架错误

### 错误 2: Hertz 路由返回 404

**症状**: 请求正确的路由仍返回 404

**原因**:
- 路由未正确注册
- 路由版本不匹配 (v1 vs v2)
- 中间件终止了请求链 (未调用 `c.Next()`)

**解决方案**:
```go
// ❌ 错误：路由未注册到正确的分组
h.GET("/api/profile", handler)  // 注册到根路径

// ✅ 正确：使用正确的 API 路径
api := h.Group("/api/v1")
api.GET("/user/profile", handler)

// 验证：启用调试模式查看已注册路由
h.DebugPrintRouteFunc = nil
h.Spin()
```

**预防**:
- 检查路由注册顺序
- 验证中间件是否调用 `c.Next(ctx)`
- 使用路由调试功能

---

### 错误 3: Hertz 中间件未生效

**症状**: 中间件逻辑没有执行

**原因**:
- 中间件未调用 `c.Next(ctx)`
- 中间件注册顺序错误
- 中间件返回了错误/重定向

**解决方案**:
```go
// ❌ 错误：中间件没有继续处理
func MyMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 处理逻辑
        // 忘记调用 c.Next(ctx)！
    }
}

// ✅ 正确：中间件必须调用 Next
func MyMiddleware() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        // 处理逻辑
        c.Next(ctx)  // 继续处理链
    }
}

// 注册顺序
h.Use(middleware.GlobalMiddleware1())  // 先执行
h.Use(middleware.GlobalMiddleware2())  // 后执行
h.GET("/path", middleware.LocalMiddleware(), handler)
```

**预防**:
- 始终在中间件中调用 `c.Next(ctx)`
- 理解中间件的洋葱模型执行顺序

---

## 数据库错误

### 错误 4: 数据库连接超时

**症状**: `context deadline exceeded` 或 `connection refused`

**原因**:
- 数据库连接池已满
- 数据库响应慢
- 网络连接问题

**解决方案**:
```go
// ❌ 错误：没有设置连接超时
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// ✅ 正确：配置连接和读写超时
dsn := "user:password@tcp(localhost:3306)/db?timeout=10s&readTimeout=10s&writeTimeout=10s"
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 配置连接池
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
sqlDB.SetConnMaxLifetime(5 * time.Minute)
```

**预防**:
- 监控数据库连接数
- 优化慢查询
- 使用连接池配置

---

### 错误 5: GORM 查询返回空结果

**症状**: 数据确实存在，但 GORM 返回 ErrRecordNotFound

**原因**:
- 未检查错误:  `if err != nil` vs `if err == gorm.ErrRecordNotFound`
- 查询条件错误
- 字段映射问题

**解决方案**:
```go
// ❌ 错误：假设 First 一定成功
var user User
db.First(&user, id)
fmt.Println(user.Email)  // 如果记录不存在，可能是零值

// ✅ 正确：检查错误
var user User
result := db.First(&user, id)
if result.Error != nil {
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        // 处理记录不存在
    }
    return nil, result.Error
}
```

**预防**:
- 总是检查 GORM 的 `result.Error`
- 理解 `First`、`Find`、`Take` 的区别
- 添加日志输出 SQL 语句: `db.WithContext(ctx).Debug()`

---

## Redis 错误

### 错误 6: Redis 连接池泄漏

**症状**: Redis 连接数不断增加，应用响应变慢

**原因**:
- 未关闭 Redis 连接
- 连接池配置不合理

**解决方案**:
```go
// ❌ 错误：未关闭连接
redis.Set(ctx, "key", "value", 0)
redis.Get(ctx, "key")

// ✅ 正确：正确使用连接池
redis := redis.NewClient(&redis.Options{
    Addr:        "localhost:6379",
    MaxRetries:  3,
    PoolSize:    10,
})
defer redis.Close()

// 操作
redis.Set(ctx, "key", "value", 0)
result := redis.Get(ctx, "key")
```

**预防**:
- 在应用启动时创建单例连接池
- 在优雅关闭时关闭连接池

---

## AI 集成错误

### 错误 7: LLM 请求超时

**症状**: `context deadline exceeded` 或 `timeout`

**原因**:
- LLM API 响应慢
- 网络问题
- Token 数过多导致生成缓慢

**解决方案**:
```go
// ❌ 错误：没有设置超时
resp := llm.Generate(context.Background(), prompt)

// ✅ 正确：设置适当的超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp := llm.Generate(ctx, prompt)

// 或者使用重试机制
resp, err := RetryWithBackoff(func() (*Response, error) {
    return llm.Generate(ctx, prompt)
}, 3)
```

**预防**:
- 根据场景调整超时时间
- 实现重试机制和降级策略
- 监控 LLM 响应时间

---

### 错误 8: Eino Agent 执行失败

**症状**: Agent 返回错误或无效结果

**原因**:
- 提示词（Prompt）不清晰
- 工具调用失败
- 模型上下文长度超限

**解决方案**:
```go
// ❌ 错误：提示词太模糊
prompt := "分析简历"

// ✅ 正确：清晰详细的提示词
prompt := `
分析以下简历，提供:
1. 技能评估 (掌握程度: 初级/中级/高级)
2. 缺失技能
3. 职业发展建议

简历内容:
` + resumeContent

// 添加错误处理
result, err := agent.Execute(ctx, prompt)
if err != nil {
    logger.Error("agent execution failed", "error", err)
    // 降级处理或重试
}
```

**预防**:
- 使用清晰的提示词模板
- 添加详细的上下文信息
- 测试不同的输入场景

---

## 并发/同步错误

### 错误 9: 竞态条件 (Race Condition)

**症状**: 不同运行中结果不同，调试时难以复现

**原因**: 多个 goroutine 同时访问共享资源

**解决方案**:
```go
// ❌ 错误：并发访问共享变量
var counter int
for i := 0; i < 1000; i++ {
    go func() {
        counter++  // 竞态条件！
    }()
}

// ✅ 正确：使用互斥锁
var (
    counter int
    mu      sync.Mutex
)
for i := 0; i < 1000; i++ {
    go func() {
        mu.Lock()
        counter++
        mu.Unlock()
    }()
}

// 或使用原子操作
var counter atomic.Int32
for i := 0; i < 1000; i++ {
    go func() {
        counter.Add(1)
    }()
}

// 检测竞态条件
// go test -race ./...
```

**预防**:
- 运行 `go test -race` 检测竞态条件
- 使用 Channel 进行通信而不是共享内存
- 尽量使用不可变数据结构

---

## 常用解决方案速查表

| 错误 | 快速修复 | 预防措施 |
|------|---------|--------|
| goroutine 泄漏 | 使用 context 取消 | 定期检查 goroutine 数 |
| 404 错误 | 检查路由注册 | 使用调试模式 |
| DB 超时 | 增加超时时间 | 配置连接池 |
| 竞态条件 | 使用 Mutex/Channel | `go test -race` |
| LLM 超时 | 增加超时或重试 | 优化提示词长度 |

---

**维护者**: 技术团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
