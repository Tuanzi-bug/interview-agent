# 配置管理

> 深入理解系统的配置管理机制

## 📋 学习任务

1. 阅读 [backend/internal/config/config.go](../../backend/internal/config/config.go)
2. 分析 [backend/config.yaml](../../backend/config.yaml)
3. 理解配置加载和环境变量扩展机制
4. 掌握配置项的作用和默认值

---

## 配置结构分析

### Config结构体

**代码位置**: `backend/internal/config/config.go`

```go
// 粘贴Config结构体定义
type Config struct {
    Host     string          `yaml:"host"`
    Port     int             `yaml:"port"`
    // ... 补充其他字段
}
```

### 配置项分类

| 分类 | 结构体 | 作用 |
|------|--------|------|
| 服务配置 | Host, Port | _____ |
| 数据库配置 | DatabaseConfig | _____ |
| Redis配置 | RedisConfig | _____ |
| Hertz配置 | HertzConfig | _____ |
| AI配置 | OpenAIConfig, EinoConfig | _____ |
| 安全配置 | SecurityConfig | _____ |
| 微信配置 | WechatConfig | _____ |
| 其他 | _____ | _____ |

---

## 详细配置项说明

### 1. 服务配置

```yaml
host: 0.0.0.0
port: 8080
```

**说明**:
- `host`: _____
- `port`: _____

**我的理解**: _____

---

### 2. 数据库配置 (DatabaseConfig)

```yaml
database:
  driver: mysql
  dsn: ${DATABASE_DSN}
  max_open_conns: 100
  max_idle_conns: 10
  conn_max_lifetime: 1h
```

**字段说明**:

| 字段 | 类型 | 作用 | 默认值 |
|------|------|------|--------|
| driver | string | _____ | _____ |
| dsn | string | _____ | _____ |
| max_open_conns | int | _____ | _____ |
| max_idle_conns | int | _____ | _____ |
| conn_max_lifetime | duration | _____ | _____ |

**连接字符串格式**:
```
user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local
```

**我的理解**: _____

---

### 3. Redis配置 (RedisConfig)

```yaml
redis:
  addr: ${REDIS_ADDR}
  password: ${REDIS_PASSWORD}
  db: 0
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  pool_size: 10
  min_idle_conns: 5
```

**字段说明**:

| 字段 | 类型 | 作用 | 说明 |
|------|------|------|------|
| addr | string | _____ | _____ |
| password | string | _____ | _____ |
| db | int | _____ | _____ |
| pool_size | int | _____ | _____ |

**我的理解**: _____

---

### 4. Hertz配置 (HertzConfig)

```yaml
hertz:
  log_level: info
  log_path: ./logs
  read_timeout: 60s
  write_timeout: 60s
```

**字段说明**:
- `log_level`: _____
- `log_path`: _____
- `read_timeout`: _____
- `write_timeout`: _____

**超时时间的影响**: _____

**我的理解**: _____

---

### 5. OpenAI配置 (OpenAIConfig)

```yaml
openai:
  api_key: ${OPENAI_API_KEY}
  base_url: https://api.openai.com/v1
  model: gpt-4
  temperature: 0.7
  max_tokens: 2000
```

**字段说明**:
- `api_key`: _____
- `base_url`: _____
- `model`: _____（支持的模型: _____）
- `temperature`: _____（范围: _____）
- `max_tokens`: _____

**我的理解**: _____

---

### 6. Embedding配置 (EmbeddingConfig)

```yaml
Embedding:
  model: text-embedding-ada-002
  dimension: 1536
```

**字段说明**:
- `model`: _____
- `dimension`: _____

**向量维度的意义**: _____

**我的理解**: _____

---

### 7. Milvus配置 (MilvusConfig)

```yaml
Milvus:
  enabled: false
  address: localhost:19530
  collection_name: interview_knowledge
  dimension: 1536
```

**字段说明**:
- `enabled`: _____
- `address`: _____
- `collection_name`: _____
- `dimension`: _____

**为什么默认disabled**: _____

**我的理解**: _____

---

### 8. 面试配置 (InterviewConfig)

```yaml
interview:
  max_duration: 60
  max_questions: 20
  default_difficulty: intermediate
```

**字段说明**:
- `max_duration`: _____（单位: _____）
- `max_questions`: _____
- `default_difficulty`: _____（可选值: _____）

**这些限制的意义**: _____

**我的理解**: _____

---

### 9. 安全配置 (SecurityConfig)

```yaml
security:
  jwt_secret: ${JWT_SECRET}
  jwt_expire: 24h
  cors:
    allow_origins: ["*"]
    allow_methods: ["GET", "POST", "PUT", "DELETE"]
    allow_headers: ["*"]
```

**JWT配置**:
- `jwt_secret`: _____
- `jwt_expire`: _____

**CORS配置**:
- 作用: _____
- 生产环境建议: _____

**我的理解**: _____

---

### 10. 微信配置 (WechatConfig)

```yaml
wechat:
  app_id: ${WECHAT_APP_ID}
  app_secret: ${WECHAT_APP_SECRET}
  redirect_url: http://localhost:3000/callback
```

**字段说明**:
- `app_id`: _____
- `app_secret`: _____
- `redirect_url`: _____

**OAuth2.0流程**: _____

