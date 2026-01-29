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
type User struct {
    ID            uint           `json:"id" gorm:"primaryKey"`
    Username      string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
    Email         string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    PasswordHash  string         `json:"-" gorm:"size:255;not null"`
    Role          string         `json:"role" gorm:"size:20;default:'user'"`
    WechatOpenID  *string        `json:"wechat_open_id" gorm:"uniqueIndex;size:100"`
    WechatUnionID *string        `json:"wechat_union_id" gorm:"uniqueIndex;size:100"`
    Nickname      string         `json:"nickname" gorm:"size:100"`
    Avatar        string         `json:"avatar" gorm:"size:255"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint | 用户ID | 主键，自增 |
| Username | string | 用户名 | 唯一索引，最大50字符，非空 |
| Email | string | 邮箱 | 唯一索引，最大100字符，非空 |
| PasswordHash | string | 密码哈希 | 最大255字符，非空，JSON序列化时隐藏 |
| Role | string | 用户角色 | 最大20字符，默认'user' |
| WechatOpenID | *string | 微信OpenID | 唯一索引，支持微信登录，可为空 |
| WechatUnionID | *string | 微信UnionID | 唯一索引，支持微信多平台登录，可为空 |
| Nickname | string | 昵称 | 最大100字符，来自微信或用户设置 |
| Avatar | string | 头像URL | 最大255字符，用户头像地址 |
| CreatedAt | time.Time | 创建时间 | 自动管理 |
| UpdatedAt | time.Time | 更新时间 | 自动管理 |
| DeletedAt | gorm.DeletedAt | 删除时间 | 软删除标记，有索引 |

**关联关系**:
- 一个用户有多个简历 (1:N) - 通过 Resume.UserID 关联
- 一个用户有多个面试记录 (1:N) - 通过 InterviewRecord.UserID 关联
- 一个用户有多个用户模型配置 (1:N) - 通过 UserModel.UserID 关联
- 一个用户有多个押题记录 (1:N) - 通过 PredictionRecord.UserID 关联

**我的理解**: User模型是系统的核心实体，采用了邮箱+密码和微信登录双认证方式。密码使用哈希存储确保安全性，PasswordHash字段通过`json:"-"`标签避免在API响应中暴露。支持软删除机制，删除的用户数据不会立即从数据库移除。微信登录字段使用指针类型(*string)表示可选，这样未绑定微信的用户该字段为NULL而非空字符串。

---

### 2. Resume (简历模型)

**文件位置**: `backend/internal/model/resume.go`

```go
type Resume struct {
    ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement;comment:简历ID"`
    UserID    uint      `json:"user_id" gorm:"index;not null;comment:用户ID"`
    Content   string    `json:"content" gorm:"type:longtext;not null;comment:简历内容"`
    FileName  string    `json:"file_name" gorm:"size:255;comment:原始文件名"`
    FileSize  int64     `json:"file_size" gorm:"comment:文件大小（字节）"`
    FileType  string    `json:"file_type" gorm:"size:50;comment:文件类型"`
    IsDefault int       `json:"is_default" gorm:"default:0;comment:是否为默认简历"`
    Deleted   int       `json:"deleted" gorm:"default:0;index;comment:删除标记"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint64 | 简历ID | 主键，自增 |
| UserID | uint | 用户ID | 外键，有索引，非空 |
| Content | string | 简历内容 | longtext类型，存储解析后的文本，非空 |
| FileName | string | 原始文件名 | 最大255字符，保留上传时的文件名 |
| FileSize | int64 | 文件大小 | 字节数，用于存储限制检查 |
| FileType | string | 文件类型 | 最大50字符，如pdf/doc/txt |
| IsDefault | int | 是否默认简历 | 0-否，1-是；用户可设置一份默认简历 |
| Deleted | int | 删除标记 | 0-未删除，1-已删除；软删除标记，有索引 |
| CreatedAt | time.Time | 创建时间 | 毫秒时间戳，自动创建 |
| UpdatedAt | time.Time | 更新时间 | 毫秒时间戳，自动更新 |

**关联关系**:
- 属于一个用户 (N:1) - 通过 UserID 外键关联
- 关联多个预测记录 (1:N) - 通过 PredictionRecord.ResumeID 关联

**JSON标签**: 所有字段都有JSON标签，用于API序列化时的字段名映射，采用snake_case命名风格

**我的理解**: Resume模型存储用户上传的简历信息。Content字段使用longtext类型支持大文本存储（最大4GB），足够存储任何简历内容。采用整数类型的软删除标记而非DeletedAt，这样可以避免GORM的自动软删除行为，给应用层更多控制权。IsDefault字段允许用户设置默认简历，在面试和押题时可以自动使用。查询时需要添加`deleted = 0`条件过滤已删除记录。

---

### 3. InterviewRecord (面试记录模型)

**文件位置**: `backend/internal/model/interview_record.go`

```go
type InterviewRecord struct {
    ID           uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID       uint      `json:"user_id" gorm:"index;not null"`
    Type         string    `json:"type" gorm:"size:255;not null"`
    Difficulty   string    `json:"difficulty" gorm:"size:128;not null"`
    Domain       string    `json:"domain" gorm:"size:255;not null"`
    CompanyName  string    `json:"company_name" gorm:"size:128"`
    PositionName string    `json:"position_name" gorm:"size:128"`
    Status       string    `json:"status" gorm:"size:50;not null;default:'pending'"`
    Duration     int64     `json:"duration" gorm:"comment:面试耗时（秒）"`
    CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime:milli"`
}
```

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint64 | 面试ID | 主键，自增 |
| UserID | uint | 用户ID | 外键，有索引，非空 |
| Type | string | 面试类型 | 最大255字符，comprehensive/technical |
| Difficulty | string | 难度级别 | 最大128字符，简单/中等/困难 |
| Domain | string | 面试领域 | 最大255字符，如"校招""社招""Java""Go" |
| CompanyName | string | 公司名称 | 最大128字符，可选 |
| PositionName | string | 岗位名称 | 最大128字符，可选 |
| Status | string | 面试状态 | 最大50字符，默认pending |
| Duration | int64 | 面试耗时 | 秒数，记录实际面试时长 |
| CreatedAt | time.Time | 创建时间 | 毫秒时间戳 |
| UpdatedAt | time.Time | 更新时间 | 毫秒时间戳 |

**面试类型枚举**:
- `comprehensive`: 综合面试 - 涵盖多个技术领域的全面面试
- `technical`: 专项面试 - 针对特定技术栈的深度面试
- 类型存储在Type字段，可扩展支持其他面试类型

**面试状态枚举**:
- `pending`: 待开始 - 面试已创建但未开始
- `in_progress`: 进行中 - 面试正在进行
- `completed`: 已完成 - 面试已结束并生成评估
- `cancelled`: 已取消 - 用户中途取消（可选状态）

**关联关系**:
- 属于一个用户 (N:1) - 通过 UserID 外键关联
- 有多条对话记录 (1:N) - 通过 InterviewDialogue.ReportID 关联
- 有一份评估报告 (1:1) - 通过 InterviewEvaluation.ReportID 关联
- 有一份答题报告 (1:1) - 通过 AnswerReport.ReportID 关联

**我的理解**: InterviewRecord是面试的主记录表，记录面试的元信息。Domain字段非常灵活，可以存储"校招/社招"这样的分类，也可以存储"Java/Go"这样的技术栈。Duration字段用于统计面试时长，可用于数据分析和用户反馈。Status字段跟踪面试生命周期，前端可以根据状态显示不同的UI（如继续面试、查看报告等）。使用uint64作为ID类型而非uint，为超大规模应用预留空间。

---

### 4. InterviewDialogue (面试对话模型)

**文件位置**: `backend/internal/model/interview_dialogue.go`

```go
type InterviewDialogue struct {
    ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_id"`
    ReportID  uint64    `json:"report_id" gorm:"not null;index:idx_report_id"`
    Question  string    `json:"question" gorm:"type:text"`
    Answer    string    `json:"answer" gorm:"type:text"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli"`
}
```

**表名**: `interview_dialogues`（注意是复数形式）

**字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint64 | 对话ID | 主键，自增 |
| UserID | uint | 用户ID | 非空，有索引idx_user_id |
| ReportID | uint64 | 面试记录ID | 非空，有索引idx_report_id，外键关联 |
| Question | string | 提问内容 | text类型，智能体的提问内容 |
| Answer | string | 回答内容 | text类型，用户的回答内容 |
| CreatedAt | time.Time | 创建时间 | 毫秒时间戳，记录对话发生时间 |

**数据模型设计特点**:
- 采用**一问一答**的存储方式，一条记录包含问题和答案
- Question字段存储智能体的提问内容
- Answer字段存储用户的回答内容
- 相比传统的Role+Content方式，这种设计更符合面试场景的语义

**角色隐含关系**:
- Question → 智能体（assistant角色）
- Answer → 候选人（user角色）
- 无需额外的Role字段，通过字段名即可区分

**索引设计**:
- `idx_user_id`: 快速查询某用户的所有对话
- `idx_report_id`: 快速查询某次面试的所有对话记录
- 查询时通常按 `ORDER BY id ASC` 保持对话顺序

**我的理解**: InterviewDialogue采用了简化的对话存储模型。与常见的聊天系统不同（通常用Role字段区分user/assistant），这里直接用Question/Answer两个字段表示一轮对话，使数据语义更清晰。这种设计特别适合面试场景，因为面试是结构化的一问一答，不像自由聊天那样角色可能随意切换。ReportID关联到InterviewRecord表，实现面试记录与对话的一对多关系。查询对话时按ID升序排列即可还原完整的面试过程。**注意**：实际表名是`interview_dialogues`（复数形式）。

---

### 5. InterviewEvaluation (面试评估模型)

**文件位置**: `backend/internal/model/interview_evaluation.go`

```go
type InterviewEvaluation struct {
    ID         uint64                 `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID     uint                   `json:"user_id" gorm:"not null;index:idx_user_id"`
    ReportID   uint64                 `json:"report_id" gorm:"not null;index:idx_report_id"`
    Comment    string                 `json:"comment" gorm:"type:text"`
    Score      float64                `json:"score" gorm:"type:decimal(5,2)"`
    Dimensions []*EvaluationDimension `json:"dimensions" gorm:"type:json;serializer:json"`
    Deleted    int                    `json:"deleted" gorm:"default:0;index:idx_deleted"`
    CreatedAt  time.Time              `json:"created_at" gorm:"autoCreateTime:milli"`
    UpdatedAt  time.Time              `json:"updated_at" gorm:"autoUpdateTime:milli"`
}

type EvaluationDimension struct {
    DimensionName string  `json:"dimension_name"`
    Evaluation    string  `json:"evaluation"`
    Score         float64 `json:"score"`
}
```

**字段说明**:

| 字段 | 类型 | 说明 | 评分范围 |
|------|------|------|---------|
| ID | uint64 | 评估ID | - |
| UserID | uint | 用户ID | - |
| ReportID | uint64 | 面试记录ID | 外键，关联interview_record |
| Comment | string | 总体评价 | text类型，AI生成的总结性评语 |
| Score | float64 | 总体评分 | decimal(5,2)，范围通常0-100 |
| Dimensions | []*EvaluationDimension | 各维度评估 | JSON数组，多维度评分 |
| Deleted | int | 删除标记 | 0-未删除，1-已删除 |
| CreatedAt | time.Time | 创建时间 | 毫秒时间戳 |
| UpdatedAt | time.Time | 更新时间 | 毫秒时间戳 |

**EvaluationDimension（评估维度）字段**:

| 字段 | 类型 | 说明 |
|------|------|------|
| DimensionName | string | 维度名称（如"技术能力""沟通表达"） |
| Evaluation | string | 该维度的详细评价文本 |
| Score | float64 | 该维度的评分 |

**评估维度**（存储在Dimensions JSON数组中）:
- **技术能力**: 评估候选人对技术知识的掌握深度和广度，算法思维，代码质量
- **表达能力**: 评估候选人的语言组织能力、逻辑清晰度、沟通效率
- **问题理解**: 评估候选人对问题的理解速度和准确性，需求分析能力
- **实战经验**: 评估候选人的项目经历、实际问题解决能力
- **学习潜力**: 评估候选人的学习态度、知识面的广度、技术热情
- **综合素质**: 评估候选人的团队协作、抗压能力等软技能

**JSON存储优势**:
- 灵活性：可以动态调整评估维度，无需修改表结构
- 可扩展：不同面试类型可以有不同的维度配置
- 原子性：评估和维度作为一个整体读写，无需事务处理多表
- 查询方便：MySQL 5.7+支持JSON字段的查询和索引

**我的理解**: InterviewEvaluation采用了单表+JSON的设计模式，将多维度评估存储在Dimensions JSON字段中，这比传统的一对多关系表更灵活。总体评分Score是各维度评分的加权平均或综合评定。Comment字段存储AI生成的总结性评语，通常包括优势总结、待改进点、面试建议等。使用软删除机制（Deleted字段）而非GORM的DeletedAt，与Resume表保持一致的删除策略。ReportID作为外键关联到InterviewRecord，实现一次面试对应一份评估报告的1:1关系。

---

### 6. PredictionRecord & PredictionQuestion (押题记录和题目模型)

**文件位置**: `backend/internal/model/prediction.go`

```go
// 押题记录主表
type PredictionRecord struct {
    ID         uint64               `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID     uint                 `json:"user_id" gorm:"index;not null"`
    ResumeID   uint64               `json:"resume_id" gorm:"index;not null"`
    Type       string               `json:"type" gorm:"size:20"`
    Language   string               `json:"language" gorm:"size:20"`
    JobTitle   string               `json:"job_title" gorm:"size:50"`
    Difficulty string               `json:"difficulty" gorm:"size:20"`
    Company    string               `json:"company" gorm:"size:100"`
    Questions  []PredictionQuestion `json:"questions" gorm:"foreignKey:RecordID"`
    CreatedAt  time.Time            `json:"created_at" gorm:"autoCreateTime"`
}

// 押题题目明细表
type PredictionQuestion struct {
    ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
    RecordID        uint64    `json:"record_id" gorm:"index;not null"`
    Question        string    `json:"question" gorm:"type:text;not null"`
    Content         string    `json:"content" gorm:"type:text"`
    Focus           string    `json:"focus" gorm:"type:text"`
    ThinkingPath    string    `json:"thinking_path" gorm:"type:text"`
    ReferenceAnswer string    `json:"reference_answer" gorm:"type:text"`
    FollowUp        string    `json:"follow_up" gorm:"type:text"`
    Sort            int       `json:"sort"`
    CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
}
```

**PredictionRecord 字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint64 | 押题记录ID | 主键，自增 |
| UserID | uint | 用户ID | 有索引，非空 |
| ResumeID | uint64 | 简历ID | 有索引，非空，外键关联 |
| Type | string | 押题类型 | 最大20字符，校招/社招 |
| Language | string | 编程语言 | 最大20字符，java/go/python等 |
| JobTitle | string | 岗位名称 | 最大50字符，前端/后端/全栈等 |
| Difficulty | string | 难度等级 | 最大20字符，入门/进阶/高级 |
| Company | string | 目标公司 | 最大100字符，字节/阿里/腾讯等 |
| Questions | []PredictionQuestion | 题目列表 | 一对多关联，外键RecordID |
| CreatedAt | time.Time | 创建时间 | 自动创建 |

**PredictionQuestion 字段说明**:

| 字段 | 类型 | 说明 | 约束 |
|------|------|------|------|
| ID | uint64 | 题目ID | 主键，自增 |
| RecordID | uint64 | 押题记录ID | 有索引，非空，外键 |
| Question | string | 问题内容 | text类型，非空，题目描述 |
| Content | string | 考察内容 | text类型，这道题主要考什么 |
| Focus | string | 重点考察 | text类型，考察的核心知识点 |
| ThinkingPath | string | 回答思路 | text类型，建议的答题思路 |
| ReferenceAnswer | string | 参考答案 | text类型，标准或推荐答案 |
| FollowUp | string | 可能追问 | text类型，面试官可能的追问点 |
| Sort | int | 题目排序 | 数字越小越靠前，用于控制题目顺序 |
| CreatedAt | time.Time | 创建时间 | 自动创建 |

**押题类型（Type字段）**:
- `校招`: 针对应届生的校园招聘面试
- `社招`: 针对有工作经验的社会招聘面试

**关联关系**:
- PredictionRecord 属于一个用户 (N:1) - 通过 UserID 关联
- PredictionRecord 基于一份简历 (N:1) - 通过 ResumeID 关联
- PredictionRecord 包含多个题目 (1:N) - 通过 Questions 关联
- PredictionQuestion 属于一条押题记录 (N:1) - 通过 RecordID 关联

**GORM外键配置**:
```go
Questions []PredictionQuestion `gorm:"foreignKey:RecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
```
- OnUpdate:CASCADE - 主表更新时级联更新
- OnDelete:CASCADE - 主表删除时级联删除子表记录

**我的理解**: Prediction采用了主从表设计，PredictionRecord存储押题的元信息（谁、什么时候、针对什么岗位），PredictionQuestion存储具体的题目内容。这种设计便于：1) 批量生成题目；2) 题目可以独立查询和展示；3) 方便统计（如某份简历生成了多少道题）。每个题目包含完整的答题指导（Focus、ThinkingPath、ReferenceAnswer），不仅给出题目，还教用户如何准备。FollowUp字段预测面试官的追问，帮助用户做更充分的准备。Sort字段控制题目顺序，AI可以按重要性或难度排序。级联删除确保删除押题记录时相关题目也被清理，避免孤儿数据。

