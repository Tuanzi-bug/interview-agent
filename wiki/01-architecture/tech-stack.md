# 技术栈详解

> 深入了解系统使用的技术栈及其选型原因

## 📋 学习任务

1. 理解每个技术组件的作用
2. 分析技术选型的原因
3. 学习各技术的核心概念
4. 了解技术之间的协作方式

---

## 技术栈总览

```
┌─────────────────────────────────────────┐
│          前端层 (Next.js)                │
├─────────────────────────────────────────┤
│          API网关层 (Hertz)               │
├─────────────────────────────────────────┤
│   业务逻辑层 (Go + Eino AI Framework)    │
├─────────────────────────────────────────┤
│  数据层 (MySQL + Redis + Milvus)        │
└─────────────────────────────────────────┘
```

---

## 后端技术栈

### 1. Hertz Web框架

#### 基本信息
- **官方网站**: https://github.com/cloudwego/hertz
- **开源方**: 字节跳动
- **语言**: Go
- **类型**: 高性能HTTP框架

#### 核心特性
1. **高性能**: 基于字节跳动自研的高性能网络库Netpoll，QPS比Gin高40%
2. **完整中间件**: 内置Recovery、日志、跨域、限流等中间件
3. **自动代码生成**: 通过IDL（接口定义语言）自动生成路由和Handler代码
4. **HTTP2支持**: 原生支持HTTP2协议
5. **易于扩展**: 插件化设计，支持自定义中间件和扩展

#### 选型原因
✅ **优势**:
- **高性能**: 相比传统框架（Gin、Echo），性能提升明显，适合高并发场景
- **易用性**: API设计友好，与Gin类似，学习成本低
- **扩展性**: 完善的中间件机制，支持链式调用，方便功能扩展
- **企业级**: 字节跳动内部大规模使用，稳定性有保障
- **生态完善**: 与Eino等字节系工具无缝集成

❌ **劣势**:
- 社区相对较小，第三方插件不如Gin丰富
- 文档主要是中文，国际化支持一般
- 学习曲线略陡（特别是IDL代码生成部分）

#### 在项目中的应用
- **位置**: `backend/api/`、`backend/main.go`
- **作用**: 提供HTTP API服务，处理前端请求
- **核心使用**:
  ```go
  // 创建Hertz服务器
  s := server.Default(server.WithHostPorts("0.0.0.0:8888"))
  
  // 注册中间件
  s.Use(routerMiddleware.Recovery())
  s.Use(appMiddleware.JWTMiddleware())
  
  // 注册路由
  router.GeneratedRegister(s)
  
  // 启动服务
  s.Run()
  ```

**实际应用场景**:
- API路由管理（面试、用户、简历等模块）
- 请求参数验证和绑定
- 响应数据序列化（JSON）
- 中间件处理（认证、CORS、错误处理）
- 文件上传下载（简历PDF）

---

### 2. Eino AI框架

#### 基本信息
- **官方网站**: https://github.com/cloudwego/eino
- **开源方**: 字节跳动
- **语言**: Go
- **类型**: 大语言模型应用开发框架

#### 核心概念

##### Agent (智能体)
**定义**: Agent是具备自主决策能力的AI实体，能够根据目标和环境信息，调用工具、执行任务并做出决策。

**特点**:
- **自主性**: 能够独立分析问题并选择合适的工具
- **记忆能力**: 维护对话上下文，实现连续对话
- **工具调用**: 可以调用外部工具获取信息或执行操作
- **链式思考**: 支持Chain-of-Thought推理

**在项目中的应用**: 
- **面试Agent**: `chatApp/agent/interview/` - 进行专项面试（Go、Java、MySQL、Redis等）
- **综合面试Agent**: `chatApp/agent/interview/comprehensive/` - 校招/社招综合面试
- **简历分析Agent**: `chatApp/agent/resume/` - 解析和分析简历
- **押题Agent**: `chatApp/agent/prediction/` - 根据简历预测面试题
- **评估Agent**: `chatApp/agent/record_evaluation/` - 评估面试表现

##### Tool (工具)
**定义**: Tool是Agent可以调用的外部能力，用于获取信息、执行操作或与外部系统交互。

**类型**:
- **本地工具**: 直接调用系统内部函数（如数据库查询）
  - `GetResumeInfoTool`: 获取简历信息
  - `GetMianshiInfoTool`: 获取面试题库信息
- **远程工具**: 调用外部API服务
  - `GoogleSearchTool`: Google搜索工具
  - Milvus检索工具（向量相似度搜索）

**示例**:
```go
// 定义工具
func GetResumeInfo(ctx context.Context, req *GetResumeInfoRequest) (*GetResumeInfoResponse, error) {
    data, err := model.ResumeDao.GetResumeByID(req.ResumeID)
    return &GetResumeInfoResponse{Data: data}, err
}

// 创建工具
tool, _ := utils.InferTool(
    "get_resume_info",
    "获取用户的解析后的简历信息",
    GetResumeInfo,
)
```

##### Chain (链)
**定义**: Chain是将多个组件（LLM调用、工具、处理器）按顺序串联的执行流程。

**作用**: 
- 实现复杂的多步骤任务
- 每个步骤的输出作为下一步的输入
- 适合固定流程的任务（如：简历解析 → 信息提取 → 结构化输出）

