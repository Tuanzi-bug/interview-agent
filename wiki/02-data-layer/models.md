# 数据模型详解

> 理解系统中的所有数据模型及其关系

## 📋 学习任务

1. 阅读 `backend/internal/model/` 目录下的所有模型文件
2. 理解每个模型的字段和作用
3. 画出ER图（实体关系图）
4. 理解模型之间的关联关系

---

## 核心数据模型

### 1. User (用户模型)

**文件位置**: `backend/internal/model/user.go`

```go
// 粘贴User结构体代码
type User struct {
    // ...
}
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 用户ID | 主键 |
| Username | string | 用户名 | 唯一 |
| Email | string | 邮箱 | _____ |
| Password | string | 密码(加密) | _____ |
| _____ | _____ | _____ | _____ |

**关联关系**:
- 一个用户有多个简历 (1:N)
- 一个用户有多个面试记录 (1:N)

**我的理解**: _____

---

### 2. Resume (简历模型)

**文件位置**: `backend/internal/model/resume.go`

```go
// 粘贴Resume结构体代码
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 简历ID | 主键 |
| UserID | uint | 用户ID | 外键 |
| Content | string | 简历内容 | _____ |
| _____ | _____ | _____ | _____ |

**关联关系**:
- 属于一个用户 (N:1)
- 关联多个预测结果 (1:N)

**JSON标签**: _____

**我的理解**: _____

---

### 3. InterviewRecord (面试记录模型)

**文件位置**: `backend/internal/model/interview_record.go`

```go
// 粘贴InterviewRecord结构体代码
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 面试ID | 主键 |
| UserID | uint | 用户ID | 外键 |
| Type | string | 面试类型 | _____ |
| Status | string | 面试状态 | _____ |
| _____ | _____ | _____ | _____ |

**面试类型枚举**:
- `comprehensive`: 综合面试
- `technical`: 专项面试
- _____

**面试状态枚举**:
- `pending`: 待开始
- `in_progress`: 进行中
- `completed`: 已完成
- _____

**关联关系**:
- 属于一个用户 (N:1)
- 有多条对话记录 (1:N)
- 有一份评估报告 (1:1)

**我的理解**: _____

---

### 4. InterviewDialogue (面试对话模型)

**文件位置**: `backend/internal/model/interview_dialogue.go`

```go
// 粘贴InterviewDialogue结构体代码
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 对话ID | 主键 |
| InterviewID | uint | 面试ID | 外键 |
| Role | string | 角色 | _____ |
| Content | string | 对话内容 | _____ |
| _____ | _____ | _____ | _____ |

**角色类型**:
- `user`: 用户
- `assistant`: AI助手
- _____

**我的理解**: _____

---

### 5. InterviewEvaluation (面试评估模型)

**文件位置**: `backend/internal/model/interview_evaluation.go`

```go
// 粘贴InterviewEvaluation结构体代码
```

**字段说明**:

| 字段 | 类型 | 说明 | 评分范围 |
|------|------|------|---------|
| ID | uint | 评估ID | - |
| InterviewID | uint | 面试ID | - |
| TotalScore | float64 | 总分 | 0-100 |
| _____ | _____ | _____ | _____ |

**评估维度**:
- 技术能力: _____
- 表达能力: _____
- 问题理解: _____
- _____

**我的理解**: _____

---

### 6. Prediction (预测模型)

**文件位置**: `backend/internal/model/prediction.go`

```go
// 粘贴Prediction结构体代码
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 预测ID | 主键 |
| ResumeID | uint | 简历ID | 外键 |
| Category | string | 预测类别 | _____ |
| _____ | _____ | _____ | _____ |

**预测类别**:
- _____
- _____

**我的理解**: _____

---

## 实体关系图 (ER图)

```
请画出实体关系图

示例：
┌─────────┐        ┌─────────┐
│  User   │1──────n│ Resume  │
└─────────┘        └─────────┘
     │1
     │
     │n
┌─────────────────┐
│InterviewRecord  │
└─────────────────┘
     │1
     ├──n──┐
     │     │
     │1    │1
┌─────────────┐  ┌──────────────────┐
│Dialogue     │  │Evaluation        │
└─────────────┘  └──────────────────┘
```

---

## GORM特性使用

### 1. gorm.Model

```go
type Model struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

**作用**: _____

**我的理解**: _____

---

### 2. 外键关联

```go
type Resume struct {
    UserID uint
    User   User `gorm:"foreignKey:UserID"`
}
```

**语法**: _____

**级联操作**: _____

**我的理解**: _____

---

### 3. JSON标签

```go
type User struct {
    Email string `json:"email" gorm:"uniqueIndex"`
}
```

**json标签作用**: _____

**gorm标签作用**: _____

**我的理解**: _____

---

### 4. 软删除

```go
DeletedAt gorm.DeletedAt `gorm:"index"`
```

**什么是软删除**: _____

**优点**: _____

**我的理解**: _____

---

## 数据验证

### 验证规则

记录各模型的验证规则：

```go
// 示例
type User struct {
    Email string `binding:"required,email"`
    Password string `binding:"required,min=6"`
}
```

**常用验证标签**:
- `required`: _____
- `email`: _____
- `min`: _____
- `max`: _____

---

## 索引设计

### 需要索引的字段

记录哪些字段需要索引及原因：

| 表名 | 字段 | 索引类型 | 原因 |
|------|------|---------|------|
| users | email | 唯一索引 | _____ |
| users | username | 唯一索引 | _____ |
| interview_records | user_id | 普通索引 | _____ |
| _____ | _____ | _____ | _____ |

**我的理解**: _____

---

## 模型方法

### 自定义方法

记录模型上定义的自定义方法：

```go
// 示例
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 密码加密等操作
}
```

**钩子函数**:
- `BeforeCreate`: _____
- `AfterCreate`: _____
- `BeforeUpdate`: _____
- _____

**我的理解**: _____

---

## 数据库迁移

### 自动迁移

```go
db.AutoMigrate(&User{}, &Resume{}, &InterviewRecord{})
```

**AutoMigrate做了什么**: _____

**限制**: _____

**生产环境建议**: _____

**我的理解**: _____

---

## ❓ 问题记录

### Q1: 为什么使用软删除而不是硬删除？
**我的理解**: _____

### Q2: 外键约束在应用层还是数据库层？
**我的理解**: _____

### Q3: _____
**我的理解**: _____

---

## 🔬 实践验证

### 任务1: 查看数据库表结构
```sql
-- 连接MySQL
DESCRIBE users;
DESCRIBE resumes;
-- 记录实际的表结构
```

**我的发现**: _____

### 任务2: 测试模型关联
```go
// 编写测试代码
user := User{ID: 1}
db.Preload("Resumes").Find(&user)
// 验证关联加载
```

**我的发现**: _____

---

## ✅ 学习检查点

- [ ] 理解所有核心数据模型
- [ ] 画出完整的ER图
- [ ] 理解GORM的核心特性
- [ ] 掌握模型关联关系
- [ ] 理解软删除机制
- [ ] 掌握数据验证方法
- [ ] 查看过实际的数据库表结构

---

**下一步**: [Repository模式](repository-pattern.md)