---

## 实体关系图 (ER图)

```mermaid
erDiagram
    User ||--o{ Resume : "拥有"
    User ||--o{ InterviewRecord : "创建"
    User ||--o{ PredictionRecord : "生成"
    User ||--o{ UserModel : "配置"
    User ||--o{ InterviewDialogue : "参与"
    User ||--o{ InterviewEvaluation : "获得"
    User ||--o{ AnswerReport : "拥有"
    
    Resume ||--o{ PredictionRecord : "用于押题"
    
    InterviewRecord ||--o{ InterviewDialogue : "包含对话"
    InterviewRecord ||--o| InterviewEvaluation : "生成评估"
    InterviewRecord ||--o| AnswerReport : "生成报告"
    
    PredictionRecord ||--|{ PredictionQuestion : "包含题目"
    
    User {
        uint ID PK
        string Username UK
        string Email UK
        string PasswordHash
        string Role
        string WechatOpenID UK
        string WechatUnionID UK
        string Nickname
        string Avatar
        time CreatedAt
        time UpdatedAt
        time DeletedAt
    }
    
    Resume {
        uint64 ID PK
        uint UserID FK
        string Content
        string FileName
        int64 FileSize
        string FileType
        int IsDefault
        int Deleted
        time CreatedAt
        time UpdatedAt
    }
    
    InterviewRecord {
        uint64 ID PK
        uint UserID FK
        string Type
        string Difficulty
        string Domain
        string CompanyName
        string PositionName
        string Status
        int64 Duration
        time CreatedAt
        time UpdatedAt
    }
    
    InterviewDialogue {
        uint64 ID PK
        uint UserID FK
        uint64 ReportID FK
        string Question
        string Answer
        time CreatedAt
    }
    
    InterviewEvaluation {
        uint64 ID PK
        uint UserID FK
        uint64 ReportID FK
        string Comment
        float64 Score
        json Dimensions
        int Deleted
        time CreatedAt
        time UpdatedAt
    }
    
    AnswerReport {
        uint64 ID PK
        uint UserID FK
        uint64 ReportID FK
        json Records "答题记录数组"
        int Deleted
        time CreatedAt
        time UpdatedAt
    }
    
    PredictionRecord {
        uint64 ID PK
        uint UserID FK
        uint64 ResumeID FK
        string Type
        string Language
        string JobTitle
        string Difficulty
        string Company
        time CreatedAt
    }
    
    PredictionQuestion {
        uint64 ID PK
        uint64 RecordID FK
        string Question
        string Content
        string Focus
        string ThinkingPath
        string ReferenceAnswer
        string FollowUp
        int Sort
        time CreatedAt
    }
    
    UserModel {
        uint64 ID PK
        int64 UserID FK
        string Name "模型显示名称"
        string ModelKey "模型标识"
        string Protocol "协议类型"
        string BaseURL "API基础地址"
        string APIKeyEncrypted "加密API密钥"
        string ConfigJSON "额外配置"
        string SecretHint "密钥脱敏提示"
        string ProviderName "提供商名称"
        int64 MetaID "关联model_meta"
        string DefaultParams "默认参数JSON"
        int Scope "使用范围"
        int Status "状态"
        int IsDefault "是否默认"
        int Deleted "删除状态"
        int64 CreatedAt "毫秒时间戳"
        int64 UpdatedAt "毫秒时间戳"
    }
```

