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
// Config 应用程序配置结构
type Config struct {
    Host             string          `yaml:"host"`              // 服务监听地址
    Port             int             `yaml:"port"`              // 服务监听端口
    Database         DatabaseConfig  `yaml:"database"`          // 数据库配置
    Redis            RedisConfig     `yaml:"redis"`             // Redis缓存配置
    Hertz            HertzConfig     `yaml:"hertz"`             // Hertz框架配置
    Eino             EinoConfig      `yaml:"eino"`              // Eino框架配置
    Interview        InterviewConfig `yaml:"interview"`         // 面试系统配置
    Security         SecurityConfig  `yaml:"security"`          // 安全性配置
    GoogleSearch     GoogleConfig    `yaml:"google_search"`     // Google搜索配置
    OpenAI           OpenAIConfig    `yaml:"openai"`            // OpenAI配置
    Embedding        EmbeddingConfig `yaml:"Embedding"`         // Embedding服务配置
    Milvus           MilvusConfig    `yaml:"Milvus"`            // Milvus向量数据库配置
    DocumentSplitter SplitterConfig  `yaml:"DocumentSplitter"`  // 文档分割器配置
    Wechat           WechatConfig    `yaml:"wechat"`            // 微信配置
    Feishu           FeishuConfig    `yaml:"feishu"`            // 飞书配置
}
```

### 配置项分类

| 分类 | 结构体 | 作用 |
|------|--------|------|
| 服务配置 | Host, Port | 定义HTTP服务的监听地址和端口 |
| 数据库配置 | DatabaseConfig | MySQL数据库连接和连接池配置 |
| Redis配置 | RedisConfig | Redis缓存服务连接和性能配置 |
| Hertz配置 | HertzConfig | HTTP框架的日志、超时等配置 |
| AI基础配置 | OpenAIConfig, EinoConfig | AI模型的基础连接配置 |
| 向量服务配置 | EmbeddingConfig, MilvusConfig | 向量嵌入和向量数据库配置 |
| 面试配置 | InterviewConfig | 面试流程的时间和题目数量限制 |
| 安全配置 | SecurityConfig | JWT认证和CORS跨域配置 |
| 第三方集成 | WechatConfig, FeishuConfig, GoogleConfig | 微信、飞书、Google搜索集成 |
| 工具配置 | SplitterConfig | 文档分割器的参数配置 |

---

## 详细配置项说明

### 1. 服务配置

```yaml
host: 0.0.0.0
port: 8080
```

**说明**:
- `host`: 服务监听的IP地址，`0.0.0.0` 表示监听所有网络接口，允许外部访问
- `port`: 服务监听的端口号，默认8888，可通过环境变量修改

**生产环境建议**:
- 开发环境：使用 `0.0.0.0` 方便调试
- 生产环境：考虑使用 `127.0.0.1` 配合反向代理（如Nginx）提高安全性

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
| driver | string | 数据库驱动类型 | mysql |
| dsn | string | 数据库连接字符串（Data Source Name） | 见配置文件 |
| max_open_conns | int | 最大打开连接数（并发连接上限） | 100 |
| max_idle_conns | int | 最大空闲连接数（连接池保持的空闲连接） | 10 |
| conn_max_lifetime | duration | 连接最大存活时间（避免长连接问题） | 1h |

**连接字符串格式**:
```
user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local