##### Flow (流)
**定义**: Flow是更复杂的执行流程，支持条件分支、并行执行和循环。

**应用场景**: 
- 根据用户回答决定下一个问题（面试流程）
- 并行调用多个工具获取信息
- 复杂的决策树逻辑

#### 选型原因
✅ **为什么选择Eino**:
1. **性能优势**: Go语言原生实现，性能远超Python的LangChain
2. **类型安全**: 静态类型检查，减少运行时错误
3. **易于集成**: 与Go后端无缝集成，不需要跨语言调用
4. **企业级**: 字节跳动内部使用，稳定性和可维护性有保障
5. **并发友好**: Go的goroutine天然支持高并发场景
6. **资源占用低**: 相比Python，内存和CPU占用更低

**与LangChain对比**:
| 特性 | Eino | LangChain |
|------|------|-----------|
| 语言 | Go | Python |
| 性能 | ⚡ 高性能，低延迟 | 🐌 相对较慢 |
| 并发 | ✅ 原生goroutine | ❌ GIL限制 |
| 类型安全 | ✅ 编译时检查 | ❌ 运行时检查 |
| 生态 | 🔨 正在发展 | 🌟 非常丰富 |
| 学习曲线 | 📚 适中 | 📖 较平缓 |
| 企业支持 | 字节跳动 | LangChain Inc. |
| 部署 | 📦 单二进制文件 | 🐍 需要Python环境 |

#### 在项目中的应用
- **位置**: `backend/chatApp/`
- **智能体实现**: `backend/chatApp/agent/`
  - 面试Agent（专项、综合）
  - 简历分析Agent
  - 押题Agent
  - 评估Agent
- **工具实现**: `backend/chatApp/tool/`
  - 简历信息工具
  - 面试题库工具
  - Google搜索工具
  - Milvus检索工具

**典型使用流程**:
```
用户提交简历 → 简历解析Agent → 提取关键信息
                ↓
        押题Agent（调用简历工具）
                ↓
        生成预测面试题
                ↓
        面试Agent（调用题库工具）
                ↓
        进行面试对话
                ↓
        评估Agent
                ↓
        生成评估报告
```

---

### 3. GORM (数据库ORM)

#### 基本信息
- **官方网站**: https://gorm.io
- **类型**: Go ORM库
- **支持数据库**: MySQL, PostgreSQL, SQLite等

#### 核心特性
- 自动迁移: _____
- 关联关系: _____
- 事务支持: _____
- 钩子函数: _____

#### 在项目中的应用
```go
// 示例：Model定义
type User struct {
    gorm.Model              // 包含ID, CreatedAt, UpdatedAt, DeletedAt
    Username  string `gorm:"uniqueIndex;not null"`
    Email     string `gorm:"uniqueIndex"`
    Password  string `gorm:"not null"`
    Nickname  string
    Avatar    string
}

// 示例：查询
db.Where("id = ?", id).First(&user)

// 关联查询
db.Preload("Resumes").Find(&user)

// 事务操作
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    return tx.Create(&resume).Error
})
```

#### 常用操作
- **创建**: `db.Create(&user)` - 插入新记录
- **查询**: 
  - `db.Find(&users)` - 查询所有
  - `db.Where("age > ?", 18).Find(&users)` - 条件查询
  - `db.First(&user)` - 查询第一条
- **更新**: 
  - `db.Save(&user)` - 保存所有字段
  - `db.Model(&user).Update("name", "新名字")` - 更新单个字段
- **删除**: 
  - `db.Delete(&user)` - 软删除
  - `db.Unscoped().Delete(&user)` - 硬删除

**高级特性**:
- **关联**: `db.Preload("Resumes").Find(&user)` - 预加载关联数据
- **钩子**: `BeforeCreate`、`AfterCreate` - 生命周期回调
- **事务**: `db.Transaction()` - 保证数据一致性
- **迁移**: `db.AutoMigrate(&User{})` - 自动同步表结构

**实际应用**:
- 用户CRUD操作
- 面试记录管理
- 简历信息存储
- 评估结果保存

---

### 4. JWT (JSON Web Token)

#### 基本信息
- **用途**: 身份认证和授权
- **格式**: Header.Payload.Signature

#### JWT结构
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.  # Header
eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6I...  # Payload
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_ad...  # Signature
```

#### 工作流程
1. **用户登录** → 服务器验证用户名密码 → 生成JWT Token
2. **客户端保存JWT** → 保存在localStorage或Cookie
3. **发送请求** → 在Authorization Header中携带Token
4. **服务器验证JWT** → 解析并验证签名和过期时间
5. **返回数据** → 验证通过后返回受保护的资源

#### JWT结构详解
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9  ← Header（算法和类型）
  .
eyJ1c2VyX2lkIjoxMjMsInVzZXJuYW1lIjoi...  ← Payload（用户数据）
  .
SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV...  ← Signature（签名）

Header: {"alg":"HS256","typ":"JWT"}
Payload: {"user_id":123,"username":"张三","exp":1234567890}
Signature: HMACSHA256(base64(header) + "." + base64(payload), secret)
```

