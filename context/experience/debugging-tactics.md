# 调试技巧

**最后更新**: 2026-02-01  
**维护者**: 技术团队

## 概述

实用的调试策略和工具使用指南。

## 日志分析

### 结构化日志

```go
import "go.uber.org/zap"

// 创建结构化日志
logger, _ := zap.NewProduction()
defer logger.Sync()

// 详细日志记录
logger.Info("user created",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
    zap.Duration("duration", time.Since(start)),
)

// 错误日志
logger.Error("failed to create user",
    zap.Error(err),
    zap.String("email", email),
    zap.Stack("stacktrace"),
)
```

### 调试日志级别

```go
// 开发环境：启用 DEBUG
// 生产环境：只有 INFO 及以上

if isDebug {
    logger.Debug("request details",
        zap.Any("headers", req.Headers),
        zap.String("body", body),
    )
}
```

---

## 性能分析

### CPU 分析

```bash
# 运行带 CPU 分析的测试
go test -cpuprofile=cpu.prof -bench=. ./...

# 分析结果
go tool pprof cpu.prof
(pprof) top10              # 查看前 10 个函数
(pprof) list FunctionName  # 查看具体函数
```

### 内存分析

```bash
# 运行带内存分析的程序
go run -memprofile=mem.prof main.go

# 分析
go tool pprof mem.prof
(pprof) top -cum  # 按累计内存排序
(pprof) alloc_space  # 查看分配
```

### goroutine 分析

```go
import "runtime"

// 检查 goroutine 数量
func CheckGoroutines() {
    numGoroutines := runtime.NumGoroutine()
    log.Printf("Active goroutines: %d", numGoroutines)
}

// 定期检查（发现泄漏）
ticker := time.NewTicker(1 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
}
```

---

## 竞态条件检测

```bash
# 启用竞态检测运行测试
go test -race ./...

# 启用竞态检测运行程序
go run -race main.go

# 竞态检测输出示例
# WARNING: DATA RACE
# Write at 0x00c0001b2180 by goroutine 18:
#      main.incrementCounter()
#           main.go:25 +0x44
```

---

## 调试工具

### Delve 调试器

```bash
# 编译调试信息
go build -gcflags="all=-N -l" -o app main.go

# 启动 Delve
dlv exec ./app

# 常用命令
(dlv) break main.main      # 设置断点
(dlv) continue             # 继续执行
(dlv) next                 # 单步执行
(dlv) step                 # 进入函数
(dlv) locals               # 查看局部变量
(dlv) print variableName   # 打印变量
(dlv) backtrace            # 打印堆栈
```

### 远程调试

```go
// 在 main 中启动 Delve 调试服务器
import "github.com/go-delve/delve/service"

go func() {
    listener, err := net.Listen("tcp", ":2345")
    if err != nil {
        log.Fatal(err)
    }
    _ = service.NewServer(listener)
}()

// 客户端连接
// dlv connect localhost:2345
```

---

## 数据库调试

### SQL 日志

```go
// 启用 GORM SQL 日志
db := db.WithContext(ctx).Debug()

// 执行查询
var users []User
db.Find(&users)
// 输出：
// [2026-02-01 10:30:00] [39.123ms] SELECT * FROM `users`
```

### 慢查询日志

```go
// MySQL 配置
mysql> SET GLOBAL slow_query_log = 'ON';
mysql> SET GLOBAL slow_query_log_file = '/var/log/mysql/slow.log';
mysql> SET GLOBAL long_query_time = 0.5;  # 0.5 秒

# 分析慢查询
mysqldumpslow -s c -t 10 /var/log/mysql/slow.log  # 计数前 10
mysqldumpslow -s t -t 10 /var/log/mysql/slow.log  # 时间前 10
```

---

## API 调试

### 使用 curl 测试

```bash
# 基本请求
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"pass"}'

# 带认证
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer TOKEN"

# 显示响应头
curl -i http://localhost:8080/api/user/profile

# 保存响应到文件
curl -o response.json http://localhost:8080/api/data
```

### 使用 Postman/Insomnia

- 快速测试 API
- 保存请求集合
- 自动化测试流程

---

## 日志样本

### 良好的调试日志示例

```go
logger.Info("interview created",
    zap.String("interview_id", interview.ID),
    zap.String("user_id", userID),
    zap.String("type", interview.Type),
    zap.Duration("processing_time", duration),
)

logger.Error("failed to evaluate answer",
    zap.String("interview_id", interviewID),
    zap.String("question_id", questionID),
    zap.Error(err),
    zap.String("error_type", fmt.Sprintf("%T", err)),
    zap.Stack("stacktrace"),
)
```

---

**维护者**: 技术团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
