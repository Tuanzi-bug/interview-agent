# 第二阶段：核心模块

> 目标：深入理解数据层、业务逻辑层和API层的实现

## 📚 本阶段学习内容

### 数据层 (Repository Pattern)
1. [数据模型详解](models.md)
2. [Repository模式](repository-pattern.md)
3. [数据库设计](database-design.md)

### 业务逻辑层
参见 [第三阶段目录](../03-service-layer/README.md)

### API层
参见 [第四阶段目录](../04-api-layer/README.md)

---

## 学习路径

```
Day 3-4: 数据层
  ├── 理解数据模型 (models.md)
  ├── 学习Repository模式 (repository-pattern.md)
  └── 分析数据库设计 (database-design.md)

Day 5-6: 业务逻辑层
  ├── 理解Service层职责
  ├── 分析业务流程
  └── 追踪完整业务逻辑

Day 7: API层
  ├── 理解Handler模式
  ├── 学习路由设计
  └── 掌握API规范
```

---

## 核心概念

### 分层架构

```
┌─────────────────────┐
│   API Layer         │ ← 处理HTTP请求和响应
│   (Handler)         │
├─────────────────────┤
│   Service Layer     │ ← 业务逻辑处理
│   (Business Logic)  │
├─────────────────────┤
│   Repository Layer  │ ← 数据访问抽象
│   (Data Access)     │
├─────────────────────┤
│   Model Layer       │ ← 数据模型定义
│   (Data Structure)  │
└─────────────────────┘
```

### 各层职责

**API层 (Handler)**:
- 接收HTTP请求
- 参数验证
- 调用Service层
- 返回HTTP响应

**Service层**:
- 业务逻辑处理
- 数据校验
- 事务控制
- 调用Repository层

**Repository层**:
- 数据库操作
- 缓存管理
- 数据访问抽象

**Model层**:
- 数据结构定义
- 数据库表映射

---

## 数据流分析

### 典型请求流程

```
HTTP Request
  ↓
Handler (API层)
  ├─ 参数解析
  ├─ 参数验证
  └─ 调用Service
      ↓
Service (业务逻辑层)
  ├─ 业务逻辑处理
  ├─ 数据校验
  └─ 调用Repository
      ↓
Repository (数据访问层)
  ├─ 数据库查询
  ├─ 缓存处理
  └─ 返回数据
      ↓
Service
  ├─ 数据处理
  └─ 返回结果
      ↓
Handler
  ├─ 封装响应
  └─ 返回HTTP Response
```

---

## 学习重点

### 1. 理解Repository模式
- [ ] 什么是Repository模式
- [ ] 为什么使用Repository模式
- [ ] 如何实现Repository接口

### 2. 掌握数据模型
- [ ] GORM模型定义
- [ ] 模型关联关系
- [ ] 数据验证

### 3. 理解业务逻辑
- [ ] Service层的职责
- [ ] 业务流程设计
- [ ] 错误处理

### 4. 掌握API设计
- [ ] RESTful API规范
- [ ] 请求/响应格式
- [ ] 错误响应设计

---

## 实践任务

### 任务1: 追踪用户登录流程
- [ ] 从Handler开始
- [ ] 追踪到Service
- [ ] 追踪到Repository
- [ ] 画出完整流程图

### 任务2: 追踪面试创建流程
- [ ] 理解创建面试的完整过程
- [ ] 记录每一层的处理逻辑
- [ ] 分析数据持久化过程

### 任务3: 分析数据模型关系
- [ ] 画出ER图
- [ ] 理解外键关联
- [ ] 理解级联操作

---

## ✅ 阶段目标

完成本阶段后，你应该能够：

- [ ] 清楚地说明数据在各层之间如何流转
- [ ] 理解Repository模式的优势
- [ ] 独立实现一个简单的CRUD功能
- [ ] 理解数据库表设计
- [ ] 掌握Service层的业务逻辑处理
- [ ] 能够设计RESTful API接口

---

**开始学习**: [数据模型详解](models.md)