**关系说明**:

| 关系 | 类型 | 说明 |
|------|------|------|
| User → Resume | 1:N | 一个用户可以上传多份简历 |
| User → InterviewRecord | 1:N | 一个用户可以进行多次面试 |
| User → PredictionRecord | 1:N | 一个用户可以生成多次押题 |
| User → UserModel | 1:N | 一个用户可以配置多个AI模型 |
| Resume → PredictionRecord | 1:N | 一份简历可以用于多次押题 |
| InterviewRecord → InterviewDialogue | 1:N | 一次面试包含多轮对话 |
| InterviewRecord → InterviewEvaluation | 1:1 | 一次面试生成一份评估报告 |
| InterviewRecord → AnswerReport | 1:1 | 一次面试生成一份答题报告 |
| PredictionRecord → PredictionQuestion | 1:N | 一次押题包含多道题目 |

**主要外键关系**:
- Resume.UserID → User.ID
- InterviewRecord.UserID → User.ID
- InterviewDialogue.UserID → User.ID
- InterviewDialogue.ReportID → InterviewRecord.ID
- InterviewEvaluation.UserID → User.ID
- InterviewEvaluation.ReportID → InterviewRecord.ID
- AnswerReport.UserID → User.ID
- AnswerReport.ReportID → InterviewRecord.ID
- PredictionRecord.UserID → User.ID
- PredictionRecord.ResumeID → Resume.ID
- PredictionQuestion.RecordID → PredictionRecord.ID
- UserModel.UserID → User.ID

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

