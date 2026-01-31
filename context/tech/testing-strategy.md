# 测试策略

**最后更新**: 2026-02-01  
**维护者**: QA 团队

## 概述

本文档定义了面试吧平台的测试策略，包括单元测试、集成测试和 E2E 测试的最佳实践。

## 测试层次

```
┌──────────────────────┐
│   E2E 测试           │ 端到端测试，完整业务流
│   (低频率，高价值)    │
├──────────────────────┤
│   集成测试           │ 服务间集成，DB 交互
│   (中等频率)          │
├──────────────────────┤
│   单元测试           │ 单一函数/方法
│   (高频率，快速)      │
└──────────────────────┘
```

## 单元测试

### 测试框架
- **主框架**: Go 标准库 `testing`
- **断言库**: `github.com/stretchr/testify/assert`
- **Mock 工具**: `github.com/stretchr/testify/mock`

### 编写规范

```go
// user_service_test.go
func TestUserLogin(t *testing.T) {
    // Arrange - 准备测试数据
    testEmail := "test@example.com"
    testPassword := "password123"
    
    // Act - 执行被测代码
    user, err := service.Login(testEmail, testPassword)
    
    // Assert - 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, testEmail, user.Email)
}
```

### Mock 对象

```go
// Mock 外部依赖
type MockDB struct {
    mock.Mock
}

func (m *MockDB) GetUser(id string) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// 使用 Mock
func TestGetUserInfo(t *testing.T) {
    mockDB := new(MockDB)
    mockDB.On("GetUser", "123").Return(&User{ID: "123"}, nil)
    
    // 使用 mockDB
}
```

## 集成测试

### 测试数据库

```go
// 使用 Docker 容器运行测试 MySQL
func TestWithDB(t *testing.T) {
    // 使用 testcontainers 启动 MySQL
    container, db, err := setupTestDB()
    if err != nil {
        t.Fatalf("failed to setup db: %v", err)
    }
    defer container.Terminate(context.Background())
    
    // 运行迁移
    if err := runMigrations(db); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }
    
    // 测试
    user := &User{Email: "test@example.com"}
    db.Create(user)
    
    var retrieved User
    db.First(&retrieved, user.ID)
    assert.Equal(t, user.Email, retrieved.Email)
}
```

### 集成测试示例

```go
// interview_service_test.go
func TestCreateInterview(t *testing.T) {
    // 设置测试环境
    db := setupTestDB()
    cache := setupTestCache()
    svc := service.NewInterviewService(db, cache)
    
    // 创建用户
    user := createTestUser(db)
    
    // 创建面试
    interview := svc.CreateInterview(&CreateInterviewRequest{
        UserID:   user.ID,
        Type:     "comprehensive",
        Duration: 60,
    })
    
    // 验证结果
    assert.NotNil(t, interview)
    assert.Equal(t, user.ID, interview.UserID)
    
    // 验证数据库
    var stored Interview
    db.First(&stored, interview.ID)
    assert.Equal(t, interview.ID, stored.ID)
}
```

## E2E 测试

### API 端到端测试

```go
// 启动完整应用，测试 API 流程
func TestInterviewFlow(t *testing.T) {
    // 启动应用
    app := startTestApp()
    defer app.Stop()
    
    client := &http.Client{}
    
    // 1. 用户登录
    loginResp := post(t, client, "/api/auth/login", LoginRequest{
        Email:    "test@example.com",
        Password: "password",
    })
    assert.Equal(t, 200, loginResp.StatusCode)
    token := extractToken(loginResp)
    
    // 2. 创建面试
    interviewResp := post(t, client, "/api/interview/create", 
        CreateInterviewRequest{Type: "comprehensive"},
        header("Authorization", "Bearer "+token))
    assert.Equal(t, 201, interviewResp.StatusCode)
    interviewID := extractID(interviewResp)
    
    // 3. 获取问题
    questionResp := get(t, client, "/api/interview/"+interviewID+"/question",
        header("Authorization", "Bearer "+token))
    assert.Equal(t, 200, questionResp.StatusCode)
    
    // 4. 提交答案
    answerResp := post(t, client, "/api/interview/"+interviewID+"/answer",
        AnswerRequest{Question: "...", Answer: "..."},
        header("Authorization", "Bearer "+token))
    assert.Equal(t, 200, answerResp.StatusCode)
}
```

## 测试覆盖率

### 目标覆盖率
- 关键业务逻辑: **> 90%**
- 服务层: **> 80%**
- 工具函数: **> 70%**
- 整体平均: **> 75%**

### 计算覆盖率

```bash
# 生成覆盖率报告
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 显示覆盖率百分比
go tool cover -func=coverage.out
```

## 持续集成

### CI 流程

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
    
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.24
      
      - name: Run tests
        run: go test ./... -v -cover
      
      - name: Run integration tests
        run: go test ./... -v -cover -tags=integration
```

## 性能测试

### Benchmark 测试

```go
// performance_test.go
func BenchmarkLogin(b *testing.B) {
    svc := setupService()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        svc.Login("test@example.com", "password")
    }
}

// 运行 benchmark
// go test -bench=. -benchmem
```

---

**维护者**: QA 团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
