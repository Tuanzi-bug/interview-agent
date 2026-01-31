# 面试吧 AI 智能面试平台 - 项目知识库

**项目**: Go-Eino Interview Agent  
**初始化**: 2026-02-01  
**最后更新**: 2026-02-01

## 项目概述

这是面试吧平台的核心知识库，捕获系统架构、技术模式、常见问题和最佳实践。知识库帮助团队成员快速上手并有效协作。

### 项目简介

- **名称**: 面试吧 - AI 智能面试平台
- **目标**: 利用 AI 技术帮助求职者轻松应对各类面试挑战
- **核心技术**: Hertz(字节 Web 框架) + Eino(大模型应用框架)
- **关键功能**: 简历分析、AI 面试、答案评估、推荐预测

## 知识库结构

### 技术知识 (`tech/`)
系统架构、框架集成、构建配置和技术模式

- [技术栈概览](tech/README.md) - 框架、工具和技术选型
- [架构设计](tech/architecture.md) - 系统架构决策和设计模式
- [Eino AI 框架](tech/eino-framework.md) - Eino 框架集成和最佳实践
- [Hertz 框架](tech/hertz-framework.md) - Hertz Web 框架集成
- [测试策略](tech/testing-strategy.md) - Go 测试模式和 E2E 测试
- [构建系统](tech/build-system.md) - 构建配置和优化

### 经验知识 (`experience/`)
常见错误、调试技巧、工作流优化和实战经验

- [错误模式](experience/error-patterns.md) - 常见错误和解决方案
- [测试模式](experience/test-patterns.md) - 有效的测试模式
- [调试技巧](experience/debugging-tactics.md) - 调试策略和工具
- [工作流优化](experience/workflow-improvements.md) - 流程优化和最佳实践

## 快速导航

### 获取信息
1. 检查此 README 了解主题概览
2. 导航到相关领域文件获取详细信息
3. 使用搜索查找特定模式或问题

### 添加新知识
使用 `/learn` 命令捕获学习：
```
/learn [你的学习/模式/解决方案]
```

### 查找已有知识
- **搜索技术问题**: 查看 `tech/` 目录
- **搜索常见错误**: 查看 `experience/error-patterns.md`
- **查找测试模式**: 查看 `tech/testing-strategy.md` 和 `experience/test-patterns.md`
- **了解架构**: 查看 `tech/architecture.md`

## 项目统计

| 指标 | 数值 |
|-----|------|
| **技术域文件** | 6 |
| **经验域文件** | 4 |
| **总域文件** | 10 |
| **最后学习添加** | 2026-02-01 |

## 技术栈速览

### 后端框架
- **Hertz**: 字节跳动高性能 HTTP 框架
- **Eino**: 大语言模型应用框架
- **Go**: 1.24+ 版本

### 数据存储
- **MySQL**: 结构化数据（用户、面试记录等）
- **Redis**: 缓存和会话管理
- **Milvus**: 向量检索和文档管理

### 关键依赖
- **gorm**: ORM 框架
- **jwt**: 认证
- **thrift**: IDL 定义
- **openai/ark**: LLM 集成

## 核心服务

| 服务 | 职责 | 关键技术 |
|------|------|---------|
| 用户服务 | 注册、登录、认证 | JWT, GORM |
| 简历服务 | 上传、解析、分析 | Eino Agent |
| 面试服务 | 问题生成、记录管理 | Hertz, GORM |
| AI 评估服务 | 答案评估、反馈 | Eino Chain, LLM |
| 向量检索服务 | 文档索引和检索 | Milvus, Eino |

## 使用示例

### 查找特定主题

**如何集成 Eino 框架?**
→ 查看 [Eino 框架文档](tech/eino-framework.md)

**遇到 Hertz 路由问题?**
→ 查看 [错误模式](experience/error-patterns.md) 中的 Hertz 相关错误

**如何编写有效测试?**
→ 查看 [测试策略](tech/testing-strategy.md) 和 [测试模式](experience/test-patterns.md)

**系统如何设计?**
→ 查看 [架构设计](tech/architecture.md)

## 维护指南

### 何时更新
- 发现新的技术模式
- 实现重大功能变更后
- 用户要求使用 **update memory bank** 命令
- 需要澄清上下文时

### 更新流程
1. 查看相关的域文件
2. 使用 `/learn` 命令添加新知识
3. 由系统自动更新相应文件
4. 保持文档的一致性和可访问性

### 文档质量标准
- ✅ 清晰的主题标题
- ✅ 可扫描的组织方式（列表、表格、代码块）
- ✅ 工作示例和实际代码
- ✅ 跨文档的双向链接
- ✅ 明确的应用场景说明

## 相关资源

### 项目文档
- [项目 README](../../README.md) - 项目概述和入门指南
- [后端架构设计](../../doc/后端架构设计.md) - 详细的架构文档
- [技术实现方案](../../doc/技术实现方案.md) - 实现细节
- [项目优化建议](../../doc/项目优化建议_索引.md) - 优化建议汇总

### Wiki 文档
- [Wiki 首页](../../wiki/README.md) - Wiki 知识库入口
- [系统架构](../../wiki/01-architecture/) - 架构详细说明
- [数据层设计](../../wiki/02-data-layer/) - 数据库设计
- [服务层](../../wiki/03-service-layer/) - 业务逻辑

### 需求文档
- [需求索引](../requirements/INDEX.md) - 所有需求的中央索引
- [项目需求文档](../../doc/需求文档.md) - 功能需求详情

---

## 快速开始

### 新开发者入门
1. 阅读本 README 了解整体结构
2. 查看 [技术栈概览](tech/README.md) 了解技术选择
3. 查看 [架构设计](tech/architecture.md) 理解系统设计
4. 查看 [Hertz 框架](tech/hertz-framework.md) 和 [Eino 框架](tech/eino-framework.md)
5. 查看相关的测试和调试资源

### 遇到问题
1. 检查 [错误模式](experience/error-patterns.md) 查找类似问题
2. 查看 [调试技巧](experience/debugging-tactics.md) 了解排查方法
3. 使用搜索功能查找特定主题

### 添加新学习
完成任务后，使用 `/learn` 记录发现的模式、解决的问题或工作流改进。

---

**维护者**: GitHub Copilot  
**版本**: 1.0  
**创建日期**: 2026-02-01  
**最后更新**: 2026-02-01  

📚 **[返回顶部](#面试吧-ai-智能面试平台---项目知识库)**