**作用**: 
- GORM提供的标准模型嵌入结构，包含4个常用字段
- `ID`: uint类型主键，自动递增
- `CreatedAt`: 创建时间，插入记录时自动设置
- `UpdatedAt`: 更新时间，更新记录时自动刷新
- `DeletedAt`: 软删除时间戳，删除时不真正删除而是设置此字段

**本项目的使用情况**:
- User模型使用了自定义的DeletedAt字段：`DeletedAt gorm.DeletedAt`
- 其他模型（Resume、InterviewRecord等）使用单独的字段定义，未直接嵌入gorm.Model
- 原因：更灵活的字段控制（如使用uint64作为ID，使用int作为删除标记等）

**我的理解**: gorm.Model是GORM提供的便捷嵌入结构，适合快速开发。但本项目为了更精细的控制（如ID类型、删除机制、时间戳精度），选择手动定义这些字段。例如InterviewRecord使用uint64的ID支持更大规模，Resume使用int类型的Deleted字段而非DeletedAt提供更灵活的软删除控制。这是根据业务需求做的架构权衡。

---

### 2. 外键关联

```go
// 定义关联
type PredictionRecord struct {
    ResumeID  uint64               `json:"resume_id" gorm:"index;not null"`
    Questions []PredictionQuestion `json:"questions" gorm:"foreignKey:RecordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type PredictionQuestion struct {
    RecordID uint64 `json:"record_id" gorm:"index;not null"`
}
```

**语法说明**:
- `foreignKey:RecordID` - 指定外键字段名称
- `references:ID` - 指定引用的主表字段（可省略，默认为主键）
- `constraint:OnUpdate:CASCADE` - 主表更新时级联更新
- `constraint:OnDelete:CASCADE` - 主表删除时级联删除
- `constraint:OnDelete:SET NULL` - 主表删除时设置外键为NULL

**级联操作示例**:
```go
// CASCADE: 删除押题记录时，自动删除所有关联的题目
db.Delete(&PredictionRecord{ID: 1})
// 自动执行: DELETE FROM prediction_question WHERE record_id = 1