#### 在项目中的应用
- **中间件**: `backend/internal/middleware/jwt.go`
- **生成Token**:
  ```go
  // 生成JWT Token
  func GenerateToken(userID uint, username, role string) (string, error) {
      claims := JWTClaims{
          UserID:   userID,
          Username: username,
          Role:     role,
          RegisteredClaims: jwt.RegisteredClaims{
              ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
              IssuedAt:  jwt.NewNumericDate(time.Now()),
              Issuer:    "interview-agent",
          },
      }
      
      token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      return token.SignedString([]byte(config.Global.Security.JWTSecret))
  }
  ```

- **验证Token**:
  ```go
  // JWT中间件验证
  func JWTMiddleware() app.HandlerFunc {
      return func(c context.Context, ctx *app.RequestContext) {
          // 1. 从Header提取Token
          tokenString := ctx.GetHeader("Authorization")
          
          // 2. 解析Token
          token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, 
              func(token *jwt.Token) (interface{}, error) {
                  return []byte(config.Global.Security.JWTSecret), nil
              })
          
          // 3. 验证Token有效性
          if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
              // 将用户信息存入上下文
              ctx.Set("user_id", claims.UserID)
              ctx.Set("username", claims.Username)
              ctx.Next(c)
          } else {
              ctx.JSON(401, gin.H{"error": "Unauthorized"})
              ctx.Abort()
          }
      }
  }
  ```

**白名单机制**:
```go
// 不需要JWT验证的路由
var authSkipList = []string{
    "/api/v1/user/login",
    "/api/v1/user/register",
    "/api/v1/user/wechat/login",
    "/api/v1/health",
}
```

**安全特性**:
- ✅ 无状态认证，不需要服务器端存储
- ✅ 防篡改（签名验证）
- ✅ 过期时间控制
- ✅ 支持刷新Token机制
- ⚠️ Token一旦签发，无法主动撤销（可通过黑名单机制解决）

---

## 数据存储技术

### 5. MySQL数据库

#### 选型原因
✅ **优势**:
- 成熟稳定: _____
- ACID支持: _____
- 生态丰富: _____

#### 在项目中的应用
- **存储内容**:
  - **用户数据**: 用户基本信息、认证信息、角色权限
  - **面试记录**: 面试会话、对话历史、面试配置
  - **简历信息**: 简历原文、解析结果、关键信息提取
  - **评估结果**: 面试评分、能力分析、改进建议
  - **题库数据**: 面试题目、参考答案、难度标签
  - **押题记录**: 预测题目、生成依据、准确率统计

#### 数据库设计
**主要表**:
- **users**: 用户ID、用户名、邮箱、密码（加密）、昵称、头像
- **user_models**: 用户自定义模型配置、API Key管理
- **resumes**: 简历ID、用户ID、文件路径、解析内容（JSON）、状态
- **interview_records**: 面试ID、用户ID、面试类型、开始/结束时间、状态、配置
- **interview_dialogues**: 对话ID、面试ID、角色（面试官/候选人）、内容、时间戳
- **interview_evaluations**: 评估ID、面试ID、总分、各维度评分、评语、建议
- **answer_reports**: 答案报告ID、对话ID、得分、分析、改进建议
- **prediction_records**: 押题记录ID、用户ID、简历ID、生成时间
- **prediction_questions**: 预测题目ID、记录ID、题目内容、依据、优先级

**关系图**: 
```
users (1) ────→ (n) resumes
  │
  ├────────────→ (n) interview_records
  │                   │
  │                   ├→ (n) interview_dialogues
  │                   │       │
  │                   │       └→ (n) answer_reports
  │                   │
  │                   └→ (n) interview_evaluations
  │
  └────────────→ (n) prediction_records
                      │
                      └→ (n) prediction_questions

resumes (1) ──→ (n) prediction_records
```

**索引设计**:
- `idx_user_id`: 用户ID索引（快速查询用户相关数据）
- `idx_interview_status`: 面试状态索引（查询进行中的面试）
- `idx_created_at`: 创建时间索引（按时间排序）
- `idx_resume_user`: 复合索引（resume_id + user_id）

**数据量估算**（假设1万用户）:
- users: ~10,000条
- resumes: ~30,000条（每人平均3份）
- interview_records: ~50,000条（每人平均5次）
- interview_dialogues: ~500,000条（每次面试平均10轮对话）
- interview_evaluations: ~50,000条（每次面试1份评估）

---

### 6. Redis缓存

#### 选型原因
✅ **优势**:
- 高性能: _____
- 数据结构丰富: _____
- 持久化支持: _____

#### 在项目中的应用
**用途**:
1. **缓存层**:
   - **用户信息缓存**: 减少数据库查询，提升响应速度
     ```
     Key: user:{user_id}
     Value: {"id":123,"username":"张三",...}
     TTL: 1小时
     ```
   - **面试会话缓存**: 缓存当前面试状态和上下文
     ```
     Key: interview:session:{interview_id}
     Value: {"current_question":5,"context":[...]}
     TTL: 面试结束后1小时
     ```
   - **简历缓存**: 缓存热门简历数据
   - **API响应缓存**: 缓存相同请求的响应
   
2. **消息队列（Pub/Sub）**:
   - **异步任务处理**: 避免阻塞主请求
     ```
     Channel: interview:messages:evaluation_report
     Message: {"type":"evaluation_report","interview_id":123}
     ```
   - **AI评估队列**: 面试结束后异步生成评估报告
   - **主题评估队列**: 对各个面试主题进行深度分析
   - **消息发布**:
     ```go
     client.Publish(ctx, "interview:messages:evaluation_report", message)
     ```
   - **消息订阅**:
     ```go
     pubsub := client.Subscribe(ctx, "interview:messages:*")
     for msg := range pubsub.Channel() {
         handleMessage(msg)
     }
     ```

