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
├── go.mod               # Go 模块文件
├── internal/
│   ├── config/          # 配置管理
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   ├── middleware/      # 中间件
│   ├── service/         # 业务逻辑层
│   └── handler/         # API 处理器
└── pkg/
    ├── eino/            # Eino 框架集成
    └── hertz/           # Hertz 框架配置
```

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
   git clone https://github.com/your-username/ai-eino-interview-agent.git
   cd ai-eino-interview-agent
   ```

2. 配置环境
   编辑 `config.yaml` 文件，填写相关配置：
   - 数据库连接信息
   - Redis 连接信息
   - API Key 等

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