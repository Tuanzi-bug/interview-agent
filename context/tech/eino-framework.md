# Eino 框架集成指南

**最后更新**: 2026-02-01  
**维护者**: AI 架构团队

## 概述

Eino 是字节跳动开源的大语言模型应用框架，为构建 LLM 应用提供了完整的解决方案。在面试吧平台中，Eino 用于实现 AI 面试功能。

### 核心特性

- 🏗️ **组件化设计**: 模块化组件易于复用和扩展
- 📊 **图编排引擎**: 支持复杂的 AI 流程编排
- 🔄 **流式处理**: 原生支持流式输出，改善 UX
- 🧠 **Agent 框架**: 内置 Agent 框架简化开发
- 🔌 **多模型支持**: 集成多个 LLM 提供商

## 项目中的 Eino 应用

### 简历分析 Agent

**职责**: 深度分析用户简历，提供专业建议

**技术实现**:
```go
// 定义 Resume Analyzer Agent
type ResumeAnalyzerAgent struct {
    llm      model.LanguageModel
    tools    []tool.Tool
    retriever vector.Retriever
}

// Agent 执行流程
1. 接收简历文本
2. 使用向量检索获取相关知识库
3. 调用 LLM 进行分析
4. 生成构造化分析结果
```

**关键功能**:
- 简历结构化解析
- 技能点识别和评估
- 职业发展路径分析
- 面试重点预测

### 面试问题生成 Agent

**职责**: 根据面试类型和难度生成个性化面试题

**技术实现**:
```go
// Question Generator Agent
1. 获取用户简历和选择的职位类型
2. 从题库检索相关问题模板
3. 基于用户背景定制问题
4. 调整难度和覆盖面试重点
```

**支持的面试类型**:
- 综合面试 (通用技能、项目经验)
- 专项面试 (Java、Go、MySQL、Redis 等)
- 行为面试 (STAR 方法)
- 系统设计面试

### 答案评估 Agent

**职责**: 评估用户答案质量，提供详细反馈

**技术实现**:
```go
// Answer Evaluator Agent
1. 接收用户答案和标准答案/评估标准
2. 调用 LLM 进行多维度评估
3. 生成评分、缺点分析、改进建议
4. 返回结构化评估结果
```

**评估维度**:
- 正确性 (Correctness)
- 完整性 (Completeness)
- 清晰度 (Clarity)
- 深度 (Depth)
- 示例质量 (Example Quality)

## Eino 架构基础

### 核心组件

```
┌─────────────────────────────────────┐
│         Eino 框架                    │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐  │
│  │ Chain - 链式组合             │  │
│  │ (组合多个组件完成任务)       │  │
│  └────────────┬────────────────┘  │
│               │                    │
│  ┌────────────▼────────────────┐  │
│  │ Component - 组件             │  │
│  │ - Model (LLM)               │  │
│  │ - Tool (工具)               │  │
│  │ - Retriever (检索器)         │  │
│  │ - Transformer (转换器)       │  │
│  └─────────────────────────────┘  │
│               │                    │
│  ┌────────────▼────────────────┐  │
│  │ Graph - 图编排               │  │
│  │ (节点和边的编排)             │  │
│  └─────────────────────────────┘  │
│                                     │
└─────────────────────────────────────┘
```

### 主要接口

#### 1. LanguageModel (LLM 模型)

```go
type LanguageModel interface {
    // 调用 LLM 获取响应
    Generate(ctx context.Context, request *GenerateRequest) (*GenerateResponse, error)
    
    // 流式调用
    GenerateStream(ctx context.Context, request *GenerateRequest, handler StreamHandler) error
    
    // 多轮对话
    Chat(ctx context.Context, messages []Message) (*ChatResponse, error)
}

// 支持的模型
- OpenAI ChatGPT
- ARK (字节内部模型)
- 其他兼容 API 的模型
```

#### 2. Tool (工具)

```go
type Tool interface {
    // 工具名称
    Name() string
    
    // 工具描述
    Description() string
    
    // 工具参数定义
    InputSchema() interface{}
    
    // 执行工具
    Execute(ctx context.Context, input interface{}) (interface{}, error)
}

// 常用工具
- WebSearch (网页搜索)
- Calculator (计算器)
- CustomTool (自定义工具)
```

#### 3. Retriever (检索器)