**我的理解**: _____

---

## 配置加载机制

### LoadConfig函数分析

**代码位置**: `backend/internal/config/config.go`

```go
// 粘贴LoadConfig函数代码
func LoadConfig(configPath string) (*Config, error) {
    // ...
}
```

**加载步骤**:
1. _____
2. _____
3. _____

**错误处理**: _____

**我的理解**: _____

---

## 环境变量扩展

### ExpandEnv机制

**代码位置**: `backend/internal/config/config.go`

```go
// 粘贴ExpandEnv相关代码
func (c *Config) ExpandEnv() {
    // ...
}
```

**语法**: `${VAR_NAME}` 或 `${VAR_NAME:default}`

**示例**:
```yaml
database:
  dsn: ${DATABASE_DSN:root:password@tcp(localhost:3306)/interview}
```

**工作原理**: _____

**我的理解**: _____

---

## 环境变量清单

### 必需的环境变量

| 环境变量 | 说明 | 示例值 | 是否必需 |
|---------|------|--------|---------|
| OPENAI_API_KEY | OpenAI API密钥 | sk-xxx | ✅ |
| DATABASE_DSN | 数据库连接串 | root:pass@tcp(localhost:3306)/db | ✅ |
| JWT_SECRET | JWT签名密钥 | your-secret-key | ✅ |
| REDIS_ADDR | Redis地址 | localhost:6379 | ✅ |
| REDIS_PASSWORD | Redis密码 | _____ | ❌ |
| WECHAT_APP_ID | 微信AppID | _____ | ❌ |
| WECHAT_APP_SECRET | 微信AppSecret | _____ | ❌ |

### .env文件模板

```bash
# OpenAI配置
OPENAI_API_KEY=sk-your-api-key-here

# 数据库配置
DATABASE_DSN=root:password@tcp(localhost:3306)/interview_db?charset=utf8mb4&parseTime=True&loc=Local

# Redis配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# 安全配置
JWT_SECRET=your-very-secret-key-here

# 微信配置（可选）
WECHAT_APP_ID=
WECHAT_APP_SECRET=
```

---

## 配置验证

### 配置检查清单

- [ ] 所有必需的环境变量都已设置
- [ ] 数据库连接信息正确
- [ ] Redis连接信息正确
- [ ] OpenAI API Key有效
- [ ] JWT Secret足够安全（至少32字符）
- [ ] CORS配置符合生产环境要求
- [ ] 日志路径有写入权限
- [ ] 超时时间设置合理

### 配置测试

```go
// 编写一个简单的配置测试
func TestConfig() {
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        // 处理错误
    }
    // 验证配置项
}
```

---

## 不同环境的配置

### 开发环境

```yaml
# config.dev.yaml
hertz:
  log_level: debug
  
openai:
  temperature: 0.9  # 更随机，方便测试
```

### 生产环境

```yaml
# config.prod.yaml
hertz:
  log_level: warn
  
security:
  cors:
    allow_origins: ["https://yourdomain.com"]
```

### 测试环境

```yaml
# config.test.yaml
database:
  dsn: root:password@tcp(localhost:3306)/interview_test
```

**配置选择方式**: _____

**我的理解**: _____

---

## 配置最佳实践

### 安全性

1. **敏感信息使用环境变量**
   - ✅ 正确: `api_key: ${OPENAI_API_KEY}`
   - ❌ 错误: `api_key: sk-xxx直接写在配置文件`

2. **不要提交.env文件**
   - 添加到.gitignore
   - 提供.env.example模板

3. **JWT Secret强度**
   - 至少32字符
   - 包含大小写字母、数字、特殊字符

### 性能优化

1. **数据库连接池**
   ```yaml
   max_open_conns: 100  # 根据服务器性能调整
   max_idle_conns: 10   # 保持一定空闲连接
   ```

2. **Redis连接池**
   ```yaml
   pool_size: 10        # 根据并发量调整
   ```

3. **超时时间**
   ```yaml
   read_timeout: 60s    # 考虑AI响应时间
   write_timeout: 60s
   ```

---

## ❓ 问题记录

### Q1: 为什么某些配置使用环境变量，某些直接写在yaml？
**我的理解**: _____

### Q2: 如何动态更新配置而不重启服务？
**我的理解**: _____

### Q3: _____
**我的理解**: _____

---

## 🔬 实践验证

### 任务1: 修改配置并测试
- [ ] 修改OpenAI的temperature参数
- [ ] 观察AI响应的变化
- [ ] 记录你的发现: _____

### 任务2: 环境变量测试
- [ ] 创建.env文件
- [ ] 设置环境变量
- [ ] 验证系统能正确读取
- [ ] 记录过程: _____

### 任务3: 配置验证
- [ ] 故意配置错误的值（如负数的端口）
- [ ] 观察系统的错误处理
- [ ] 记录你的发现: _____

---

## ✅ 学习检查点

- [ ] 理解所有配置项的作用
- [ ] 掌握环境变量扩展机制
- [ ] 能够为不同环境编写配置文件
- [ ] 理解配置的安全最佳实践
- [ ] 知道如何排查配置相关问题
- [ ] 实际配置过开发环境

---

**下一步**: [技术栈详解](tech-stack.md) 或 [第二阶段：核心模块](../02-data-layer/README.md)