// SET NULL: 删除简历时，将押题记录的ResumeID设为NULL
// Questions []PredictionQuestion `gorm:"foreignKey:ResumeID;constraint:OnDelete:SET NULL"`

// RESTRICT: 如果有关联数据，禁止删除主表记录
// Questions []PredictionQuestion `gorm:"foreignKey:RecordID;constraint:OnDelete:RESTRICT"`
```

**本项目的外键处理**:
1. **数据库层外键**:
   - PredictionQuestion → PredictionRecord: 定义了外键约束和级联删除
   
2. **应用层外键**:
   - 大多数关联在应用层处理，数据库层不设外键约束
   - 原因：提高灵活性，避免外键锁，便于分库分表
   - 例如：Resume.UserID、InterviewRecord.UserID 只是普通索引字段

**Preload关联查询**:
```go
// 查询押题记录及其所有题目
var record PredictionRecord
db.Preload("Questions", func(db *gorm.DB) *gorm.DB {
    return db.Order("sort ASC")
}).First(&record, id)

// 查询用户及其所有简历
var user User
db.Preload("Resumes").Find(&user, userID)
```

**我的理解**: GORM支持定义外键约束，但本项目大部分表采用应用层外键管理。只有PredictionQuestion明确使用了数据库外键（CASCADE级联删除），因为题目完全依赖押题记录存在。其他关联如User-Resume、User-InterviewRecord等都在应用层通过索引+查询实现，这样更灵活，避免了数据库外键带来的性能开销和迁移复杂度。使用Preload可以方便地加载关联数据，避免N+1查询问题。

---

### 3. JSON标签与GORM标签

```go
type User struct {
    ID           uint           `json:"id" gorm:"primaryKey"`
    Email        string         `json:"email" gorm:"uniqueIndex;size:100;not null"`
    PasswordHash string         `json:"-" gorm:"size:255;not null"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
```

**json标签作用**:
- 控制结构体序列化为JSON时的字段名
- `json:"email"` - 字段序列化为"email"
- `json:"-"` - 序列化时忽略该字段（如密码、DeletedAt）
- `json:"email,omitempty"` - 字段为零值时不序列化

**gorm标签作用**:
- 定义数据库表结构和约束
- `primaryKey` - 主键
- `autoIncrement` - 自增
- `uniqueIndex` - 唯一索引
- `index` - 普通索引
- `size:100` - 字段长度
- `not null` - 非空约束
- `default:'value'` - 默认值
- `type:text` - 字段类型
- `comment:说明` - 字段注释
- `foreignKey:UserID` - 外键定义
- `serializer:json` - JSON序列化（用于存储复杂对象）

**常用GORM标签组合**:
```go
// 主键ID
ID uint64 `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`

// 唯一字段
Email string `json:"email" gorm:"uniqueIndex;size:100;not null;comment:邮箱"`

// 外键索引
UserID uint `json:"user_id" gorm:"index;not null;comment:用户ID"`

// 大文本
Content string `json:"content" gorm:"type:longtext;not null;comment:内容"`

// JSON字段
Dimensions []*EvaluationDimension `json:"dimensions" gorm:"type:json;serializer:json"`

// 时间戳（毫秒）
CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:milli;comment:创建时间"`

// 软删除标记（整数）
Deleted int `json:"deleted" gorm:"default:0;index:idx_deleted;comment:删除标记"`

// 密码字段（隐藏）
PasswordHash string `json:"-" gorm:"size:255;not null"`
```

**autoCreateTime 和 autoUpdateTime**:
```go
// 秒级时间戳
CreatedAt time.Time `gorm:"autoCreateTime"`

// 毫秒级时间戳
CreatedAt time.Time `gorm:"autoCreateTime:milli"`

// 纳秒级时间戳
CreatedAt time.Time `gorm:"autoCreateTime:nano"`

// 自动更新
UpdatedAt time.Time `gorm:"autoUpdateTime:milli"`
```

**我的理解**: JSON标签和GORM标签各司其职。JSON标签控制API层的数据表示（如PasswordHash使用`json:"-"`避免密码泄露），GORM标签控制数据库层的存储细节（如字段类型、索引、约束）。本项目广泛使用了注释标签`comment`，这是良好的实践，有助于数据库文档生成。毫秒时间戳`autoCreateTime:milli`提供了足够的时间精度，适合记录对话、操作日志等需要精确时间的场景。使用`serializer:json`可以直接将Go结构体存储为JSON字段，避免手动序列化/反序列化，简化代码。

---

### 4. 软删除

```go
// GORM方式：使用DeletedAt时间戳
type User struct {
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 应用层方式：使用整数标记
type Resume struct {
    Deleted int `json:"deleted" gorm:"default:0;index;comment:删除标记"`
}
```

**什么是软删除**:
软删除是指在删除数据时不真正从数据库中删除记录，而是通过标记字段表示数据已删除。查询时默认过滤掉已删除的数据，需要时可以恢复。

**两种实现方式**:

**1. GORM的DeletedAt方式**（User模型）:
```go
DeletedAt gorm.DeletedAt `gorm:"index"`

// 删除操作
db.Delete(&User{ID: 1})
// 实际执行: UPDATE users SET deleted_at = '2025-01-29 10:00:00' WHERE id = 1

// 查询自动过滤
db.Find(&users)
// 自动添加: WHERE deleted_at IS NULL

// 查询包含已删除记录
db.Unscoped().Find(&users)

// 永久删除
db.Unscoped().Delete(&User{ID: 1})
```

**2. 整数标记方式**（Resume、InterviewEvaluation等）:
```go
Deleted int `gorm:"default:0;index"`

// 删除操作（应用层）
db.Model(&Resume{}).Where("id = ?", id).Update("deleted", 1)

// 查询时手动过滤
db.Where("deleted = 0").Find(&resumes)

// 恢复（应用层）
db.Model(&Resume{}).Where("id = ?", id).Update("deleted", 0)
```

**优点**:
1. **数据安全**: 防止误删除，可以恢复数据
2. **审计追踪**: 保留历史记录，用于审计和分析
3. **关联完整性**: 避免级联删除导致的数据丢失
4. **业务逻辑**: 某些业务需要"删除"但保留数据（如订单、交易记录）
5. **性能优化**: 延迟删除，批量清理过期数据

**缺点**:
1. **存储空间**: 占用更多磁盘空间
2. **查询复杂**: 需要添加过滤条件
3. **唯一索引**: 需要特殊处理（如email唯一但允许多个已删除用户使用相同email）

**本项目的使用**:
- **User**: 使用GORM的DeletedAt，自动处理
- **Resume、InterviewEvaluation、AnswerReport**: 使用整数标记Deleted，应用层控制
- **InterviewRecord、InterviewDialogue**: 不使用软删除，可以真正删除

**唯一索引与软删除**:
```go
// 问题：email唯一索引会阻止重复邮箱（即使已删除）
Email string `gorm:"uniqueIndex"`

// 解决方案1：组合唯一索引（email + deleted_at）
Email string `gorm:"uniqueIndex:idx_email_deleted"`
DeletedAt gorm.DeletedAt `gorm:"uniqueIndex:idx_email_deleted"`

// 解决方案2：删除时修改email（添加时间戳后缀）
// email: user@example.com → user@example.com_deleted_1643000000
```

**我的理解**: 本项目混合使用了两种软删除方式。User使用GORM的DeletedAt时间戳，享受GORM的自动查询过滤便利。Resume等模型使用整数标记，提供更灵活的控制（如区分不同的删除状态：0-正常，1-用户删除，2-管理员删除，3-系统清理等）。软删除是双刃剑，增加了数据安全性但也带来了查询复杂度。对于不需要恢复的数据（如临时会话、日志），可以使用硬删除节省空间。生产环境应该定期归档或清理长期软删除的数据。

---

## 数据验证

### 验证规则

记录各模型的验证规则：

```go
// 示例：用户注册请求验证
type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6,max=100"`
}

// 面试创建请求验证
type CreateInterviewRequest struct {
    Type       string `json:"type" binding:"required,oneof=comprehensive technical"`
    Difficulty string `json:"difficulty" binding:"required,oneof=easy medium hard"`
    Domain     string `json:"domain" binding:"required,max=255"`
}

// 简历上传验证
type UploadResumeRequest struct {
    FileName string `json:"file_name" binding:"required,max=255"`
    FileSize int64  `json:"file_size" binding:"required,max=10485760"` // 10MB
    FileType string `json:"file_type" binding:"required,oneof=pdf doc docx txt"`
    Content  string `json:"content" binding:"required"`
}
```

**常用验证标签**:

| 标签 | 作用 | 示例 |
|------|------|------|
| `required` | 字段必填 | `binding:"required"` |
| `email` | 邮箱格式验证 | `binding:"required,email"` |
| `min` | 最小值/最小长度 | `binding:"min=6"` |
| `max` | 最大值/最大长度 | `binding:"max=100"` |
| `len` | 精确长度 | `binding:"len=11"` |
| `gt` | 大于 | `binding:"gt=0"` |
| `gte` | 大于等于 | `binding:"gte=18"` |
| `lt` | 小于 | `binding:"lt=100"` |
| `lte` | 小于等于 | `binding:"lte=150"` |
| `oneof` | 枚举值 | `binding:"oneof=male female"` |
| `url` | URL格式 | `binding:"url"` |
| `uuid` | UUID格式 | `binding:"uuid"` |
| `alphanum` | 只允许字母和数字 | `binding:"alphanum"` |
| `numeric` | 只允许数字 | `binding:"numeric"` |

**自定义验证器**:
```go
import "github.com/go-playground/validator/v10"

// 注册自定义验证器
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterValidation("username", validateUsername)
}