3. **会话管理**:
   - **用户登录状态**: 存储在线用户
     ```
     Key: session:{session_id}
     Value: {"user_id":123,"login_time":...}
     TTL: 24小时
     ```
   - **面试进行中的状态**: 实时更新面试进度
   - **JWT黑名单**: 存储已登出的Token

4. **分布式锁**:
   - **防止并发创建**: 同一用户同时创建多个面试
   ```go
   lock, _ := client.SetNX(ctx, "lock:interview:{user_id}", 1, 10*time.Second)
   ```

5. **限流计数**:
   - **API限流**: 限制用户请求频率
   ```go
   count, _ := client.Incr(ctx, "rate_limit:{user_id}:{minute}")
   client.Expire(ctx, key, 60*time.Second)
   ```

#### 数据结构使用
- **String**: 简单的键值对存储
  - 用户Token: `SET token:{user_id} {token} EX 86400`
  - 计数器: `INCR view_count:{interview_id}`
  
- **Hash**: 存储对象的多个字段
  ```redis
  HSET user:123 name "张三" age 25 city "北京"
  HGET user:123 name
  ```
  - 用户信息、面试配置
  
- **List**: 有序列表，支持队列操作
  ```redis
  LPUSH task_queue {"task_id":1,...}
  RPOP task_queue
  ```
  - 任务队列、消息队列（简单场景）
  
- **Set**: 无序集合，支持交集、并集
  ```redis
  SADD online_users 123 456 789
  SMEMBERS online_users
  ```
  - 在线用户集合、标签集合
  
- **Sorted Set**: 有序集合，支持范围查询
  ```redis
  ZADD leaderboard 95 user:123 92 user:456
  ZRANGE leaderboard 0 10 WITHSCORES
  ```
  - 排行榜、定时任务调度

**缓存策略**:
- **缓存穿透**: 布隆过滤器 + 缓存空值
- **缓存击穿**: 分布式锁 + 热点数据永不过期
- **缓存雪崩**: 随机TTL + 多级缓存
- **缓存更新**: Cache-Aside模式（先更新DB，再删除缓存）

**性能优化**:
- 使用Pipeline批量操作
- 避免大Key（限制在10KB以内）
- 合理设置过期时间
- 使用连接池复用连接

---

### 7. Milvus向量数据库

#### 基本信息
- **类型**: 向量数据库
- **用途**: 向量相似度检索
- **是否必需**: ❌ 可选

#### 核心概念
- **向量**: _____
- **Embedding**: _____
- **相似度搜索**: _____

#### 在项目中的应用
- **位置**: `backend/internal/eino/milvus/`、`backend/chatApp/tool/milvus_retriever_tool.go`
- **用途**: 
  - **简历知识检索**: 根据JD（职位描述）找到最匹配的简历候选人
    - 存储：简历关键信息的向量表示
    - 查询：给定职位要求，找相似度最高的简历
  - **面试题库检索**: 根据简历内容和面试方向，检索相关题目
    - 存储：面试题目的向量表示
    - 查询：基于候选人背景，推荐合适难度的题目
  - **知识库问答**: 存储技术文档，实现智能问答
  - **相似对话检索**: 找到历史类似的面试对话，复用评估逻辑

#### 工作流程
```
【存储阶段】
原始文本 → Embedding模型 → 向量(1536维) → 存储到Milvus Collection

示例：
"精通Go并发编程，熟悉Goroutine和Channel"
     ↓ text-embedding-ada-002
[0.123, -0.456, 0.789, ...] (1536个数字)
     ↓
Milvus: collection="resumes", id=123

【检索阶段】
查询文本 → 向量化 → 相似度搜索 → 返回Top-K最相似文档

示例：
"招聘Go后端工程师，要求并发编程经验"
     ↓
向量化
     ↓
Milvus搜索（余弦相似度）
     ↓
返回最匹配的5份简历（相似度>0.8）
```

**核心概念**:
- **Collection（集合）**: 类似关系数据库的表
  - `resumes`: 存储简历向量
  - `interview_questions`: 存储面试题向量
  - `documents`: 存储技术文档向量
  
- **Embedding（向量嵌入）**: 
  - 将文本转换为固定维度的向量
  - 语义相近的文本，向量距离也近
  - 模型：text-embedding-ada-002（OpenAI）或国产模型
  - 维度：通常1536维或768维
  
- **相似度度量**:
  - **余弦相似度（COSINE）**: 推荐，范围[-1,1]，值越大越相似
  - **欧氏距离（L2）**: 直线距离，值越小越相似
  - **内积（IP）**: 向量点积，适合归一化向量

**配置示例**:
```yaml
Milvus:
  Address: "localhost:19530"
  CollectionName: "interview_knowledge"
  TopK: 5                    # 返回前5个最相似结果
  MetricType: "COSINE"       # 使用余弦相似度
  Dimensions: 1536           # 向量维度
```