示例: root:root@tcp(mysql:3306)/interview_agent?charset=utf8mb4&parseTime=True&loc=Local
```

**参数说明**:
- `charset=utf8mb4`: 使用utf8mb4字符集，支持emoji等4字节字符
- `parseTime=True`: 自动将MySQL的TIME/DATE类型转换为Go的time.Time
- `loc=Local`: 使用本地时区

**性能调优建议**:
- `max_open_conns`: 根据并发量和数据库性能设置，通常10-200之间
- `max_idle_conns`: 保持一定数量避免频繁创建连接，通常为max_open_conns的10-20%
- `conn_max_lifetime`: 避免长时间连接导致的问题，建议30分钟到2小时

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
| addr | string | Redis服务器地址 | 格式: host:port，如 redis:6379 |
| password | string | Redis认证密码 | 如无密码可留空 |
| db | int | Redis数据库索引 | 0-15，默认0 |
| dial_timeout | duration | 连接超时时间 | 建立连接的最长等待时间 |
| read_timeout | duration | 读取超时时间 | 从Redis读取数据的超时 |
| write_timeout | duration | 写入超时时间 | 向Redis写入数据的超时 |
| pool_size | int | 连接池大小 | 最大连接数，默认10 |
| min_idle_conns | int | 最小空闲连接数 | 连接池保持的最小空闲连接 |

**Redis在系统中的作用**:
- 用户会话管理（Session缓存）
- 面试进度状态缓存
- 热点数据缓存（减少数据库压力）
- 分布式锁（如面试并发控制）

**性能建议**:
- `pool_size`: 根据并发量调整，通常10-50
- `min_idle_conns`: 保持5-10个避免冷启动
- 超时时间不宜过长，避免阻塞其他请求

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
- `log_level`: 日志级别（debug/info/warn/error），控制日志输出详细程度
- `log_path`: 日志文件保存路径，默认./logs目录
- `read_timeout`: 读取请求体的超时时间，防止慢客户端攻击
- `write_timeout`: 写入响应的超时时间
- `idle_timeout`: 空闲连接的超时时间，Keep-Alive连接的最大空闲时间

**日志级别选择**:
- **debug**: 开发环境，输出详细调试信息
- **info**: 生产环境推荐，记录重要操作
- **warn**: 只记录警告和错误
- **error**: 只记录错误

**超时时间的影响**:
- **过短**: 可能导致正常请求被中断（特别是AI调用耗时较长）
- **过长**: 可能导致资源占用过多，降低系统并发能力
- **推荐值**: read_timeout=10s, write_timeout=10s, idle_timeout=60s
- **AI接口**: 由于AI响应可能较慢，建议适当延长超时时间

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
- `api_key`: OpenAI API密钥，从OpenAI平台获取，必须保密
- `base_url`: API基础URL，默认为官方地址，支持自定义代理或第三方兼容服务
- `model_name`: 使用的模型名称（如gpt-4, gpt-3.5-turbo等）

**支持的模型**:
- **gpt-4**: 最强大，适合复杂推理和代码生成
- **gpt-3.5-turbo**: 性价比高，适合常规对话
- **gpt-4-turbo**: 速度更快，成本更低

**注意事项**:
- API Key必须通过环境变量配置，不要直接写在配置文件中
- base_url可用于配置国内镜像或代理服务
- 不同模型的定价和速度差异较大，需根据业务需求选择

**在Eino配置中有更详细的AI参数配置**（如temperature、max_tokens等）

---

### 6. Embedding配置 (EmbeddingConfig)

```yaml
Embedding:
  model: text-embedding-ada-002
  dimension: 1536