// 自定义验证函数
func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // 用户名只能包含字母、数字、下划线
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
    return matched
}

// 使用自定义验证
type User struct {
    Username string `binding:"required,username,min=3,max=50"`
}
```

**模型层验证**:
```go
// 在模型中添加验证逻辑
func (r *Resume) Validate() error {
    if r.FileSize > 10*1024*1024 {
        return errors.New("文件大小超过10MB")
    }
    if r.FileType != "pdf" && r.FileType != "doc" && r.FileType != "docx" {
        return errors.New("不支持的文件类型")
    }
    if len(r.Content) == 0 {
        return errors.New("简历内容不能为空")
    }
    return nil
}
```

---

## 索引设计

### 需要索引的字段

记录哪些字段需要索引及原因：

| 表名 | 字段 | 索引类型 | 原因 |
|------|------|---------|------|
| user | email | 唯一索引 | 登录时根据email查询用户，需要快速查询且保证唯一性 |
| user | username | 唯一索引 | 用户名登录和唯一性检查 |
| user | wechat_open_id | 唯一索引 | 微信登录时根据OpenID查询用户 |
| user | deleted_at | 普通索引 | 软删除查询过滤 |
| resume | user_id | 普通索引 | 频繁根据用户ID查询用户的所有简历 |
| resume | deleted | 普通索引 | 查询时需要过滤已删除的简历 |
| interview_record | user_id | 普通索引 | 查询用户的面试历史，高频操作 |
| interview_record | status | 普通索引 | 按状态筛选面试（如查询进行中的面试） |
| interview_dialogue | user_id | 普通索引(idx_user_id) | 查询用户的对话记录 |
| interview_dialogue | report_id | 普通索引(idx_report_id) | 查询某次面试的所有对话，高频操作 |
| interview_evaluation | user_id | 普通索引(idx_user_id) | 查询用户的评估记录 |
| interview_evaluation | report_id | 普通索引(idx_report_id) | 根据面试ID查询评估，1:1关系 |
| interview_evaluation | deleted | 普通索引(idx_deleted) | 过滤已删除的评估 |
| answer_report | user_id | 普通索引(idx_user_id) | 查询用户的答题报告 |
| answer_report | report_id | 普通索引(idx_report_id) | 根据面试ID查询报告 |
| answer_report | deleted | 普通索引(idx_deleted) | 过滤已删除的报告 |
| prediction_record | user_id | 普通索引 | 查询用户的押题历史 |
| prediction_record | resume_id | 普通索引 | 查询某份简历的押题记录 |
| prediction_question | record_id | 普通索引 | 查询某次押题的所有题目 |
| user_model | user_id | 普通索引 | 查询用户的模型配置列表 |

**索引设计原则**:

1. **高频查询字段**: 如user_id、email等经常作为WHERE条件的字段
2. **外键字段**: 关联查询时提高性能
3. **唯一性约束**: 如email、username需要保证唯一
4. **软删除标记**: deleted、deleted_at等查询过滤字段
5. **状态字段**: status等按状态筛选的字段
6. **时间字段**: created_at等按时间范围查询的字段

**复合索引**:
```go
// 示例：按用户ID和状态查询面试
type InterviewRecord struct {
    UserID uint   `gorm:"index:idx_user_status"`
    Status string `gorm:"index:idx_user_status"`
}

