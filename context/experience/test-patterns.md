# 测试模式

**最后更新**: 2026-02-01  
**维护者**: QA 团队

## 概述

记录在开发过程中发现的有效测试模式和反模式。

## 测试组织

### 表驱动测试

**优势**: 清晰、易于添加新测试用例、代码量少

```go
func TestUserLogin(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        password string
        wantErr bool
        wantUser bool
    }{
        {
            name:     "valid credentials",
            email:    "test@example.com",
            password: "correct",
            wantErr:  false,
            wantUser: true,
        },
        {
            name:     "wrong password",
            email:    "test@example.com",
            password: "wrong",
            wantErr:  true,
            wantUser: false,
        },
        {
            name:     "user not found",
            email:    "notexist@example.com",
            password: "any",
            wantErr:  true,
            wantUser: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := Login(tt.email, tt.password)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if (user != nil) != tt.wantUser {
                t.Errorf("Login() user = %v, wantUser %v", user, tt.wantUser)
            }
        })
    }
}
```

---

## Mock 策略

### 接口 Mock

**最佳实践**: 针对接口编程，便于 Mock

```go
// 定义接口
type UserRepository interface {
    GetUser(id string) (*User, error)
    SaveUser(user *User) error
}

// 实现接口
type DatabaseRepository struct {
    db *gorm.DB
}

// Mock 实现
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetUser(id string) (*User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

// 在测试中使用
func TestUserService(t *testing.T) {
    mockRepo := new(MockRepository)
    mockRepo.On("GetUser", "123").Return(&User{ID: "123", Email: "test@example.com"}, nil)
    
    svc := NewUserService(mockRepo)
    user, err := svc.GetUserInfo("123")
    
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
    mockRepo.AssertExpectations(t)
}
```

### HTTP 客户端 Mock

```go
func TestExternalAPICall(t *testing.T) {
    // 创建 Mock HTTP 客户端
    mockClient := &MockHTTPClient{}
    mockClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
        return req.URL.Path == "/api/external"
    })).Return(&http.Response{
        StatusCode: 200,
        Body:       ioutil.NopCloser(strings.NewReader(`{"result":"ok"}`)),
    }, nil)
    
    // 使用 Mock 客户端
    client := NewClient(mockClient)
    result, err := client.CallExternalAPI()
    
    assert.NoError(t, err)
    assert.Equal(t, "ok", result)
}
```

---

## 断言技巧

### 使用 testify/assert

```go
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    // 基本断言
    assert.NoError(t, err)
    assert.Error(t, err)
    
    // 值断言
    assert.Equal(t, expected, actual)
    assert.NotEqual(t, unexpected, actual)
    
    // 比较
    assert.Greater(t, actual, expected)
    assert.Less(t, actual, expected)
    
    // 集合
    assert.Contains(t, []int{1, 2, 3}, 2)
    assert.NotContains(t, []int{1, 2, 3}, 4)
    assert.Len(t, slice, 3)
    
    // 条件
    assert.True(t, condition)
    assert.False(t, condition)
    
    // 指针
    assert.Nil(t, nilValue)
    assert.NotNil(t, notNilValue)
}
```

---

## 并发测试

### 数据竞态检测

```go
// 运行测试时启用竞态检测
// go test -race ./...

func TestConcurrentAccess(t *testing.T) {
    var counter int
    var mu sync.Mutex
    
    // 并发写入
    for i := 0; i < 100; i++ {
        go func() {
            mu.Lock()
            counter++
            mu.Unlock()
        }()
    }
    
    // 给 goroutine 时间完成
    time.Sleep(100 * time.Millisecond)
    
    mu.Lock()
    assert.Equal(t, 100, counter)
    mu.Unlock()
}
```

### WaitGroup 同步