#### 为什么是可选的
1. **不是核心功能**: 
   - 基础面试功能不依赖向量检索
   - 可以用传统数据库查询替代（虽然效果差一些）
   
2. **资源消耗大**:
   - 需要额外部署Milvus、Etcd、MinIO服务
   - 内存占用：向量数据常驻内存
   - 存储占用：1万条1536维向量约60MB
   
3. **成本考虑**:
   - Embedding API调用费用
   - 服务器资源成本
   
4. **渐进式启用**:
   - 初期：使用关键词匹配
   - 数据量增大后：启用向量检索提升准确率
   - 灵活性：可以随时开启/关闭

**何时需要启用Milvus**:
- ✅ 简历数量 > 1000份
- ✅ 需要语义搜索（不是简单关键词匹配）
- ✅ 需要智能推荐面试题
- ✅ 要实现知识库问答
- ✅ 追求更高的匹配准确率

**替代方案**:
- **小规模**: MySQL全文索引 + 关键词匹配
- **中规模**: Elasticsearch语义搜索
- **大规模**: Milvus专业向量数据库

---

## AI服务

### 8. OpenAI API

#### 支持的模型
- **GPT-4**: _____
- **GPT-3.5-turbo**: _____
- **text-embedding-ada-002**: _____

#### API调用方式
```go
// 使用Eino封装的OpenAI客户端
import (
    "github.com/cloudwego/eino-ext/components/model/openai"
    "github.com/cloudwego/eino/components/model"
)

// 创建聊天模型
chatModel, _ := openai.NewChatModel(ctx, &openai.ChatModelConfig{
    APIKey:      config.Global.OpenAI.APIKey,
    BaseURL:     config.Global.OpenAI.BaseURL,
    Model:       "gpt-4",
    Temperature: 0.7,
    MaxTokens:   2000,
})

// 发送聊天请求
resp, _ := chatModel.Generate(ctx, []*model.Message{
    model.SystemMessage("你是一位专业的面试官"),
    model.UserMessage("请出一道Go并发编程的面试题"),
})

fmt.Println(resp.Content)
```

#### 在项目中的应用
- **对话生成**: `chatApp/chat/openAi.go`
  - 面试官提问
  - 根据候选人回答生成追问
  - 生成面试总结和评语
  
- **文本理解**: 
  - 分析候选人回答的质量
  - 提取简历关键信息
  - 理解用户意图
  
- **Embedding生成**: 
  - 使用`text-embedding-ada-002`模型
  - 为简历、面试题生成向量
  - 用于Milvus相似度检索
  ```go
  embedding, _ := embeddingModel.Generate(ctx, "文本内容")
  // 返回1536维向量
  ```

- **结构化输出**:
  - 使用Function Calling提取结构化数据
  - 生成JSON格式的评估报告

#### Token和成本控制
- **Token计算**: 
  - 1 Token ≈ 4个字符（英文）或 1.5个汉字
  - 输入Token：用户消息 + 系统提示词 + 历史对话
  - 输出Token：AI生成的回复
  - 计算工具：tiktoken库
  
- **成本估算**（GPT-4价格）:
  ```
  输入：$0.03 / 1K tokens
  输出：$0.06 / 1K tokens
  
  单次面试（10轮对话）：
  - 输入：~5000 tokens × $0.03 = $0.15
  - 输出：~3000 tokens × $0.06 = $0.18
  - 总计：~$0.33 ≈ ¥2.4
  
  月成本（1000次面试）：¥2400
  ```
  
- **优化策略**:
  1. **模型选择**:
     - 简单任务用`gpt-3.5-turbo`（便宜10倍）
     - 复杂推理用`gpt-4`
     - 国产替代：文心一言、通义千问（更便宜）
  
  2. **提示词优化**:
     - 精简系统提示词
     - 避免冗余的示例
     - 使用Few-shot而非Many-shot
  
  3. **上下文管理**:
     - 只保留最近N轮对话
     - 总结旧对话，压缩上下文
     - 移除不必要的历史消息
  
  4. **缓存机制**:
     - 缓存常见问题的回答
     - 相似问题复用答案
  
  5. **流式输出**:
     - 使用Stream模式，用户体验更好
     - 可以提前终止，节省Token
  
  6. **批量处理**:
     - 多个任务合并一次请求
     - 减少API调用次数

**安全和限流**:
- 设置`max_tokens`限制单次输出长度
- 用户级别的速率限制
- 异常检测和熔断机制
- API Key轮换和负载均衡

---

## 前端技术栈

### 9. Next.js

#### 基本信息
- **类型**: React框架
- **版本**: 14+
- **特性**: 服务端渲染、静态生成

#### 在项目中的应用
- **位置**: `frontend/`
- **主要页面**:
  - `/login` - 登录页面
  - `/register` - 注册页面  
  - `/dashboard` - 用户控制台
  - `/interview` - 面试页面
  - `/resume` - 简历管理
  - `/history` - 面试历史
  - `/evaluation` - 评估报告
  
- **核心库**:
  - **Next.js 14**: App Router + Server Components
  - **React 18**: 用户界面库
  - **Ant Design**: UI组件库
  - **React Query**: 数据获取和缓存
  - **Zustand**: 状态管理（轻量级）
  - **Axios**: HTTP客户端
  - **TypeScript**: 类型安全
  