```

**字段说明**:
- `APIKey`: API密钥认证（与AccessKey/SecretKey二选一）
- `AccessKey/SecretKey`: AK/SK认证方式（与APIKey二选一）
- `Model`: Embedding模型的端点ID，从Ark平台获取
- `BaseURL`: Embedding服务的API地址
- `Region`: 服务区域（如cn-beijing）
- `Timeout`: 请求超时时间，默认30秒
- `RetryTimes`: 失败重试次数，默认3次
- `Dimensions`: 输出向量的维度，必须与模型实际输出匹配
- `User`: 用户标识，可选，用于请求追踪

**向量维度的意义**:
- 向量维度决定了文本表示的精细程度
- 维度越高，表示能力越强，但计算和存储成本也越高
- 常见维度: 768(BERT), 1536(text-embedding-ada-002), 2560(某些国产模型)
- **重要**: Dimensions必须与选择的模型实际输出维度一致

**支持环境变量**:
所有敏感字段都支持 `${VAR_NAME}` 语法从环境变量读取

---

### 7. Milvus配置 (MilvusConfig)

```yaml
#Milvus:
#  # 连接配置
#  Address: "${MILVUS_ADDRESS}"              # Milvus 服务地址，从环境变量读取（.env 中设置）
#  Username: "${MILVUS_USERNAME}"            # 从环境变量 MILVUS_USERNAME 读取
#  Password: "${MILVUS_PASSWORD}"            # 从环境变量 MILVUS_PASSWORD 读取
#  DatabaseName: "default"                   # 数据库名称
#  CollectionName: "documents"               # 集合名称
#
#  # 检索配置
#  TopK: 5                                   # 返回的最相似文档数量
#  MetricType: "COSINE"                      # 距离度量类型: L2, IP, COSINE
#
#  # 超时配置
#  ConnectTimeout: 10s                       # 连接超时
#  SearchTimeout: 30s                        # 搜索超时
```

**字段说明**:
- `Address`: Milvus服务地址（格式: host:port）
- `Username/Password`: 认证信息（可选）
- `DatabaseName`: 数据库名称，默认"default"
- `CollectionName`: 默认集合名称，用于存储向量
- `Collections`: 多集合映射配置，支持不同类型数据使用不同集合
- `TopK`: 检索时返回的最相似文档数量
- `MetricType`: 距离度量方式（L2/IP/COSINE）
- `ConnectTimeout/SearchTimeout`: 连接和搜索超时配置

**距离度量类型**:
- **L2**: 欧几里得距离，值越小越相似
- **IP**: 内积（点积），值越大越相似
- **COSINE**: 余弦相似度，值越大越相似（推荐）

**为什么配置文件中注释掉**:
- Milvus是可选组件，不是系统必需的
- 对于小规模应用，可以不使用向量数据库
- 需要额外部署Milvus服务，增加了部署复杂度
- 开发环境可以先不启用，生产环境根据需要配置

**GetCollection方法**: 支持多集合管理，可以根据名称获取不同的集合

---

### 8. 面试配置 (InterviewConfig)

```yaml
interview:
  max_duration: 60
  max_questions: 20
  default_difficulty: intermediate
```

**字段说明**:
- `max_duration`: 单次面试最大时长（格式: 30m表示30分钟）
- `question_timeout`: 单个问题的超时时间（10m表示10分钟）
- `max_questions`: 单次面试最多问题数量（默认10题）
- `min_questions`: 单次面试最少问题数量（默认5题）

**时间格式说明**:
- 使用Go的Duration格式: 1h（1小时）、30m（30分钟）、90s（90秒）
- 可以组合使用: 1h30m（1小时30分钟）

**这些限制的意义**:
1. **用户体验**: 避免面试过长导致疲劳
2. **资源控制**: 限制AI调用次数和时长，控制成本
3. **质量保证**: 合理的题目数量保证面试质量
4. **超时保护**: 防止面试进程hang住占用资源

**推荐配置**:
- 快速面试: max_duration=15m, max_questions=5
- 标准面试: max_duration=30m, max_questions=10
- 深度面试: max_duration=60m, max_questions=20

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
- `jwt_secret`: JWT签名密钥，必须保密且足够复杂（建议32字符以上）
- `jwt_expiration`: JWT令牌有效期（24h表示24小时）

**JWT工作原理**:
1. 用户登录成功后，服务器生成JWT令牌
2. 令牌包含用户信息，使用jwt_secret签名
3. 客户端将令牌保存在本地（localStorage或cookie）
4. 后续请求在Header中携带令牌（Authorization: Bearer <token>）
5. 服务器验证签名和有效期，识别用户身份

**CORS配置**:
- **作用**: 跨域资源共享，允许前端从不同域名访问后端API
- **allow_origins**: 允许的源域名，"*"表示允许所有（开发环境方便，生产环境不安全）
- **allow_methods**: 允许的HTTP方法
- **allow_headers**: 允许的请求头
- **allow_credentials**: 是否允许携带Cookie

**生产环境建议**:
- `jwt_secret`: 使用随机生成的强密钥，定期轮换
- `jwt_expiration`: 根据安全要求设置，建议不超过24小时
- `allow_origins`: 明确指定前端域名，不要使用"*"
  ```yaml
  allow_origins: ["https://yourdomain.com", "https://www.yourdomain.com"]
  ```
- `allow_credentials`: 如果需要携带Cookie，必须明确指定origins，不能用"*"

---

### 10. 微信配置 (WechatConfig)

```yaml
wechat:
  app_id: ${WECHAT_APP_ID}
  app_secret: ${WECHAT_APP_SECRET}
  redirect_url: http://localhost:3000/callback