// 生成SQL: CREATE INDEX idx_user_status ON interview_record(user_id, status)
// 适用查询: WHERE user_id = ? AND status = ?
```

**索引优化建议**:
- 避免过多索引：每个索引都会增加写入开销
- 选择性高的字段：区分度高的字段索引效果更好
- 考虑查询模式：根据实际SQL语句设计索引
- 定期分析：使用EXPLAIN分析查询性能

**我的理解**: 索引是数据库性能优化的关键。本项目在所有外键字段（user_id、resume_id、report_id）都建立了索引，因为这些是高频关联查询的关键。唯一索引（email、username）既保证数据唯一性又加速查询。软删除相关的字段（deleted、deleted_at）都加了索引，因为几乎每次查询都需要过滤已删除数据。命名索引（如idx_user_id）比自动生成的索引名更清晰，便于维护。生产环境应该定期监控慢查询，根据实际查询模式调整索引策略。

---

## 模型方法

### 自定义方法

记录模型上定义的自定义方法：

```go
// 表名定义
func (User) TableName() string {
    return "user"
}

func (Resume) TableName() string {
    return "resume"
}

// 模型验证
func (r *Resume) Validate() error {
    if r.FileSize > 10*1024*1024 {
        return errors.New("文件大小不能超过10MB")
    }
    return nil
}

// 密码加密（可以在BeforeCreate钩子中调用）
func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(u.PasswordHash), 
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    u.PasswordHash = string(hashedPassword)
    return nil
}

// 密码验证
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.PasswordHash), 
        []byte(password),
    )
    return err == nil
}
```

**钩子函数**（GORM Hooks）:

```go
// BeforeCreate - 创建记录前执行
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 密码加密
    if u.PasswordHash != "" {
        return u.HashPassword()
    }
    return nil
}

// AfterCreate - 创建记录后执行
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 发送欢迎邮件
    // go sendWelcomeEmail(u.Email)
    return nil
}

// BeforeUpdate - 更新记录前执行
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 检查是否需要重新加密密码
    if tx.Statement.Changed("PasswordHash") {
        return u.HashPassword()
    }
    return nil
}

// AfterUpdate - 更新记录后执行
func (u *User) AfterUpdate(tx *gorm.DB) error {
    // 清除缓存
    // cache.Delete(fmt.Sprintf("user:%d", u.ID))
    return nil
}

// BeforeDelete - 删除记录前执行
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // 检查是否有关联数据
    var count int64
    tx.Model(&Resume{}).Where("user_id = ?", u.ID).Count(&count)
    if count > 0 {
        return errors.New("用户还有简历，不能删除")
    }
    return nil
}

// AfterDelete - 删除记录后执行
func (u *User) AfterDelete(tx *gorm.DB) error {
    // 清除相关数据
    // tx.Where("user_id = ?", u.ID).Delete(&Session{})
    return nil
}

// AfterFind - 查询记录后执行
func (u *User) AfterFind(tx *gorm.DB) error {
    // 加载额外数据或转换格式
    return nil
}
```

**钩子执行顺序**:

**创建**:
```
Begin Transaction
  BeforeSave
    BeforeCreate
      // 保存到数据库
    AfterCreate
  AfterSave
Commit / Rollback
```

**更新**:
```
Begin Transaction
  BeforeSave
    BeforeUpdate
      // 更新数据库
    AfterUpdate
  AfterSave
Commit / Rollback
```

**删除**:
```
Begin Transaction
  BeforeDelete
    // 从数据库删除
  AfterDelete
Commit / Rollback
```

**钩子使用注意事项**:
- 钩子中的错误会导致事务回滚
- 不要在钩子中执行耗时操作，应该使用异步
- BeforeCreate/BeforeUpdate中可以修改字段值
- 使用tx.Statement.Changed()检查字段是否变更

**我的理解**: GORM的钩子函数提供了在数据库操作的关键节点插入业务逻辑的机制。常见用途包括：1) 数据验证（BeforeCreate/BeforeUpdate）；2) 密码加密（BeforeCreate/BeforeUpdate）；3) 自动填充字段（BeforeCreate设置UUID）；4) 缓存失效（AfterUpdate/AfterDelete）；5) 关联数据处理。本项目主要在DAO层手动控制逻辑，未大量使用钩子，这样保持了更明确的控制流。钩子适合处理通用的、重复的逻辑，但要注意不要过度使用导致代码难以跟踪。

---

## 数据库迁移

### 自动迁移

```go
// backend/main.go 中的自动迁移
db.AutoMigrate(
    &model.User{},
    &model.Resume{},
    &model.InterviewRecord{},
    &model.InterviewDialogue{},
    &model.InterviewEvaluation{},
    &model.AnswerReport{},
    &model.PredictionRecord{},
    &model.PredictionQuestion{},
    &model.UserModel{},
)
```

**AutoMigrate做了什么**:
1. **创建表**: 如果表不存在，根据结构体定义创建表
2. **添加字段**: 如果结构体新增了字段，会自动添加到表中
3. **创建索引**: 根据gorm标签创建相应的索引
4. **创建外键**: 根据关联定义创建外键约束
5. **修改字段类型**: 如果字段类型变更，尝试修改（有限制）

**AutoMigrate会做**:
- ✅ 创建新表
- ✅ 添加新字段
- ✅ 创建缺失的索引
- ✅ 创建外键约束
- ✅ 修改字段大小（增加）

**AutoMigrate不会做**:
- ❌ 删除未使用的字段
- ❌ 修改字段类型（不兼容的变更）
- ❌ 删除索引
- ❌ 重命名字段或表
- ❌ 删除表

**限制**:
1. **无法删除字段**: 如果结构体中删除了字段，数据库表中的字段不会被删除
2. **类型变更限制**: 不兼容的类型变更可能失败
3. **数据迁移**: 不处理数据转换和迁移
4. **复杂变更**: 多步骤变更无法自动处理

**生产环境建议**:

**开发阶段**: 可以使用AutoMigrate快速迭代
```go
if os.Getenv("ENV") == "development" {
    db.AutoMigrate(&model.User{}, &model.Resume{}, ...)
}
```

**生产环境**: 建议使用专业迁移工具

**1. golang-migrate**:
```bash
# 安装
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 创建迁移文件
migrate create -ext sql -dir db/migrations -seq create_users_table