- **API调用**:
  ```typescript
  // services/api.ts
  import axios from 'axios'
  
  const api = axios.create({
    baseURL: 'http://localhost:8888/api/v1',
    timeout: 30000,
  })
  
  // 添加Token
  api.interceptors.request.use(config => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  })
  
  // 面试API
  export const interviewAPI = {
    start: (data) => api.post('/interview/start', data),
    answer: (data) => api.post('/interview/answer', data),
    getHistory: () => api.get('/interview/history'),
  }
  ```

#### 与后端的交互
```
【用户登录流程】
前端: 输入用户名密码
  ↓ POST /api/v1/user/login
后端: 验证密码 → 生成JWT
  ↓ {"token":"xxx", "user":{...}}
前端: 保存Token → 跳转Dashboard

【面试流程】
前端: 点击"开始面试"
  ↓ POST /api/v1/interview/start {type:"go"}
后端: 创建面试记录 → 生成第一个问题
  ↓ {"interview_id":123, "question":"..."}
前端: 显示问题 → 用户回答
  ↓ POST /api/v1/interview/answer {answer:"..."}
后端: Agent评估回答 → 生成下一题
  ↓ {"feedback":"...", "next_question":"..."}
前端: 显示反馈和下一题 → 循环...
  ↓ POST /api/v1/interview/end
后端: 生成评估报告（异步）
  ↓ {"status":"processing"}
前端: 轮询获取报告 → 显示结果
```

**前端架构**:
```
Pages (路由)
  ↓
Components (组件)
  ↓
Services (API调用)
  ↓
Store (状态管理)
```

**响应式设计**:
- 支持桌面端和移动端
- Tailwind CSS实现响应式布局
- Ant Design组件自适应

**性能优化**:
- Server Components减少客户端JS
- 图片优化（Next.js Image组件）
- 代码分割和懒加载
- React Query缓存减少请求

---

### 10. TypeScript


#### 在项目中的应用
- **类型定义**: `frontend/src/types/`
  ```typescript
  // types/user.ts
  export interface User {
    id: number
    username: string
    email: string
    nickname: string
    avatar?: string
    role: 'user' | 'admin'
    created_at: string
  }
  
  // types/interview.ts
  export interface Interview {
    id: number
    user_id: number
    type: 'go' | 'java' | 'mysql' | 'redis'
    status: 'pending' | 'in_progress' | 'completed'
    start_time: string
    end_time?: string
    config: InterviewConfig
  }
  
  export interface InterviewConfig {
    difficulty: 'junior' | 'intermediate' | 'senior'
    max_questions: number
    duration: number // 分钟
  }
  ```
  
- **API类型**:
  ```typescript
  // types/api.ts
  export interface ApiResponse<T = any> {
    code: number
    message: string
    data: T
  }
  
  export interface LoginRequest {
    username: string
    password: string
  }
  
  export interface LoginResponse {
    token: string
    user: User
  }
  
  export interface StartInterviewRequest {
    type: string
    resume_id?: number
    config?: Partial<InterviewConfig>
  }
  ```
  
- **组件Props**:
  ```typescript
  // components/InterviewCard.tsx
  interface InterviewCardProps {
    interview: Interview
    onStart?: (id: number) => void
    onView?: (id: number) => void
    className?: string
  }
  
  export const InterviewCard: React.FC<InterviewCardProps> = ({
    interview,
    onStart,
    onView,
    className
  }) => {
    // 组件实现
  }
  ```

**类型安全的好处**:
- ✅ 编译时捕获错误，减少运行时bug
- ✅ IDE智能提示和自动补全
- ✅ 重构更安全，影响范围清晰
- ✅ 代码即文档，类型定义就是最好的说明
- ✅ 团队协作更顺畅，减少沟通成本

**实际应用示例**:
```typescript
// 使用泛型确保类型安全
const useApi = <T,>(fn: () => Promise<ApiResponse<T>>) => {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)
  
  const execute = async () => {
    setLoading(true)
    try {
      const response = await fn()
      setData(response.data) // 类型为T
    } catch (err) {
      setError(err as Error)
    } finally {
      setLoading(false)
    }
  }
  
  return { data, loading, error, execute }
}

// 使用
const { data: interviews } = useApi<Interview[]>(interviewAPI.getHistory)
// interviews的类型自动推断为 Interview[] | null
```

---

## 开发工具

### 11. Docker

#### 用途
- **容器化部署**: 将应用打包成Docker镜像，实现"一次构建，到处运行"
- **环境一致性**: 开发、测试、生产环境完全一致，避免"我电脑上能跑"问题
- **服务编排**: 使用docker-compose管理多个服务（MySQL、Redis、Backend、Frontend）
- **资源隔离**: 每个服务独立运行，互不干扰
- **快速部署**: 一条命令启动整个系统

#### 在项目中的应用

**后端Dockerfile** (`backend/Dockerfile`):
```dockerfile
# 第一阶段：构建
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o main ./main.go

# 第二阶段：运行
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
EXPOSE 8888
CMD ["./main"]
```