```

**字段说明**:
- `app_id`: 微信开放平台应用ID
- `app_secret`: 微信开放平台应用密钥
- `redirect_url`: OAuth2.0回调地址，用户授权后跳转的URL

**微信登录OAuth2.0流程**:
1. 用户点击"微信登录"按钮
2. 前端跳转到微信授权页面（带上app_id和redirect_url）
3. 用户在微信授权页面确认授权
4. 微信重定向回redirect_url，带上授权code
5. 后端使用code + app_id + app_secret 换取access_token
6. 使用access_token获取用户信息
7. 后端创建或更新用户，生成JWT令牌返回前端

**配置说明**:
- app_id和app_secret需要在微信开放平台申请
- redirect_url必须在微信开放平台配置的白名单中
- 本地开发可以使用内网穿透工具（如ngrok）

**安全注意**:
- app_secret必须保密，不能暴露给前端
- 生产环境使用HTTPS协议
- redirect_url要严格校验，防止重定向攻击

---

## 配置加载机制

### LoadConfig函数分析

**代码位置**: `backend/internal/config/config.go`

```go
// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
    // 1. 读取配置文件
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    // 2. 解析YAML到Config结构体
    var cfg Config
    err = yaml.Unmarshal(data, &cfg)
    if err != nil {
        return nil, err
    }

    // 3. 保存到全局变量（方便其他模块访问）
    Global = cfg
    log.Println("配置加载成功")
    return &cfg, nil
}
```

**加载步骤**:
1. **读取文件**: 使用`os.ReadFile`读取YAML配置文件内容
2. **解析YAML**: 使用`yaml.Unmarshal`将YAML内容反序列化到Config结构体
3. **保存全局**: 将配置保存到`Global`变量，供全局访问
4. **返回配置**: 返回配置指针和可能的错误

**错误处理**:
- 文件不存在: 返回"no such file or directory"错误
- YAML格式错误: 返回"yaml: unmarshal errors"错误
- 建议在main函数中处理配置加载失败，直接退出程序

**使用示例**:
```go
// 在main.go中加载配置
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
    log.Fatalf("加载配置失败: %v", err)
}

// 展开环境变量
cfg.ExpandEnv()

// 访问配置
port := cfg.Port
dbDSN := cfg.Database.DSN
```

---

## 环境变量扩展

### ExpandEnv机制

**代码位置**: `backend/internal/config/config.go`

```go
// ExpandEnv 展开配置中的环境变量引用
// 支持 ${VAR_NAME} 和 $VAR_NAME 两种语法
func (c *Config) ExpandEnv() {
    // 展开 Embedding 配置
    c.Embedding.APIKey = expandEnvVar(c.Embedding.APIKey)
    c.Embedding.AccessKey = expandEnvVar(c.Embedding.AccessKey)
    c.Embedding.SecretKey = expandEnvVar(c.Embedding.SecretKey)
    c.Embedding.Model = expandEnvVar(c.Embedding.Model)
    c.Embedding.BaseURL = expandEnvVar(c.Embedding.BaseURL)
    c.Embedding.Region = expandEnvVar(c.Embedding.Region)
    
    // 展开 Milvus 配置
    c.Milvus.Address = expandEnvVar(c.Milvus.Address)
    c.Milvus.Username = expandEnvVar(c.Milvus.Username)
    c.Milvus.Password = expandEnvVar(c.Milvus.Password)
    
    // 展开 Feishu 配置
    c.Feishu.WebhookURL = expandEnvVar(c.Feishu.WebhookURL)
}

