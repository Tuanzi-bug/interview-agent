# 技术知识库

系统架构、框架集成、技术决策和最佳实践

**最后更新**: 2026-02-01

## 领域概览

| 领域 | 描述 | 关键内容 |
|------|------|---------|
| [架构设计](architecture.md) | 系统架构和设计决策 | 微服务拆分、数据流、模块化设计 |
| [Eino 框架](eino-framework.md) | Eino AI 框架集成 | Agent 设计、Chain 编排、流式处理 |
| [Hertz 框架](hertz-framework.md) | Hertz Web 框架集成 | 路由、中间件、请求处理 |
| [测试策略](testing-strategy.md) | Go 测试和 E2E 测试 | 单元测试、集成测试、Mock 策略 |
| [构建系统](build-system.md) | 构建配置和优化 | Docker、Make、CI/CD |

## 技术栈总览

### 核心框架
- **Hertz** (v0.10.3) - 字节跳动高性能 HTTP 框架
  - 特点：高易用性、高性能、高扩展性
  - 用途：API 服务、路由管理、中间件
  
- **Eino** (v0.5.11+) - 大语言模型应用框架
  - 特点：组件化设计、图编排引擎、流式处理
  - 用途：AI Agent、Chain 编排、LLM 应用

### 编程语言和工具
- **Go**: 1.24+ (高性能、并发、编译型语言)
- **MySQL**: 结构化数据存储
- **Redis**: 缓存和会话管理
- **Milvus**: 向量数据库（文档检索）

### AI 模型集成
- **OpenAI**: ChatGPT API 集成
- **ARK** (字节 Ark): 内部 LLM 模型服务
- **Google Search**: 搜索工具集成

### 关键依赖库

| 库 | 版本 | 用途 |
|----|----|------|
| gorm | v1.25.11 | ORM - 数据库操作 |
| milvus-sdk-go | v2.4.2 | 向量检索 |
| redis/go-redis | v9.16.0 | Redis 客户端 |
| jwt | v5.2.2 | JWT 认证 |
| thrift | v0.21.0 | IDL 接口定义 |
| yaml | v3.0.1 | 配置管理 |
| sonic | v1.14.2 | JSON 处理 |

## 系统分层

```
┌─────────────────────────────────────────┐
│         前端 (Next.js)                   │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    API 层 (Hertz Router & Handler)      │
│  - 路由定义                              │
│  - 请求验证                              │
│  - 响应封装                              │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    服务层 (Business Logic)               │
│  - 用户服务                              │
│  - 简历服务                              │
│  - 面试服务                              │
│  - AI 评估服务                           │
│  - 向量检索服务                          │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    AI 层 (Eino Agent & Chain)            │
│  - 简历分析 Agent                        │
│  - 问题生成 Agent                        │
│  - 评估 Agent                            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    数据层 (Repository)                   │
│  - MySQL (GORM)                          │
│  - Redis (Cache)                         │
│  - Milvus (Vector DB)                    │
└─────────────────────────────────────────┘
```

## 核心服务模块

### 1. 用户服务 (User Service)
- **职责**: 用户认证、授权、信息管理
- **技术**: JWT, GORM, Redis
- **关键端点**: /auth/login, /user/profile, /user/update

### 2. 简历服务 (Resume Service)
- **职责**: 简历上传、解析、AI 分析
- **技术**: Eino Agent, PDF Parser, Milvus
- **关键端点**: /resume/upload, /resume/analyze, /resume/qa

### 3. 面试服务 (Interview Service)
- **职责**: 面试创建、问题生成、记录管理
- **技术**: Hertz, GORM, Eino
- **关键端点**: /interview/start, /interview/question, /interview/submit

### 4. AI 评估服务 (Evaluation Service)
- **职责**: 答案评估、反馈生成、报告
- **技术**: Eino Chain, LLM, GORM
- **关键端点**: /evaluate/answer, /evaluate/report

## 数据模型

### 主要实体
- **User**: 用户账户信息
- **Resume**: 简历记录
- **Interview**: 面试记录
- **Question**: 面试问题
- **Answer**: 用户答案
- **Evaluation**: 评估结果

### 数据存储策略

| 数据类型 | 存储选择 | 原因 |
|---------|---------|------|
| 用户、面试记录 | MySQL | 结构化、需要 ACID |
| 会话、临时数据 | Redis | 高性能、临时存储 |
| 简历向量、文档 | Milvus | 向量检索、相似度匹配 |
| 文件 | 对象存储 | 文件管理、备份 |

## 配置管理

### 环境配置
- 位置: `backend/config.yaml`
- 内容: 数据库、Redis、LLM 密钥、Milvus 地址
- 覆盖: 环境变量可覆盖配置文件

### 构建配置
- Docker: 容器化部署
- Docker Compose: 本地开发和测试
- Makefile: 构建和部署脚本

## 关键决策

### 为什么选择 Hertz?
✅ 高性能 HTTP 框架
✅ 字节跳动生产级验证
✅ 与 Eino 生态兼容
✅ Go 标准库风格 API

### 为什么选择 Eino?
✅ 专为大模型应用设计
✅ 组件化、易于扩展
✅ 图编排引擎功能强大
✅ 流式处理支持

### 为什么选择 GORM + MySQL?
✅ Go 标准 ORM
✅ 关系型数据适合业务数据
✅ 成熟的生态和支持

## 常见技术问题

**Q: 如何集成新的 LLM 模型?**
→ 查看 [Eino 框架文档](eino-framework.md) 中的模型集成章节

**Q: Hertz 路由如何配置?**
→ 查看 [Hertz 框架文档](hertz-framework.md) 中的路由管理

**Q: 如何实现 AI Agent?**
→ 查看 [Eino 框架文档](eino-framework.md) 中的 Agent 设计

**Q: 如何进行单元测试?**
→ 查看 [测试策略文档](testing-strategy.md)

## 相关资源

### 官方文档
- [Hertz 官方文档](https://www.hertzframework.com/)
- [Eino 官方文档](https://github.com/cloudwego/eino)
- [GORM 文档](https://gorm.io/)

### 项目文档
- [架构设计详述](architecture.md)
- [系统架构 Wiki](../../wiki/01-architecture/)
- [后端架构设计](../../doc/后端架构设计.md)

### 最佳实践
- 查看 [测试策略](testing-strategy.md) 了解测试方法
- 查看 [构建系统](build-system.md) 了解部署流程

---

**维护者**: GitHub Copilot  
**最后更新**: 2026-02-01  
**版本**: 1.0