```go
type Retriever interface {
    // 检索相似文档
    Retrieve(ctx context.Context, query string, opts ...Option) ([]*Document, error)
    
    // 向量检索
    RetrieveVectors(ctx context.Context, vectors [][]float32) ([]*Document, error)
}

// 在项目中
- Milvus 检索器 (向量数据库)
- 知识库检索器
```

## Chain 组合模式

### 简单 Chain 示例

```go
// 1. 创建 LLM 模型
llm := openai.NewModel(apiKey)

// 2. 创建 Chain
chain := chain.NewChain(
    // 输入处理
    chain.NewInput("question"),
    
    // LLM 处理
    chain.NewLLMStep(llm, "Answer this: {question}"),
    
    // 输出处理
    chain.NewOutput("answer"),
)

// 3. 运行 Chain
result := chain.Run(ctx, map[string]interface{}{
    "question": "Go 的 goroutine 是什么?",
})

fmt.Println(result["answer"])
// 输出: Goroutine 是 Go 语言中的轻量级线程...
```

### RAG Chain (检索增强生成)

```go
// 检索增强生成模式
chain := chain.NewChain(
    // 输入
    chain.NewInput("query"),
    
    // 从向量库检索相关文档
    chain.NewRetrieverStep(retriever, "{query}", k=3),
    
    // 构建增强提示
    chain.NewTransformerStep(func(context map[string]interface{}) map[string]interface{} {
        docs := context["retrieved_docs"]
        return map[string]interface{}{
            "context": FormatDocuments(docs),
            "query": context["query"],
        }
    }),
    
    // 使用 LLM 生成答案
    chain.NewLLMStep(llm, 
        "Context: {context}\nQuestion: {query}\nAnswer:"),
    
    // 输出
    chain.NewOutput("answer"),
)
```

## Agent 开发

### Agent 基础框架

```go
type Agent struct {
    model     model.LanguageModel
    tools     map[string]tool.Tool
    retriever retriever.Retriever
    memory    memory.Memory
}

// Agent 执行流程
func (a *Agent) Execute(ctx context.Context, task string) (interface{}, error) {
    // 1. 思考 - LLM 理解任务
    thought := a.model.Generate(ctx, task)
    
    // 2. 行动 - 选择工具或检索
    action := ParseAction(thought)
    
    // 3. 观察 - 执行工具获取结果
    observation := ExecuteTool(action)
    
    // 4. 反思 - 是否完成或继续循环
    if IsComplete(observation) {
        return observation, nil
    }
    
    // 继续循环
    return a.Execute(ctx, FormulateNextTask(thought, observation))
}
```

### 简历分析 Agent 实现

```go
type ResumeAnalyzerAgent struct {
    llm       model.LanguageModel
    tools     map[string]tool.Tool
    retriever retriever.Retriever
}

func (ra *ResumeAnalyzerAgent) Analyze(ctx context.Context, resume string) (*Analysis, error) {
    // 1. 解析简历结构
    structured := ParseResumeStructure(resume)
    
    // 2. 调用工具进行分析
    skillAnalysis := ra.AnalyzeSkills(ctx, structured.Skills)
    careerPath := ra.PredictCareerPath(ctx, structured.Experience)
    interviewTopics := ra.GenerateInterviewTopics(ctx, structured)
    
    // 3. 生成综合分析
    analysis := &Analysis{
        SkillAnalysis:  skillAnalysis,
        CareerPath:     careerPath,
        InterviewTopics: interviewTopics,
        Recommendations: ra.GenerateRecommendations(ctx, skillAnalysis),
    }
    
    return analysis, nil
}

// 技能分析工具
func (ra *ResumeAnalyzerAgent) AnalyzeSkills(ctx context.Context, skills []string) *SkillAnalysis {
    // 使用向量检索获取行业标准
    industryStandard := ra.retriever.Retrieve(ctx, "industry standard skills")
    
    // 调用 LLM 进行分析
    prompt := fmt.Sprintf(`
    分析以下技能，与行业标准比较:
    用户技能: %v
    行业标准: %v
    
    提供:
    1. 技能评估 (掌握程度)
    2. 缺失技能
    3. 优势技能
    `, skills, industryStandard)
    
    resp := ra.llm.Generate(ctx, prompt)
    return ParseSkillAnalysis(resp)
}
```

## 流式处理

### 为什么使用流式?