# 执行迁移
migrate -path db/migrations -database "mysql://user:pass@tcp(host:3306)/dbname" up

# 回滚
migrate -path db/migrations -database "mysql://user:pass@tcp(host:3306)/dbname" down 1
```

迁移文件示例 (↑上迹):
```sql
-- 000001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_email` (`email`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

迁移文件示例 (↓下迹):
```sql
-- 000001_create_users_table.down.sql
DROP TABLE IF EXISTS `user`;
```

**2. goose**:
```bash
# 安装
go install github.com/pressly/goose/v3/cmd/goose@latest

# 创建迁移
goose -dir db/migrations create add_user_avatar sql

# 执行迁移
goose -dir db/migrations mysql "user:pass@/dbname" up
```

**3. GORM Migrator（编程式）**:
```go
// 自定义迁移逻辑
type Migration struct {
    ID        uint
    Version   string
    AppliedAt time.Time
}

func RunMigrations(db *gorm.DB) error {
    // 检查迁移状态
    if !db.Migrator().HasTable(&Migration{}) {
        db.AutoMigrate(&Migration{})
    }
    
    // 执行具体迁移
    migrations := []func(*gorm.DB) error{
        migration_001_create_users,
        migration_002_add_user_avatar,
    }
    
    for _, m := range migrations {
        if err := m(db); err != nil {
            return err
        }
    }
    return nil
}
```

**迁移最佳实践**:
1. **版本控制**: 迁移文件纳入Git版本控制
2. **可逆性**: 每个迁移都应有up和down脚本
3. **幂等性**: 多次执行同一迁移结果一致
4. **小步快跑**: 每次迁移只做一件事
5. **测试**: 在测试环境充分测试后再上线
6. **备份**: 执行迁移前备份数据库

**我的理解**: AutoMigrate适合开发阶段快速迭代，但生产环境应该使用专业的迁移工具。原因：1) 版本控制：明确记录每个变更；2) 可逆性：支持回滚操作；3) 审计：所有变更有迹可查；4) 团队协作：迁移文件统一管理；5) 安全：避免意外的结构变更。本项目目前使用AutoMigrate，在项目成熟后应该迁移到golang-migrate或goose。

---

## ❓ 问题记录

### Q1: 为什么使用软删除而不是硬删除？
**我的理解**: 软删除有多个优势：1) **数据安全** - 防止误删除，可以恢复数据；2) **审计追踪** - 保留历史记录用于审计和数据分析；3) **关联完整性** - 避免级联删除导致关联数据丢失；4) **业务需求** - 某些业务场景需要“删除”但保留数据（如订单、交易记录）。但软删除也有缺点：增加存储空间、查询需要额外过滤条件、唯一索引需要特殊处理。对于不需要恢复的数据（如日志、临时数据）可以使用硬删除节省空间。

### Q2: 外键约束在应用层还是数据库层？
**我的理解**: 本项目主要采用**应用层外键**的策略。只有PredictionQuestion明确使用了数据库外键（CASCADE级联删除），因为题目完全依赖押题记录存在。应用层外键的优势：1) **灵活性** - 可以实现更复杂的关联逻辑；2) **性能** - 避免数据库外键锁，提高并发性能；3) **扩展性** - 便于分库分表、跨数据库查询；4) **可控** - 删除逻辑由业务代码控制，更加明确。数据库外键的优势：1) **数据一致性** - 数据库层保证强制约束；2) **简化代码** - 不需要在业务代码中处理级联操作。选择哪种方式取决于项目规模和复杂度：小型项目可用数据库外键，大型项目建议应用层外键。

### Q3: 为什么InterviewEvaluation使用JSON字段存储Dimensions？
**我的理解**: 使用JSON字段存储评估维度是**灵活性与结构化的权衡**。优势：1) **灵活性** - 可以动态调整评估维度，不同面试类型可以有不同的评估标准；2) **简化表结构** - 避免创建Dimension子表和复杂的关联查询；3) **原子性** - 评估和维度作为一个整体读写，无需事务处理多表；4) **性能** - 一次查询获取所有信息，不需要JOIN。缺点：1) **查询复杂** - 无法直接按维度查询排序（需要JSON函数）；2) **索引限制** - JSON字段的索引支持比普通字段弱；3) **数据验证** - 需要在应用层验证JSON结构。MySQL 5.7+对JSON的支持已经很完善，这种设计是合理的。如果未来需要复杂的维度统计分析，可以考虑拆分为独立表。

### Q4: uint64和uint的选择？
**我的理解**: 本项目混合使用了uint和uint64作为ID类型。**uint**（User.ID）适用于用户表等增长缓慢的表，最大4亿+记录；**uint64**（InterviewRecord.ID、Resume.ID）适用于高频创建的表，支持18京亿+记录。选择因素：1) **数据增长速度** - 面试记录、对话记录增长快，需要更大的ID范围；2) **未来扩展** - 为大规模应用预留空间；3) **一致性** - 同类表使用相同类型。但这也带来了不一致：User.ID是uint但Resume.UserID也是uint，这是合理的，因为外键必须与主键类型一致。现代应用中，直接使用uint64作为所有ID类型是更好的实践，避免类型转换问题。

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
