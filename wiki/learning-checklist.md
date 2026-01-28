# 学习清单与进度追踪

> 按照此清单逐步学习，完成一项打一个勾 ✅

## 📅 开始时间
**开始日期**: 2026-01-24  
**预计完成**: _____ （根据你的学习速度填写）

---

## 第一阶段：架构理解 (预计2天)

### Day 1: 整体认知
- [ ] 阅读 [README.md](../README.md)，理解项目概述
- [ ] 阅读 [doc/后端架构设计.md](../doc/后端架构设计.md)
- [ ] 阅读 [doc/技术实现方案.md](../doc/技术实现方案.md)
- [ ] 完成 Wiki: [01-architecture/README.md](01-architecture/README.md)
  - [ ] 记录系统的核心功能
  - [ ] 画出系统架构图
  - [ ] 理解技术栈选型原因

### Day 2: 启动流程与配置
- [ ] 阅读 [backend/main.go](../backend/main.go) (全部代码)
- [ ] 阅读 [backend/internal/config/config.go](../backend/internal/config/config.go)
- [ ] 理解配置文件 [backend/config.yaml](../backend/config.yaml)
- [ ] 完成 Wiki: [01-architecture/startup-flow.md](01-architecture/startup-flow.md)
  - [ ] 绘制启动流程图
  - [ ] 记录初始化步骤
  - [ ] 理解各组件依赖关系
- [ ] 完成 Wiki: [01-architecture/config-management.md](01-architecture/config-management.md)
  - [ ] 记录配置项说明
  - [ ] 理解环境变量机制

**阶段总结**: 
- 完成日期: _____
- 学习心得: _____

---

## 第二阶段：核心模块 (预计3-5天)

### Day 3-4: 数据层 (Repository Layer)

#### 数据模型
- [ ] 阅读 `backend/internal/model/` 目录下所有文件
  - [ ] user.go - 用户模型
  - [ ] interview_record.go - 面试记录
  - [ ] interview_evaluation.go - 面试评估
  - [ ] resume.go - 简历模型
  - [ ] prediction.go - 预测模型
- [ ] 查看 [backend/db_schema.sql](../backend/db_schema.sql/) 数据库设计
- [ ] 完成 Wiki: [02-data-layer/models.md](02-data-layer/models.md)
  - [ ] 画出ER图
  - [ ] 记录每个模型的字段和作用
  - [ ] 理解模型之间的关系

#### Repository层
- [ ] 阅读 `backend/internal/repository/` 目录
  - [ ] user/ - 用户仓储
  - [ ] interviews/ - 面试仓储
  - [ ] prediction/ - 预测仓储
- [ ] 完成 Wiki: [02-data-layer/repository-pattern.md](02-data-layer/repository-pattern.md)
  - [ ] 理解Repository模式
  - [ ] 记录CRUD操作实现
  - [ ] 理解缓存策略

### Day 5-6: 业务逻辑层 (Service Layer)
- [ ] 阅读 `backend/internal/service/` 目录
  - [ ] user/ - 用户服务
  - [ ] interviews/ - 面试服务
  - [ ] prediction/ - 预测服务
- [ ] 完成 Wiki: [03-service-layer/README.md](03-service-layer/README.md)
  - [ ] 记录服务层的职责
  - [ ] 理解业务逻辑处理
- [ ] 完成 Wiki: [03-service-layer/business-flow.md](03-service-layer/business-flow.md)
  - [ ] 画出关键业务流程图
  - [ ] 记录面试创建流程
  - [ ] 记录简历分析流程

### Day 7: API层
- [ ] 阅读 `backend/api/` 目录
  - [ ] handler/interview/ - 处理器
  - [ ] router/ - 路由配置
  - [ ] response/response.go - 响应封装
- [ ] 完成 Wiki: [04-api-layer/handler-pattern.md](04-api-layer/handler-pattern.md)
  - [ ] 理解Handler模式
  - [ ] 记录请求处理流程
- [ ] 完成 Wiki: [04-api-layer/api-documentation.md](04-api-layer/api-documentation.md)
  - [ ] 整理API接口列表
  - [ ] 记录请求/响应格式

**阶段总结**: 
- 完成日期: _____
- 学习心得: _____

---

## 第三阶段：AI核心 (预计3-4天)

### Day 8-9: 智能体架构