**前端Dockerfile** (`frontend/Dockerfile`):
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./
EXPOSE 3000
CMD ["npm", "start"]
```

**docker-compose.yml** - 服务编排:
```yaml
services:
  # MySQL数据库
  mysql:
    image: mysql:8.0
    ports:
      - "3307:3306"
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: interview_agent
    volumes:
      - mysql-data:/var/lib/mysql
  
  # Redis缓存
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --requirepass root
  
  # 后端服务
  backend:
    build: ./backend
    ports:
      - "8888:8888"
    depends_on:
      - mysql
      - redis
    environment:
      - DATABASE_DSN=root:root@tcp(mysql:3306)/interview_agent
      - REDIS_ADDR=redis:6379
  
  # 前端服务
  frontend:
    build: ./frontend
    ports:
      - "3000:3000"
    depends_on:
      - backend

volumes:
  mysql-data:
  redis-data:
```

**常用命令**:
```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 停止服务
docker-compose down

# 重新构建
docker-compose build

# 进入容器
docker exec -it app-backend sh
```

**多阶段构建的优势**:
- ✅ 最终镜像只包含运行时需要的文件，体积更小
- ✅ 构建依赖不会进入生产镜像，更安全
- ✅ Go编译后的二进制文件+Alpine基础镜像，总大小仅~20MB

---

### 12. Git

#### 在项目中的应用
- **版本控制**: 
  - 记录代码变更历史
  - 随时回滚到任意版本
  - 代码审查和追溯
  
- **分支管理**:
  ```
  main (生产分支)
    ↓
  develop (开发分支)
    ↓
  feature/xxx (功能分支)
  bugfix/xxx (修复分支)
  hotfix/xxx (紧急修复)
  ```
  
- **协作开发**:
  - Pull Request代码审查
  - 冲突解决
  - Git Flow工作流
  
**Git工作流**:
```bash
# 1. 创建功能分支
git checkout -b feature/interview-agent

# 2. 开发并提交
git add .
git commit -m "feat: 添加面试Agent功能"

# 3. 推送到远程
git push origin feature/interview-agent

# 4. 创建PR，代码审查
# 5. 合并到develop分支
git checkout develop
git merge feature/interview-agent

# 6. 测试通过后，合并到main
git checkout main
git merge develop
git tag v1.0.0
```

**提交规范** (Conventional Commits):
```
feat: 新功能
fix: 修复bug
docs: 文档更新
style: 代码格式（不影响功能）
refactor: 重构
test: 测试相关
chore: 构建/工具配置

示例：
feat(interview): 添加Go专项面试功能
fix(auth): 修复JWT过期时间计算错误
docs(readme): 更新部署文档
```

---

## 技术栈协作图

```
用户请求
  ↓
Next.js前端
  ↓ HTTP
Hertz API网关
  ↓
Go业务逻辑 + Eino AI
  ↓              ↓