✅ **改善用户体验**: 不用等待完整响应，即时看到输出  
✅ **降低延迟感**: 长文本分块显示  
✅ **内存效率**: 逐块处理，不缓存完整响应

### 流式实现

```go
// 前端 WebSocket 连接
func (h *Handler) StreamAnswer(c context.Context) {
    // 创建流处理器
    streamHandler := func(chunk string) {
        // 发送到前端 WebSocket
        h.websocket.Send(chunk)
    }
    
    // 调用 LLM 流式生成
    err := h.llm.GenerateStream(c, prompt, streamHandler)
}

// 后端 Hertz 处理
func StreamHandler(ctx *context.Context) {
    // 建立 SSE 连接
    ctx.SetHeader("Content-Type", "text/event-stream")
    
    // 流式生成答案
    h.evaluator.EvaluateStream(ctx, answer, func(chunk string) {
        fmt.Fprintf(ctx, "data: %s\n\n", chunk)
        ctx.Writer.Flush()
    })
}
```

## 模型集成

### 支持的模型

#### OpenAI

```go
import "github.com/cloudwego/eino-ext/components/model/openai"

model := openai.NewModel(&openai.Config{
    APIKey: os.Getenv("OPENAI_API_KEY"),
    Model:  "gpt-4",
})

resp := model.Generate(ctx, "Your prompt")
```

#### ARK (字节内部模型)

```go
import "github.com/cloudwego/eino-ext/components/model/ark"

model := ark.NewModel(&ark.Config{
    APIKey: os.Getenv("ARK_API_KEY"),
    Model:  "your-model-id",
})

resp := model.Generate(ctx, "Your prompt")
```

### 添加新模型

```go
type CustomModel struct {
    apiKey string
}

func (m *CustomModel) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
    // 调用你的 API
    resp := callYourAPI(req.Prompt)
    return &GenerateResponse{
        Text: resp.Text,
    }, nil
}
```

## 错误处理和重试

### 错误类型

```go
// LLM 错误
- RateLimitError: 速率限制
- AuthenticationError: 认证失败
- ContextLengthError: 上下文过长
- APIError: API 错误

// Chain 错误
- InvalidInputError: 输入格式错误
- ExecutionError: 执行失败
- TimeoutError: 执行超时
```

### 重试策略

```go
func RetryWithBackoff(fn func() error, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if !IsRetryable(err) {
            return err
        }
        
        lastErr = err
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    
    return lastErr
}

// 使用
err := RetryWithBackoff(func() error {
    return agent.Execute(ctx, task)
}, 3)
```

## 性能优化

### 缓存策略

```go
// 缓存常见问题的答案
cache := lru.New(1000)

func (a *Agent) Execute(ctx context.Context, task string) (interface{}, error) {
    // 检查缓存
    if cached, ok := cache.Get(task); ok {
        return cached, nil
    }
    
    // 执行任务
    result := a.doExecute(ctx, task)
    
    // 存储缓存
    cache.Add(task, result)
    
    return result, nil
}
```

### Token 优化

```go
// 减少 Token 使用
func BuildPrompt(context string, query string) string {
    // 只包含必要的上下文
    return fmt.Sprintf(`Context: %s
    
Question: %s

Answer:`, TruncateContext(context, 2000), query)
}

// 使用函数调用而不是自然语言
// 减少多轮对话中的 Token 使用
```

## 常见问题

**Q: 如何处理 LLM 输出不稳定?**
→ 使用温度(temperature)参数调整随机性，使用 few-shot 提示词示例

**Q: 如何优化 Agent 的执行速度?**
→ 使用缓存、并行执行工具、优化提示词、使用流式处理

**Q: 如何集成新的 LLM 模型?**
→ 实现 LanguageModel 接口，集成 API 调用逻辑

**Q: 如何处理上下文长度限制?**
→ 使用摘要、分块处理、检索最相关片段

## 相关资源

### 官方文档
- [Eino GitHub](https://github.com/cloudwego/eino)
- [Eino 官方文档](https://eino-doc.cloudwego.dev/)

### 项目代码
- [AI Agent 实现](../../backend/chatApp/agent/)
- [AI Service 封装](../../backend/chatApp/agent_service/)
- [简历分析 Agent](../../backend/chatApp/agent/resume/)

### 相关文档
- [架构设计](architecture.md)
- [Hertz 框架](hertz-framework.md)

---

**维护者**: AI 架构团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