```go
func TestParallelOperations(t *testing.T) {
    var wg sync.WaitGroup
    results := make(chan string, 10)
    
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // 执行异步操作
            result := DoSomething(id)
            results <- result
        }(i)
    }
    
    // 等待所有 goroutine 完成
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // 收集结果
    var actualResults []string
    for result := range results {
        actualResults = append(actualResults, result)
    }
    
    assert.Len(t, actualResults, 10)
}
```

---

## 测试工具

### 使用 testing.T 的便利功能

```go
func TestWithHelpers(t *testing.T) {
    // Helper 函数 - 不会计入测试覆盖率的行数
    assertUserEqual := func(t *testing.T, expected, actual *User) {
        t.Helper()
        if expected.ID != actual.ID {
            t.Errorf("User ID mismatch: expected %s, got %s", expected.ID, actual.ID)
        }
    }
    
    user, _ := GetUser("123")
    expected := &User{ID: "123", Email: "test@example.com"}
    assertUserEqual(t, expected, user)
}
```

### 测试命名

```go
// ✅ 好的测试名称：清晰地表达测试意图
func TestUserLogin_WithValidCredentials_ReturnsUser(t *testing.T)
func TestUserLogin_WithWrongPassword_ReturnsError(t *testing.T)
func TestInterviewService_CreateInterview_WithValidRequest_CreatesRecord(t *testing.T)

// ❌ 不好的名称
func TestLogin(t *testing.T)
func Test1(t *testing.T)
```

---

## 集成测试模式

### 使用 Testcontainers

```go
func TestWithDatabase(t *testing.T) {
    // 启动 MySQL 容器
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "mysql:8.0",
        ExposedPorts: []string{"3306/tcp"},
        Env: map[string]string{
            "MYSQL_ROOT_PASSWORD": "root",
            "MYSQL_DATABASE":      "test",
        },
    }
    
    container, err := testcontainers.GenericContainer(ctx, 
        testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    // 获取连接信息
    host, _ := container.Host(ctx)
    port, _ := container.MappedPort(ctx, "3306")
    
    // 连接数据库进行测试
    dsn := fmt.Sprintf("root:root@tcp(%s:%s)/test", host, port.Port())
    db, _ := gorm.Open(mysql.Open(dsn))
    
    // 运行迁移和测试
    db.AutoMigrate(&User{})
    user := &User{Email: "test@example.com"}
    db.Create(user)
    
    assert.NotZero(t, user.ID)
}
```

---

## 反模式

### 反模式 1: 过度 Mock

**问题**: 几乎所有依赖都被 Mock，测试变成了教科书实现

```go
// ❌ 过度 Mock - 测试没有意义
func TestBusinessLogic(t *testing.T) {
    mockDB := new(MockDB)
    mockCache := new(MockCache)
    mockLogger := new(MockLogger)
    mockHTTP := new(MockHTTP)
    
    // 现在测试基本上验证了 mock 的行为，而不是真实逻辑
}

// ✅ 适度 Mock - 只 Mock 外部依赖
func TestBusinessLogic(t *testing.T) {
    // 真实的业务逻辑
    cache := setupTestCache()  // 真实的缓存
    
    // Mock 只用于 HTTP 外部 API
    mockHTTP := new(MockHTTPClient)
    mockHTTP.On("Do", mock.Anything).Return(response)
}
```

### 反模式 2: 测试之间的依赖

```go
// ❌ 测试之间有依赖 - 运行顺序重要
var globalState *User

func TestCreateUser(t *testing.T) {
    globalState = &User{ID: "1"}
}

func TestGetUser(t *testing.T) {
    user := GetUser(globalState.ID)  // 依赖前一个测试！
}

// ✅ 测试独立 - 每个测试都完整
func TestGetUserAfterCreate(t *testing.T) {
    user := CreateUser(&CreateUserRequest{Email: "test@example.com"})
    retrieved := GetUser(user.ID)
    assert.Equal(t, user.ID, retrieved.ID)
}
```

---

**维护者**: QA 团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