// expandEnvVar 展开字符串中的环境变量引用
// 支持 ${VAR_NAME} 和 $VAR_NAME 两种语法
func expandEnvVar(s string) string {
    if s == "" {
        return s
    }

    // 匹配 ${VAR_NAME} 或 $VAR_NAME
    re := regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

    result := re.ReplaceAllStringFunc(s, func(match string) string {
        // 提取变量名
        varName := ""
        if strings.HasPrefix(match, "${") {
            // ${VAR_NAME} 格式
            varName = match[2 : len(match)-1]
        } else {
            // $VAR_NAME 格式
            varName = match[1:]
        }

        // 获取环境变量值
        value := os.Getenv(varName)
        if value != "" {
            return value
        }

        // 如果环境变量不存在，保持原样
        return match
    })

    return result
}
```

**语法**: `${VAR_NAME}` 或 `$VAR_NAME`

**示例**:
```yaml
# 方式1: ${} 语法（推荐，更清晰）
Embedding:
  APIKey: ${EMBEDDING_API_KEY}
  Model: ${EMBEDDING_MODEL}

# 方式2: $ 语法
openai:
  api_key: $OPENAI_API_KEY
```

**工作原理**:
1. 使用正则表达式匹配配置中的`${VAR}`或`$VAR`模式
2. 提取变量名（如`VAR_NAME`）
3. 调用`os.Getenv(varName)`获取环境变量的值
4. 如果环境变量存在，替换为实际值
5. 如果环境变量不存在，保持原样（显示`${VAR_NAME}`）

**注意事项**:
- 必须在配置加载后手动调用`cfg.ExpandEnv()`
- 如果环境变量不存在，不会报错，而是保持原样
- 敏感信息（API Key、密码等）应该使用环境变量
- 不支持默认值语法（如`${VAR:default}`），需要自己扩展

**使用流程**:
```go
// 1. 加载配置
cfg, _ := config.LoadConfig("config.yaml")

// 2. 展开环境变量
cfg.ExpandEnv()

// 3. 此时配置中的 ${EMBEDDING_API_KEY} 已被实际值替换
fmt.Println(cfg.Embedding.APIKey) // 输出实际的API Key
```

---

## 环境变量清单

### 必需的环境变量

| 环境变量 | 说明 | 示例值 | 是否必需 |
|---------|------|--------|---------|
| OPENAI_API_KEY | OpenAI API密钥 | sk-xxx | ✅ |
| DATABASE_DSN | 数据库连接串 | root:pass@tcp(localhost:3306)/db | ✅ |
| JWT_SECRET | JWT签名密钥 | your-secret-key | ✅ |
| REDIS_ADDR | Redis地址 | localhost:6379 | ✅ |
| REDIS_PASSWORD | Redis密码 | your-redis-password | ❌ |
| EMBEDDING_API_KEY | Embedding API密钥 | sk-xxx | ✅ |
| EMBEDDING_MODEL | Embedding模型ID | ep-xxx | ✅ |
| EMBEDDING_BASE_URL | Embedding服务地址 | https://ark.cn-beijing.volces.com/api/v3 | ✅ |
| EMBEDDING_REGION | Embedding服务区域 | cn-beijing | ✅ |
| MILVUS_ADDRESS | Milvus服务地址 | localhost:19530 | ❌ |
| MILVUS_USERNAME | Milvus用户名 | root | ❌ |
| MILVUS_PASSWORD | Milvus密码 | Milvus | ❌ |
| WECHAT_APP_ID | 微信AppID | wx1234567890abcdef | ❌ |
| WECHAT_APP_SECRET | 微信AppSecret | abc123def456 | ❌ |
| FEISHU_WEBHOOK_URL | 飞书机器人Webhook | https://open.feishu.cn/... | ❌ |

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

# Embedding服务配置（必需）
EMBEDDING_API_KEY=your-embedding-api-key
EMBEDDING_MODEL=your-model-endpoint-id
EMBEDDING_BASE_URL=https://ark.cn-beijing.volces.com/api/v3
EMBEDDING_REGION=cn-beijing

# Milvus配置（可选）
MILVUS_ADDRESS=localhost:19530
MILVUS_USERNAME=root
MILVUS_PASSWORD=Milvus

# 飞书告警配置（可选）
FEISHU_WEBHOOK_URL=
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

**配置选择方式**: 
```go
// 方式1: 通过命令行参数
./backend -config=config.prod.yaml