MySQL          OpenAI API
Redis
Milvus (可选)
```

---

## 技术选型总结

### 为什么选择Go语言
1. **高性能**: 
   - 编译型语言，性能接近C/C++
   - 并发模型（Goroutine）轻量高效
   - 适合高并发场景（面试系统需要同时处理多个会话）
   
2. **开发效率**:
   - 语法简洁，学习曲线平缓
   - 标准库丰富，"电池自带"
   - 编译快，开发调试效率高
   
3. **工程化**:
   - 强类型，编译时发现错误
   - 内置测试和性能分析工具
   - 代码格式化（gofmt）统一风格
   - 依赖管理（go mod）简单可靠
   
4. **部署优势**:
   - 编译成单一二进制文件
   - 跨平台编译
   - 无运行时依赖，部署简单
   - Docker镜像体积小（20MB左右）
   
5. **生态成熟**:
   - 云原生领域事实标准（Kubernetes、Docker、Etcd）
   - 微服务框架丰富（gRPC、Hertz、Gin）
   - AI工具链完善（Eino、LangChain-Go）

### 为什么选择字节系技术栈
1. **Hertz + Eino都是字节开源**: 
   - 同一团队维护，兼容性好
   - API设计风格一致
   - 互相深度集成，减少适配工作
   - 文档和示例丰富
   
2. **技术栈一致性**: 
   - 避免多语言混合开发的复杂性
   - 团队技能栈统一，降低学习成本
   - 代码风格和最佳实践一致
   - 工具链统一（调试、监控、日志）
   
3. **企业级验证**: 
   - 字节跳动内部大规模使用
   - 经过抖音、今日头条等产品验证
   - 性能和稳定性有保障
   - 问题快速响应和修复
   
4. **性能优势**:
   - Hertz比Gin性能高40%+
   - Eino比LangChain（Python）快数倍
   - 适合高并发AI应用场景

### 技术栈的优缺点

**优点**:
- ✅ **高性能**: Go + Hertz + Eino组合，性能卓越
- ✅ **类型安全**: Go静态类型 + TypeScript前端，减少错误
- ✅ **易于部署**: Docker容器化，一键部署
- ✅ **开发效率**: 代码简洁，工具链完善
- ✅ **可维护性**: 技术栈统一，代码风格一致
- ✅ **成本可控**: 单体服务，资源占用低
- ✅ **扩展性好**: 模块化设计，易于扩展
- ✅ **社区支持**: 字节系开源社区活跃

**缺点**:
- ❌ **生态相对小**: Hertz/Eino社区不如Gin/LangChain大
- ❌ **文档中文为主**: 国际化支持一般
- ❌ **学习资源**: 相关教程和案例相对较少
- ❌ **招聘难度**: 熟悉字节系技术栈的人才相对稀缺
- ❌ **版本迭代快**: API可能有breaking changes
- ❌ **AI成本**: OpenAI API调用费用较高
- ❌ **单点依赖**: 过度依赖OpenAI，需要有降级方案

---

## ❓ 问题记录

### Q1: 为什么不用Python + LangChain？

**回答**: Python + LangChain是更主流的选择，但本项目选择Go + Eino有以下考虑：

**性能对比**:
- Go编译型语言，性能是Python的3-10倍
- Goroutine并发模型，支持高并发场景
- 单个面试会话响应时间：Go ~200ms, Python ~500ms+
- 内存占用：Go ~50MB, Python ~200MB+

**开发效率**:
- Go语法简洁，代码量少
- 静态类型减少运行时错误
- 编译快，开发调试效率高
- 不需要处理Python的依赖地狱

**部署优势**:
- Go编译成单一二进制文件，部署简单
- Python需要安装解释器+虚拟环境+依赖包
- Docker镜像：Go ~20MB, Python ~500MB+
- 启动速度：Go ~100ms, Python ~3s+

**技术栈统一**:
- 后端已使用Go（Hertz框架）
- 避免引入Python增加复杂性
- 团队技能栈统一

**何时应该用Python + LangChain**:
- ✅ 快速原型验证（Python开发更快）
- ✅ 团队Python技能更强
- ✅ 需要LangChain丰富的生态（如更多的Tool集成）
- ✅ 性能要求不高的场景
- ✅ 需要Jupyter Notebook进行实验

### Q2: Milvus可以用其他向量数据库替代吗？

**回答**: 可以！向量数据库有多种选择，各有优劣：

**1. Milvus**（当前选择）
- ✅ 开源免费
- ✅ 性能强大，支持十亿级向量
- ✅ 功能丰富（多种索引、多种距离度量）
- ❌ 部署复杂（需要Etcd、MinIO等依赖）
- ❌ 资源占用大

**2. Qdrant**
- ✅ Rust编写，性能优秀
- ✅ 部署简单（单一二进制文件）
- ✅ RESTful API，易于集成
- ❌ 生态相对较小

**3. Weaviate**
- ✅ GraphQL API，查询灵活
- ✅ 内置ML模型，可以直接处理文本
- ❌ 资源占用大

**4. Pinecone**（商业）
- ✅ 托管服务，无需运维
- ✅ 性能和稳定性好
- ❌ 收费（免费版有限制）
- ❌ 数据在第三方

**5. Elasticsearch（向量搜索插件）**
- ✅ 如果已有ES，可以复用
- ✅ 全文检索+向量检索一体
- ❌ 向量搜索性能不如专业向量数据库

**6. PostgreSQL + pgvector**
- ✅ 如果已有PostgreSQL，可以复用
- ✅ 学习成本低
- ❌ 性能较差，适合小规模（< 10万向量）

**选择建议**:
- 小规模（< 10万）：PostgreSQL + pgvector
- 中规模（10万-100万）：Qdrant或Weaviate
- 大规模（> 100万）：Milvus
- 不想运维：Pinecone（付费）

**迁移成本**:
- 向量数据库接口相对标准化
- Eino支持多种向量数据库
- 迁移主要是改配置，代码改动很小

### Q3: 为什么要用消息队列？直接同步处理不行吗？

**回答**: 消息队列带来的好处：

**1. 异步处理，提升响应速度**
```
同步方式：
用户结束面试 → 生成评估报告（耗时30s） → 返回结果
用户等待时间：30s+

异步方式：
用户结束面试 → 发送消息到队列 → 立即返回
后台消费者 → 生成评估报告 → 更新数据库
用户等待时间：< 1s，轮询获取报告
```

**2. 削峰填谷**
- 高峰期：消息暂存队列，避免系统过载
- 低峰期：慢慢消费消息，充分利用资源

**3. 解耦服务**
- 面试模块不需要知道评估模块的实现
- 评估模块可以独立扩展、升级
- 任一模块故障不影响对方

**4. 失败重试**
- AI调用失败可以自动重试
- 保证消息最终被处理

**5. 多消费者**
- 评估报告生成
- 主题评估分析
- 用户行为分析
- 数据统计

**何时不需要消息队列**:
- ✅ 用户量很小（< 100人）
- ✅ 可以接受等待（如批量任务）
- ✅ 没有异步需求
- ✅ 想要简化架构

**为什么用Redis Pub/Sub而不是RabbitMQ/Kafka**:
见前文Redis部分的详细说明。简单说：Redis已经在用了，小规模够用，部署简单。

---

## 📚 学习资源

### 官方文档
- [ ] [Hertz文档](https://www.cloudwego.io/zh/docs/hertz/)
- [ ] [Eino文档](https://www.cloudwego.io/zh/docs/eino/)
- [ ] [GORM文档](https://gorm.io/zh_CN/docs/)
- [ ] [Next.js文档](https://nextjs.org/docs)

---

## ✅ 学习检查点

- [ ] 理解每个技术组件的作用
- [ ] 知道技术选型的原因
- [ ] 理解各技术之间的协作方式
- [ ] 掌握核心技术的基本概念
- [ ] 能够解释为什么选择这套技术栈

---

**第一阶段完成！下一步**: [第二阶段：核心模块](../02-data-layer/README.md)
