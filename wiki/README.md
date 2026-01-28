# 面试吧系统学习Wiki

> 这是一个系统化学习和记录「面试吧 AI智能面试平台」的知识库

## 📖 Wiki目录结构

```
wiki/
├── README.md                          # 本文件 - Wiki主索引
├── learning-checklist.md              # 学习清单和进度追踪
├── 01-architecture/                   # 第一阶段：架构理解
│   ├── README.md                      # 架构总览
│   ├── startup-flow.md                # 启动流程分析
│   ├── config-management.md           # 配置管理
│   └── tech-stack.md                  # 技术栈详解
├── 02-data-layer/                     # 第二阶段：数据层
│   ├── README.md                      # 数据层总览
│   ├── models.md                      # 数据模型
│   ├── repository-pattern.md          # Repository模式
│   └── database-design.md             # 数据库设计
├── 03-service-layer/                  # 第三阶段：业务逻辑层
│   ├── README.md                      # 服务层总览
│   ├── user-service.md                # 用户服务
│   ├── interview-service.md           # 面试服务
│   └── business-flow.md               # 业务流程图
├── 04-api-layer/                      # 第四阶段：API层
│   ├── README.md                      # API层总览
│   ├── handler-pattern.md             # Handler模式
│   ├── router-design.md               # 路由设计
│   └── api-documentation.md           # API文档
├── 05-ai-core/                        # 第五阶段：AI核心
│   ├── README.md                      # AI核心总览
│   ├── agent-architecture.md          # 智能体架构
│   ├── eino-framework.md              # Eino框架理解
│   ├── interview-agent.md             # 面试智能体详解
│   ├── resume-agent.md                # 简历分析智能体
│   ├── evaluation-agent.md            # 评估智能体
│   └── vector-database.md             # 向量数据库
├── 06-development-guide/              # 实战开发指南
│   ├── README.md                      # 开发指南总览
│   ├── add-new-feature.md             # 如何添加新功能
│   ├── coding-standards.md            # 编码规范
│   ├── testing-guide.md               # 测试指南
│   └── deployment-guide.md            # 部署指南
└── diagrams/                          # 架构图和流程图
    ├── system-architecture.md         # 系统架构图
    ├── data-flow.md                   # 数据流图
    └── ai-workflow.md                 # AI工作流图
```

## 🎯 学习路径

### 快速入门（3天速成）
- **Day 1**: [架构理解](01-architecture/) + [启动流程](01-architecture/startup-flow.md)
- **Day 2**: [数据层](02-data-layer/) + [业务逻辑层](03-service-layer/)
- **Day 3**: [API层](04-api-layer/) + [AI核心](05-ai-core/)

### 深度学习（1-2周）
- **Week 1**: 
  - 第一阶段：[架构理解](01-architecture/)（2天）
  - 第二阶段：[核心模块](02-data-layer/)（3天）
- **Week 2**: 
  - 第三阶段：[AI核心](05-ai-core/)（4天）
  - 第四阶段：[实战开发](06-development-guide/)（3天）

## 📝 使用方法

1. **按顺序学习**: 从第一阶段开始，逐步深入
2. **记录笔记**: 在每个md文件中填写你的学习笔记和理解
3. **画图辅助**: 在`diagrams/`目录中画出你的理解
4. **更新清单**: 完成一个模块后，更新[learning-checklist.md](learning-checklist.md)
5. **实践验证**: 通过实际代码运行验证你的理解

## 🔖 快速导航

### 核心概念
- [系统架构](01-architecture/README.md)
- [数据模型](02-data-layer/models.md)
- [AI智能体](05-ai-core/agent-architecture.md)

### 开发实战
- [添加新功能指南](06-development-guide/add-new-feature.md)
- [API开发规范](04-api-layer/handler-pattern.md)
- [测试指南](06-development-guide/testing-guide.md)

### 问题记录
在学习过程中遇到的问题，可以记录在对应模块的文件中，使用以下格式：

```markdown
## ❓ 问题记录

### Q1: 问题描述
**解决方案**: ...
**相关代码**: [文件路径](../backend/xxx.go)
```

## 📊 学习进度追踪

查看 [learning-checklist.md](learning-checklist.md) 追踪你的学习进度。

## 💡 贡献指南

这是你的个人学习Wiki，可以：
- ✅ 添加自己的理解和注释
- ✅ 补充代码示例和图表
- ✅ 记录遇到的问题和解决方案
- ✅ 添加学习心得和总结

## 🚀 开始学习

从这里开始：[学习清单](learning-checklist.md) → [第一阶段：架构理解](01-architecture/README.md)
