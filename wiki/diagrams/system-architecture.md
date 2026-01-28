# 系统架构图

> 系统的各种架构图和流程图

## 整体架构图

```mermaid
graph TB
    subgraph "客户端层"
        A[Web浏览器]
        B[移动端]
    end
    
    subgraph "前端层"
        C[Next.js应用]
    end
    
    subgraph "API网关层"
        D[Hertz Web框架]
        D1[CORS中间件]
        D2[JWT认证中间件]
        D3[日志中间件]
    end
    
    subgraph "业务逻辑层"
        E[Handler层]
        F[Service层]
        G[AI Agent层]
    end
    
    subgraph "数据访问层"
        H[Repository层]
        I[Cache层]
    end
    
    subgraph "数据存储层"
        J[(MySQL)]
        K[(Redis)]
        L[(Milvus)]
    end
    
    subgraph "外部服务"
        M[OpenAI API]
        N[微信API]
    end
    
    A --> C
    B --> C
    C --> D
    D --> D1
    D1 --> D2
    D2 --> D3
    D3 --> E
    E --> F
    F --> G
    F --> H
    H --> I
    I --> K
    H --> J
    G --> M
    G --> L
    F --> N
```

**说明**: 补充你对架构图的理解_____

---

## AI工作流图

```mermaid
graph LR
    A[用户输入] --> B[路由到对应Agent]
    B --> C{判断Agent类型}
    
    C -->|面试| D[Interview Agent]
    C -->|简历| E[Resume Agent]
    C -->|评估| F[Evaluation Agent]
    C -->|预测| G[Prediction Agent]
    
    D --> H[加载上下文]
    E --> H
    F --> H
    G --> H
    
    H --> I{是否需要工具}
    
    I -->|是| J[调用Tool]
    I -->|否| K[直接调用LLM]
    
    J --> L[简历检索Tool]
    J --> M[题库检索Tool]
    J --> N[Milvus检索Tool]
    
    L --> K
    M --> K
    N --> K
    
    K --> O[LLM处理]
    O --> P[流式返回]
    P --> Q[保存上下文]
    Q --> R[返回结果]
```

**说明**: _____

---

## 数据流图

```mermaid
sequenceDiagram
    participant U as 用户
    participant F as 前端
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as 数据库
    participant Cache as Redis
    
    U->>F: 发起请求
    F->>H: HTTP请求
    H->>H: 参数验证
    H->>S: 调用服务
    S->>R: 查询数据
    
    alt 缓存命中
        R->>Cache: 查询缓存
        Cache-->>R: 返回缓存数据
    else 缓存未命中
        R->>DB: 查询数据库
        DB-->>R: 返回数据
        R->>Cache: 写入缓存
    end
    
    R-->>S: 返回数据
    S->>S: 业务逻辑处理
    S-->>H: 返回结果
    H-->>F: HTTP响应
    F-->>U: 展示结果
```

**说明**: _____

---

## 面试流程图

```mermaid
stateDiagram-v2
    [*] --> 创建面试
    创建面试 --> 加载简历
    加载简历 --> 开始面试
    
    开始面试 --> 生成问题
    生成问题 --> 等待回答
    等待回答 --> 接收回答
    
    接收回答 --> 分析回答
    分析回答 --> 决策
    
    决策 --> 追问: 需要深入
    决策 --> 生成问题: 切换话题
    决策 --> 结束面试: 达到限制
    
    追问 --> 等待回答
    
    结束面试 --> 生成评估
    生成评估 --> 保存记录
    保存记录 --> [*]
```

**说明**: _____

---

## 数据库ER图

```
┌─────────────┐
│    users    │
├─────────────┤
│ id (PK)     │
│ username    │
│ email       │
│ password    │
│ created_at  │
│ updated_at  │
└──────┬──────┘
       │ 1
       │
       │ N
┌──────┴──────────┐         ┌─────────────────┐
│    resumes      │         │ interview_records│
├─────────────────┤         ├──────────────────┤
│ id (PK)         │         │ id (PK)          │
│ user_id (FK)    │         │ user_id (FK)     │
│ title           │         │ type             │
│ content         │         │ status           │
│ file_url        │         │ start_time       │
│ created_at      │         │ end_time         │
└────────┬────────┘         └────────┬─────────┘
         │ 1                         │ 1
         │                           │
         │ N                         │ N
┌────────┴─────────┐         ┌───────┴──────────────┐
│   predictions    │         │ interview_dialogues  │
├──────────────────┤         ├──────────────────────┤
│ id (PK)          │         │ id (PK)              │
│ resume_id (FK)   │         │ interview_id (FK)    │
│ category         │         │ role                 │
│ questions        │         │ content              │
│ created_at       │         │ sequence             │
└──────────────────┘         └──────────────────────┘
                                     
                             ┌──────────────────────┐
                             │ interview_evaluations│
                             ├──────────────────────┤
                             │ id (PK)              │
                             │ interview_id (FK)    │
                             │ total_score          │
                             │ feedback             │
                             │ created_at           │
                             └──────────────────────┘
```

**说明**: _____

---

## 部署架构图

```
┌─────────────────────────────────────────┐
│            Load Balancer                 │
└──────────────┬──────────────────────────┘
               │
       ┌───────┴────────┐
       │                │
┌──────▼──────┐  ┌──────▼──────┐
│   Nginx 1   │  │   Nginx 2   │
└──────┬──────┘  └──────┬──────┘
       │                │
       └───────┬────────┘
               │
       ┌───────┴────────┐
       │                │
┌──────▼──────┐  ┌──────▼──────┐
│  Backend 1  │  │  Backend 2  │
│  (Hertz)    │  │  (Hertz)    │
└──────┬──────┘  └──────┬──────┘
       │                │
       └───────┬────────┘
               │
       ┌───────┴────────┬────────┐
       │                │        │
┌──────▼──────┐  ┌──────▼──┐  ┌─▼────────┐
│   MySQL     │  │  Redis  │  │  Milvus  │
│   Master    │  │ Cluster │  │          │
└──────┬──────┘  └─────────┘  └──────────┘
       │
┌──────▼──────┐
│   MySQL     │
│   Slave     │
└─────────────┘
```

**说明**: _____

---

## 你的架构图

请在这里补充你自己画的架构图和理解：

### 我理解的系统架构

```
在这里画出你的理解
```

### 我理解的数据流

```
在这里画出你的理解
```

### 我理解的AI处理流程

```
在这里画出你的理解
```

---

## 工具推荐

### 画图工具
- [Mermaid](https://mermaid.js.org/) - 代码画流程图
- [Draw.io](https://www.drawio.com/) - 在线画图工具
- [PlantUML](https://plantuml.com/) - UML图工具
- [Excalidraw](https://excalidraw.com/) - 手绘风格画图

### 使用Mermaid

在Markdown中直接写：

````markdown
```mermaid
graph TD
    A[开始] --> B[结束]
```
````

---

**返回**: [Wiki首页](../README.md)
