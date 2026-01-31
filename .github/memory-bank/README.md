# 🧠 Memory Bank - Go-Eino Interview Agent

项目持久化知识库，包含架构决策、技术模式、常见问题和需求追踪

**项目**: 面试吧 - AI 智能面试平台  
**版本**: 1.0  
**最后更新**: 2026-02-01

---

## 📚 知识库结构

### 核心知识库 (`context/`)

结构化知识库，包含技术文档和实战经验

```
context/
├── README.md                    # 知识库入口
├── tech/                        # 技术知识
│   ├── README.md               # 技术索引
│   ├── architecture.md         # 系统架构设计
│   ├── eino-framework.md       # Eino AI 框架
│   ├── hertz-framework.md      # Hertz Web 框架
│   ├── testing-strategy.md     # 测试策略
│   └── build-system.md         # 构建系统
└── experience/                  # 经验知识
    ├── README.md               # 经验索引
    ├── error-patterns.md       # 常见错误和解决方案
    ├── test-patterns.md        # 有效的测试模式
    ├── debugging-tactics.md    # 调试技巧
    └── workflow-improvements.md # 工作流优化
```

### 需求管理 (`requirements/`)

项目需求追踪和管理

```
requirements/
├── INDEX.md                     # 需求索引（中央追踪）
├── in-progress/                 # 进行中的需求
└── completed/                   # 已完成的需求
```

### 任务追踪 (`tasks/`)

系统维护和改进任务

```
tasks/
└── _index.md                    # 任务索引（1 个完成的任务）
```

---

## 🎯 快速导航

### 按用途

| 用途 | 位置 | 说明 |
|------|------|------|
| **项目入门** | [context/README.md](../../context/README.md) | 新人快速上手 |
| **理解架构** | [context/tech/architecture.md](../../context/tech/architecture.md) | 系统设计和模块 |
| **框架学习** | [context/tech/](../../context/tech/) | Eino、Hertz 使用指南 |
| **问题排查** | [context/experience/error-patterns.md](../../context/experience/error-patterns.md) | 常见错误解决 |
| **改进工作流** | [context/experience/workflow-improvements.md](../../context/experience/workflow-improvements.md) | 开发效率优化 |
| **需求追踪** | [requirements/INDEX.md](../../requirements/INDEX.md) | 项目管理和规划 |

### 按角色

**👨‍💻 新开发者**
1. 阅读 [context/README.md](../../context/README.md)
2. 查看 [context/tech/README.md](../../context/tech/README.md)
3. 学习相关框架文档
4. 遇到问题查看 error-patterns.md

**🏗️ 架构师**
1. 查看 [context/tech/architecture.md](../../context/tech/architecture.md)
2. 理解系统设计和技术决策
3. 参考 [context/tech/](../../context/tech/) 中的最佳实践

**🧪 QA/测试工程师**
1. 查看 [context/tech/testing-strategy.md](../../context/tech/testing-strategy.md)
2. 学习 [context/experience/test-patterns.md](../../context/experience/test-patterns.md)
3. 参考测试用例和断言方法

**📊 项目经理**
1. 查看 [requirements/INDEX.md](../../requirements/INDEX.md)
2. 追踪项目进度和需求
3. 使用 `/plan` 创建新需求

---

## 📖 文档概览

### 技术知识库 (6 个文档)

#### 1. 架构设计 (`architecture.md`)
**包含**: 系统架构图、微服务拆分、数据流、技术决策、扩展设计
**关键内容**: 5 个核心服务、3 层数据存储、部署架构

#### 2. Eino 框架 (`eino-framework.md`)
**包含**: 框架基础、Agent 开发、Chain 组合、流式处理、模型集成
**关键内容**: 3 个 Agent 设计（简历分析、问题生成、答案评估）、RAG 模式、错误处理

#### 3. Hertz 框架 (`hertz-framework.md`)
**包含**: 路由管理、中间件、请求/响应处理、文件上传、WebSocket
**关键内容**: 路由分组、认证中间件、统一响应格式、最佳实践

#### 4. 测试策略 (`testing-strategy.md`)
**包含**: 单元测试、集成测试、E2E 测试、覆盖率、CI/CD
**关键内容**: 测试框架、Mock 工具、性能测试

#### 5. 构建系统 (`build-system.md`)
**包含**: Docker、Docker Compose、Makefile、部署流程
**关键内容**: 前后端 Dockerfile、compose 配置、常用命令

### 经验知识库 (4 个文档)

#### 1. 错误模式 (`error-patterns.md`)
**包含**: 9+ 常见错误，每个都有原因和解决方案
**错误类型**: Go 基础、Hertz 框架、数据库、Redis、AI 集成、并发问题