// 方式2: 通过环境变量
export CONFIG_FILE=config.prod.yaml
./backend

// 方式3: 在代码中根据环境变量选择
env := os.Getenv("ENV")
var configFile string
switch env {
case "prod":
    configFile = "config.prod.yaml"
case "test":
    configFile = "config.test.yaml"
default:
    configFile = "config.dev.yaml"
}
```

**Docker环境配置**:
在docker-compose.yml中可以为不同环境指定配置:
```yaml
services:
  backend:
    image: interview-backend
    volumes:
      - ./config.prod.yaml:/app/config.yaml
    environment:
      - ENV=production
```

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
**回答**: 
- **环境变量**: 敏感信息（API Key、密码）、不同环境差异大的配置
  - 优点：安全，不会被提交到代码仓库
  - 缺点：需要在每个环境单独配置
- **YAML直接配置**: 非敏感、通用的配置
  - 优点：清晰可见，方便版本管理
  - 缺点：敏感信息容易泄露

**最佳实践**:
- 密钥、密码 → 环境变量
- 超时时间、连接池大小 → YAML配置
- 服务地址 → 环境变量（不同环境不同）

### Q2: 如何动态更新配置而不重启服务？
**回答**: 
当前系统配置是启动时一次性加载，不支持热更新。如需支持热更新：
1. **文件监听**: 使用`fsnotify`库监听配置文件变化
2. **重新加载**: 文件变化时重新解析配置
3. **优雅更新**: 使用原子操作更新全局配置
4. **通知机制**: 通知各模块配置已更新

**示例**:
```go
import "github.com/fsnotify/fsnotify"

func WatchConfig(path string) {
    watcher, _ := fsnotify.NewWatcher()
    watcher.Add(path)
    
    for {
        select {
        case event := <-watcher.Events:
            if event.Op&fsnotify.Write == fsnotify.Write {
                cfg, _ := LoadConfig(path)
                cfg.ExpandEnv()
                // 原子更新
                atomic.StorePointer(&Global, unsafe.Pointer(cfg))
            }
        }
    }
}
```

### Q3: Eino配置和OpenAI配置有什么区别？
**回答**:
- **OpenAI配置**: 基础的OpenAI API连接配置（api_key, model_name, base_url）
- **Eino配置**: Eino框架的AI配置，包含更多参数：
  - `temperature`: 控制输出随机性（0-1）
  - `max_tokens`: 最大输出token数
  - `retry_count`: 失败重试次数
  - `retry_delay`: 重试间隔时间
  
Eino是对OpenAI等大模型的封装，提供了更多的功能和配置选项。

### Q4: 如何验证配置是否正确？
**回答**:
可以添加配置验证函数：
```go
func (c *Config) Validate() error {
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    if c.Database.DSN == "" {
        return errors.New("database dsn is required")
    }
    if c.OpenAI.APIKey == "" || strings.HasPrefix(c.OpenAI.APIKey, "${") {
        return errors.New("openai api_key is not set")
    }
    // ... 更多验证
    return nil
}

// 在加载配置后调用
cfg, _ := LoadConfig("config.yaml")
cfg.ExpandEnv()
if err := cfg.Validate(); err != nil {
    log.Fatalf("配置验证失败: %v", err)
}
```

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
