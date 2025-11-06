# AI-Eino-Interview-Agent

基于字节跳动 Hertz 框架和 Eino 大语言模型框架开发的智能面试代理系统。

## 项目概述

该项目是一个智能面试系统，利用 Eino 框架提供的大语言模型能力，实现简历分析、面试问题生成、答案评估等功能。系统采用 Hertz 作为 Web 框架，提供 RESTful API 接口。

## 技术栈

- **后端框架**: [Hertz](https://github.com/cloudwego/hertz) - 字节跳动高性能 Web 框架
- **AI 框架**: [Eino](https://github.com/cloudwego/eino) - 大语言模型应用框架
- **数据库**: MySQL + GORM
- **缓存**: Redis
- **认证**: JWT

## 功能特性

- 用户注册和登录
- 简历管理（上传、更新、删除）
- 简历智能分析
- 面试创建和管理
- AI 驱动的面试问题生成
- 面试答案评估
- 面试结果报告

## 项目结构

```
├── main.go              # 应用程序入口
├── config.yaml          # 配置文件
├── config.example.yaml  # 配置示例文件
├── go.mod               # Go 模块文件
├── go.sum               # Go 依赖校验文件
├── db_schema.sql        # 数据库模式定义
├── internal/            # 内部包，不对外暴露
│   ├── config/          # 配置管理
│   │   └── config.go    # 配置结构和加载逻辑
│   ├── model/           # 数据模型定义
│   │   └── models.go    # 实体模型定义
│   ├── repository/      # 数据访问层
│   │   ├── database.go  # 数据库连接和操作
│   │   └── redis.go     # Redis连接和缓存操作
│   ├── middleware/      # 中间件
│   │   └── jwt.go       # JWT认证中间件
│   ├── service/         # 业务逻辑层
│   │   ├── user_service.go      # 用户相关业务逻辑
│   │   ├── resume_service.go    # 简历相关业务逻辑
│   │   └── interview_service.go # 面试相关业务逻辑
│   ├── handler/         # API 处理器
│   │   ├── user_handler.go      # 用户API处理
│   │   ├── resume_handler.go    # 简历API处理
│   │   └── interview_handler.go # 面试API处理
│   └── utils/           # 工具函数
├── pkg/                 # 可重用的公共包
│   ├── eino/            # Eino 框架集成
│   │   └── eino_manager.go # Eino管理器
│   └── hertz/           # Hertz 框架配置
│       ├── hertz_manager.go # Hertz管理器
│       └── swagger_config.go # Swagger配置
├── docs/                # API文档
│   ├── docs.go          # Swagger文档生成
│   ├── swagger.json     # Swagger JSON文档
│   └── swagger.yaml     # Swagger YAML文档
├── doc/                 # 项目文档
│   ├── swagger_implementation_plan.md # Swagger实现计划
│   ├── 后端架构设计.md    # 后端架构设计文档
│   ├── 技术实现方案.md    # 技术实现方案文档
│   └── 需求文档.md        # 需求文档
├── chatApp/             # 聊天应用相关代码
│   ├── chat/            # 聊天功能
│   ├── config/          # 聊天应用配置
│   └── tool/            # 聊天应用工具
├── cmd/                 # 命令行入口
│   └── api/             # API服务命令
└── frontend/            # 前端代码目录
```

### 目录功能说明

**根目录**: 包含项目的主要配置文件、入口文件和依赖管理文件
- `main.go`: 应用程序的主入口，负责初始化和启动服务
- `config.yaml`: 项目配置文件，包含数据库、Redis、API等配置
- `config.example.yaml`: 配置示例文件，提供配置模板
- `go.mod`: Go模块定义文件，管理项目依赖
- `go.sum`: 依赖版本锁定文件，确保依赖一致性
- `db_schema.sql`: 数据库表结构定义脚本

**internal/**: 内部包目录，根据Go语言约定，这些包不会被外部项目导入
- `config/`: 配置管理，负责加载和解析配置文件
- `model/`: 数据模型定义，定义数据库表对应的Go结构体
- `repository/`: 数据访问层，封装数据库和缓存操作
- `middleware/`: 中间件，包含认证、日志等横切关注点
- `service/`: 业务逻辑层，实现核心业务功能
- `handler/`: API处理器，处理HTTP请求和响应
- `utils/`: 工具函数，提供通用辅助功能

**pkg/**: 可重用的公共包，可以被其他项目导入使用
- `eino/`: Eino框架集成，封装AI大模型相关功能
- `hertz/`: Hertz框架配置，提供Web框架相关的配置和工具

**docs/**: API文档目录
- 包含自动生成的Swagger文档，用于API接口说明和测试

**doc/**: 项目文档目录
- 包含架构设计、技术方案、需求文档等项目说明文件

**chatApp/**: 聊天应用相关代码
- 实现与AI模型的对话功能，用于面试问答交互

**cmd/**: 命令行入口目录
- 包含可执行程序的入口点

**frontend/**: 前端代码目录
- 存放前端应用代码，负责用户界面实现

## 配置说明

配置文件 `config.yaml` 包含以下主要配置项：

- 服务配置（端口、主机）
- 数据库连接信息
- Redis 连接信息
- OpenAI API 配置（用于 Eino 框架）
- JWT 密钥
- CORS 配置

## 安装和运行

### 前提条件

- Go 1.20+
- MySQL 数据库
- Redis 服务
- 有效的 OpenAI API Key（或其他支持的模型 API Key）

### 安装步骤

1. 克隆项目
   ```bash
   git clone git@codeup.aliyun.com:60fadd729187b7df39056384/training_camp/go-eino-interview-agent.git
   cd ai-eino-interview-agent
   ```

2. 配置环境
   编辑 `config.yaml` 文件，填写相关配置：
   - 数据库连接信息
   - Redis 连接信息
   - API Key 等
   - 把 db_schema.sql 中的数据库 schema 导入到 MySQL 数据库中

3. 安装依赖
   ```bash
   go mod download
   ```

4. 运行服务
   ```bash
   go run main.go
   ```

## API 接口

### 用户相关
- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录
- `GET /api/v1/user/profile` - 获取用户信息
- `PUT /api/v1/user/profile` - 更新用户信息

### 简历相关
- `POST /api/v1/resume` - 创建简历
- `GET /api/v1/resume` - 获取用户所有简历
- `GET /api/v1/resume/:id` - 获取指定简历
- `PUT /api/v1/resume/:id` - 更新简历
- `DELETE /api/v1/resume/:id` - 删除简历
- `POST /api/v1/resume/:id/analyze` - 分析简历

### 面试相关
- `POST /api/v1/interview` - 创建面试
- `GET /api/v1/interview` - 获取用户所有面试
- `GET /api/v1/interview/:id` - 获取指定面试
- `POST /api/v1/interview/:id/start` - 开始面试
- `POST /api/v1/interview/:id/end` - 结束面试
- `POST /api/v1/interview/:id/answer` - 提交答案
- `GET /api/v1/interview/:id/questions` - 获取面试问题
- `GET /api/v1/interview/:id/result` - 获取面试结果

## 健康检查

- `GET /health` - 服务健康检查

## 开发说明

### 中间件

- JWT 认证中间件保护需要授权的接口
- CORS 中间件处理跨域请求

### 错误处理

系统采用统一的错误响应格式：
```json
{
  "error": "错误信息"
}
```

### 日志记录

系统使用标准库日志记录关键操作和错误信息。

## 注意事项

- 确保配置文件中的 API Key 和数据库连接信息安全存储
- 生产环境建议使用环境变量或密钥管理服务
- 定期更新依赖包以获取安全补丁

## 许可证

MIT