#### 2. 测试模式 (`test-patterns.md`)
**包含**: 表驱动测试、Mock 策略、并发测试、反模式
**关键技巧**: 接口 Mock、HTTP Mock、断言方法、WaitGroup 同步

#### 3. 调试技巧 (`debugging-tactics.md`)
**包含**: 日志分析、性能分析、竞态检测、调试工具
**工具**: Delve 调试器、pprof、race 检测

#### 4. 工作流优化 (`workflow-improvements.md`)
**包含**: Git 工作流、代码审查、发布流程、性能优化
**最佳实践**: Commit 规范、审查清单、缓存策略

---

## 🚀 使用方式

### 添加新学习

```bash
/learn [描述你发现的模式、错误或最佳实践]

# 示例
/learn We discovered that goroutines leak when not properly canceled with context
```

系统自动:
- 检测学习的类型和领域
- 在相应文件中添加内容
- 更新索引和统计

### 创建新需求

```bash
/plan [描述你想要实现的功能]

# 示例
/plan Implement user feedback and rating system for interview sessions
```

系统自动:
- 创建需求文档
- 更新需求索引
- 提供实现计划模板

### 查看统计

每个索引文件都包含最新统计:
- 文档计数
- 最后更新时间
- 内容概览

---

## 📊 项目统计

### 知识库规模

| 指标 | 数值 |
|-----|------|
| **技术文档** | 6 |
| **经验文档** | 4 |
| **需求文件** | 1 (INDEX) |
| **任务文件** | 1 (INDEX) |
| **总计** | 12 |

### 内容深度

| 类型 | 数量 | 示例 |
|------|------|------|
| **错误模式** | 9+ | goroutine 泄漏、路由 404、DB 超时等 |
| **测试模式** | 4+ | 表驱动测试、Mock 策略、并发测试 |
| **技术决策** | 4+ | 框架选择、架构设计、数据存储 |
| **最佳实践** | 10+ | 日志、缓存、性能、工作流 |

---

## ✨ 关键特性

### 🎯 索引优先策略
- 快速导航入口
- 多层级目录结构
- 清晰的分类和标签

### 📚 完整的技术栈文档
- Hertz Web 框架
- Eino AI 框架
- Go 编程最佳实践
- 构建和部署

### 🔍 实战经验总结
- 常见错误及解决方案
- 有效的测试模式
- 调试工具和技巧
- 工作流优化建议

### 📊 项目追踪
- 中央需求索引
- 任务管理
- 进度统计

### 🔗 交叉引用
- 文档之间相互链接
- 快速跳转
- 完整的上下文

---

## 🎓 学习路径建议

### 初级 (1-2 天)
1. 阅读 [context/README.md](../../context/README.md) - 了解整体结构
2. 查看 [context/tech/architecture.md](../../context/tech/architecture.md) - 理解系统设计
3. 选择一个框架学习 (Eino 或 Hertz)

### 中级 (1 周)
1. 完整学习 [context/tech/](../../context/tech/) 所有文档
2. 阅读 [context/experience/error-patterns.md](../../context/experience/error-patterns.md)
3. 参考 [context/experience/test-patterns.md](../../context/experience/test-patterns.md)

### 高级 (持续)
1. 贡献新的学习和最佳实践
2. 参与架构决策
3. 优化工作流和流程

---

## 🛠️ 维护指南

### 定期更新

- **每周**: 检查新的错误或问题
- **每月**: 更新统计信息和索引
- **每季度**: 审查和优化文档结构

### 保持质量

- ✅ 示例代码必须可运行
- ✅ 链接必须有效
- ✅ 信息必须准确
- ✅ 格式保持一致

### 版本控制

- 使用 git 追踪所有更改
- 有意义的 commit 消息
- 定期备份

---

## 📞 获取帮助

### 问题排查

1. **查看 error-patterns.md** - 大多数问题都有解决方案
2. **查看相关框架文档** - 特定框架的配置和用法
3. **使用搜索功能** - 全局文本搜索

### 贡献知识

1. 发现新的错误模式？ → 使用 `/learn` 添加
2. 想到了优化？ → 提议到 workflow-improvements.md
3. 发现了 bug？ → 创建 issue 并链接到需求

---

## 相关文档

### 内部链接
- [知识库主页](../../context/README.md)
- [需求索引](../../requirements/INDEX.md)
- [初始化报告](../../KNOWLEDGE_BASE_INITIALIZED.md)

### 项目文档
- [项目 README](../../README.md)
- [后端架构设计](../../doc/后端架构设计.md)
- [Wiki](../../wiki/README.md)

---

**维护者**: GitHub Copilot  
**版本**: 1.0  
**创建日期**: 2026-02-01  
**最后更新**: 2026-02-01  

📚 **[返回知识库](../../context/README.md)**