#### Eino框架理解
- [ ] 阅读 Eino 官方文档
- [ ] 理解 Agent、Tool、Chain 等概念
- [ ] 完成 Wiki: [05-ai-core/eino-framework.md](05-ai-core/eino-framework.md)
  - [ ] 记录Eino核心概念
  - [ ] 理解组件化设计
  - [ ] 理解流式处理机制

#### 智能体实现
- [ ] 阅读 `backend/chatApp/agent/` 目录
  - [ ] interview/ - 面试智能体
  - [ ] resume/ - 简历分析智能体
  - [ ] prediction/ - 预测智能体
  - [ ] record_evaluation/ - 评估智能体
- [ ] 完成 Wiki: [05-ai-core/agent-architecture.md](05-ai-core/agent-architecture.md)
  - [ ] 画出智能体架构图
  - [ ] 理解各智能体的职责
  - [ ] 记录智能体交互流程

### Day 10: 具体智能体深入
- [ ] 详细分析面试智能体实现
- [ ] 完成 Wiki: [05-ai-core/interview-agent.md](05-ai-core/interview-agent.md)
  - [ ] 记录问题生成逻辑
  - [ ] 理解对话管理机制
  - [ ] 记录追问策略
- [ ] 完成 Wiki: [05-ai-core/evaluation-agent.md](05-ai-core/evaluation-agent.md)
  - [ ] 记录评估标准
  - [ ] 理解打分机制
  - [ ] 记录反馈生成逻辑

### Day 11: AI工具与向量数据库
- [ ] 阅读 `backend/chatApp/tool/` 目录
- [ ] 阅读 `backend/internal/eino/milvus/` 目录
- [ ] 阅读 `backend/chatApp/chat/openAi.go`
- [ ] 完成 Wiki: [05-ai-core/vector-database.md](05-ai-core/vector-database.md)
  - [ ] 理解向量检索原理
  - [ ] 记录Milvus使用方式
  - [ ] 理解Embedding生成

**阶段总结**: 
- 完成日期: _____
- 学习心得: _____

---

## 第四阶段：实战准备 (预计2-3天)

### Day 12: 完整流程追踪
- [ ] 追踪一个完整的用户请求流程
  - [ ] 用户登录流程
  - [ ] 创建面试流程
  - [ ] 面试对话流程
  - [ ] 获取评估报告流程
- [ ] 画出完整的数据流图
- [ ] 完成 Wiki: [diagrams/data-flow.md](diagrams/data-flow.md)

### Day 13: 运行与调试
- [ ] 配置开发环境
- [ ] 运行系统
- [ ] 使用Postman测试API
- [ ] 打断点调试代码
- [ ] 记录遇到的问题和解决方案

### Day 14: 开发指南编写
- [ ] 完成 Wiki: [06-development-guide/add-new-feature.md](06-development-guide/add-new-feature.md)
  - [ ] 记录新功能开发流程
  - [ ] 提供代码模板
- [ ] 完成 Wiki: [06-development-guide/coding-standards.md](06-development-guide/coding-standards.md)
  - [ ] 总结代码规范
  - [ ] 记录最佳实践

**阶段总结**: 
- 完成日期: _____
- 学习心得: _____

---

## 🎯 学习成果验收

完成以上所有学习后，你应该能够：

- [ ] 清楚地画出系统架构图
- [ ] 说明任意一个API的完整处理流程
- [ ] 解释AI智能体的工作原理
- [ ] 独立添加一个新的API接口
- [ ] 独立实现一个简单的智能体功能
- [ ] 理解数据库表设计和关联关系
- [ ] 能够调试和解决常见问题

---

## 📝 学习笔记

### 重要发现
> 记录学习过程中的重要发现和Aha时刻

1. _____
2. _____
3. _____

### 遇到的问题
> 记录遇到的问题和解决方案

| 问题描述 | 解决方案 | 相关文档 |
|---------|---------|---------|
| _____ | _____ | _____ |
| _____ | _____ | _____ |

### 改进建议
> 对系统的改进建议

1. _____
2. _____
3. _____

---

## 📊 学习进度统计

- **第一阶段**: ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
- **第二阶段**: ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
- **第三阶段**: ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
- **第四阶段**: ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%
- **总体进度**: ⬜⬜⬜⬜⬜⬜⬜⬜⬜⬜ 0%

---

**下一步**: 开始 [第一阶段：架构理解](01-architecture/README.md)
