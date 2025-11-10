# 面试助手系统 Agent 模块架构与功能文档

本文档详细描述了 `backend/chatApp/agent` 目录下各个Agent的功能、结构和工作流程。这些Agent共同构成了面试助手系统的核心智能处理组件。

## 目录结构

```
agent/
├── MainAgent.go          # 主Agent配置器
├── RouterAgent.go        # 智能任务分发器
├── QuestionGeneratorAgent.go  # 面试问题生成器
├── AnswerEvalAgent.go    # 答案评估器
├── resumeAgent.go         # 简历分析器
├── InterView_Procrss.go  # 面试流程协处理器
├── InterView_Report.go   # 面试报告生成器
```

## 1. MainAgent.go

### 功能描述
MainAgent 是整个Agent系统的入口点，负责创建并配置所有子Agent，建立它们之间的层级关系，并初始化运行器。

### 核心实现
```go
func NewMainAgent() *adk.Runner {
    ResumeAnalysisAgent := NewResumeAnalysisAgent()
    QuestionGeneratorAgent := NewQuestionGeneratorAgent()
    RouterAgent := NewRouterAgent()
    AnswerEvalAgent := NewAnswerEvalAgent()

    ctx := context.Background()
    a, err := adk.SetSubAgents(ctx, RouterAgent, []adk.Agent{ResumeAnalysisAgent, QuestionGeneratorAgent, AnswerEvalAgent})
    // 错误处理和日志记录
    runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: a})
    return runner
}
```

### 工作流程
1. 创建各个专业子Agent实例
2. 将RouterAgent设置为主控制器，其他Agent作为子Agent
3. 创建并返回Runner实例，用于启动Agent系统

## 2. RouterAgent.go

### 功能描述
RouterAgent 是智能任务分发器，负责分析用户请求并将其转发给最适合的专业Agent处理。

### 核心实现
```go
func NewRouterAgent() adk.Agent {
    // 创建配置和模型
    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "RouterAgent",
        Description: "智能任务分发器，将用户请求转交给最合适的专家助手。",
        Instruction: `你是一个智能任务路由器。请分析用户请求，并将其委派给最擅长的专家助手处理...`,
        Model:       chat.CreatOpenAiChatModel(ctx),
    })
    // 错误处理
    return a
}
```

### 任务分配规则
1. 简历分析请求 → ResumAnalysisAgent
2. 基于简历生成面试问题 → QuestionGeneratorAgent
3. 评估面试表现 → InterviewAgent（注：实际实现中可能使用AnswerEvalAgent）
4. 无法处理的请求 → 直接告知用户

## 3. QuestionGeneratorAgent.go

### 功能描述
QuestionGeneratorAgent 是一名资深技术面试官智能体，负责根据候选人简历背景生成面试问题。

### 核心实现
```go
func NewQuestionGeneratorAgent() adk.Agent {
    // 创建配置和模型
    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "QuestionGeneratorAgent",
        Description: "根据前面的简历去分析结果去提问问题",
        Instruction: `你是一名资深的技术面试官，负责围绕候选人的背景开展模拟面试...`,
        Model:       chat.CreatOpenAiChatModel(ctx),
        ToolsConfig: adk.ToolsConfig{
            ToolsNodeConfig: compose.ToolsNodeConfig{
                Tools: []componenttool.BaseTool{
                    tool2.CreateGenQuestionTool(),
                    tool2.NewAskForInputTool(),
                },
            },
        },
        MaxIterations: 12,
    })
    // 错误处理
    return a
}
```

### 工作流程
1. 调用 `gen_question` 工具获取问题的增强方案
2. 生成问题后调用 `ask_for_clarification` 工具获取用户回答
3. 循环以上过程，每轮生成一个问题并等待回答

### 关键约束
- 保持专业、鼓励且以候选人为中心的语气
- 明确当前面试阶段（如"技术深挖"、"总结反馈"等）

## 4. AnswerEvalAgent.go

### 功能描述
AnswerEvalAgent 负责根据简历分析和面试问题，对用户的答案进行专业评估。

### 核心实现
```go
func NewAnswerEvalAgent() adk.Agent {
    // 创建配置和模型
    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "AnswerEvalAgent",
        Description: "根据前面的简历分析和问题,对用户的答案进行评估",
        Instruction: `你是一名资深的技术面试官 你要去调用 "answer_eval" 工具去获取评估的方向生成评估结果...`,
        Model:       chat.CreatOpenAiChatModel(ctx),
        ToolsConfig: adk.ToolsConfig{
            ToolsNodeConfig: compose.ToolsNodeConfig{
                Tools: []componenttool.BaseTool{
                    tool2.CreateAnswerEvalTool(),
                },
            },
        },
        MaxIterations: 12,
    })
    // 错误处理
    return a
}
```

### 工作流程
1. 调用 `answer_eval` 工具获取评估的方向
2. 根据评估方向生成评估结果
3. 返回评估结果给用户

### 关键约束
- 评估结果必须使用中文

## 5. resumAgent.go

### 功能描述
ResumAnalysisAgent 是简历分析专家智能体，负责解析用户提供的简历（包括PDF格式），并进行结构化评估。

### 核心实现
```go
func NewResumAnalysisAgent() adk.Agent {
    // 创建配置和模型
    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "ResumAnalysisAgent",
        Description: "一个可以解析简历pdf分析简历的智能体",
        Instruction: `你是一名资深的简历分析专家，负责对用户的简历进行分析,并输出对应的分析结果...`,
        Model:       chat.CreatOpenAiChatModel(ctx),
        ToolsConfig: adk.ToolsConfig{
            ToolsNodeConfig: compose.ToolsNodeConfig{
                Tools: []tool.BaseTool{tool2.CreatePDFToTextTool()},
            },
        },
        MaxIterations: 10,
    })
    // 错误处理
    return a
}
```

### 工作流程
1. 当用户提供简历时，使用 `pdf_to_text` 工具提取文本（如需要）
2. 从多个维度对简历进行结构化评估：
   - 模块完整度
   - 技能匹配度
   - 量化成果
   - 语言表达等
3. 提供详细反馈和可执行的改进建议
4. 进行简历评分（0-100分）

## 6. InterView_Procrss.go

### 功能描述
InterView_Procrss 是面试流程协处理器，负责协调整个面试流程，包括简历分析、问题生成、获取回答、评估回答和生成报告。

### 核心实现
```go
func NewInterviewProcessAgent() *adk.Runner {
    // 创建各个专业子Agent
    ResumeAnalysisAgent := NewResumAnalysisAgent()
    QuestionGeneratorAgent := NewQuestionGeneratorAgent()
    AnswerEvalAgent := NewAnswerEvalAgent()
    InterviewReportAgent := NewInterviewReportAgent()

    // 创建循环Agent处理问答和评估环节
    loopAgent, err := adk.NewLoopAgent(ctx, &adk.LoopAgentConfig{
        Name:        "InterviewProcessAgent",
        Description: "根据最大迭代次数 生成问题 → 向用户获取回答 → 评估回答 ",
        SubAgents:     []adk.Agent{QuestionGeneratorAgent, AnswerEvalAgent},
        MaxIterations: 10,
    })

    // 创建顺序Agent定义完整面试流程
    sequentialAgent, err := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
        Name:        "InterviewProcessAgent",
        Description: " 模拟面试流程: 分析简历 -> 生成问题 -> 向用户获取回答 -> 评估回答 -> 生成面试报告",
        SubAgents:   []adk.Agent{ResumeAnalysisAgent, loopAgent, InterviewReportAgent},
    })
    // 错误处理和日志记录
    runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: sequentialAgent})
    return runner
}
```

### 工作流程
1. 创建所有必要的子Agent
2. 配置LoopAgent处理问答循环（生成问题→获取回答→评估回答）
3. 配置SequentialAgent定义完整面试流程顺序
4. 创建并返回Runner实例

### 流程顺序
1. 分析简历（ResumeAnalysisAgent）
2. 循环执行：生成问题→获取回答→评估回答（loopAgent）
3. 生成面试报告（InterviewReportAgent）

## 7. InterView_Report.go

### 功能描述
InterviewReportAgent 是面试报告生成专家，负责分析用户的面试记录并生成详细的面试报告。

### 核心实现
```go
func NewInterviewReportAgent() adk.Agent {
    // 创建配置和模型
    a, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:        "InterviewReportAgent",
        Description: "一个可以解析面试记录生成面试报告的智能体",
        Instruction: `你是一名资深的面试报告专家，负责对用户的面试记录进行分析,并输出对应的面试报告。`,
        Model:       chat.CreatOpenAiChatModel(ctx),
        ToolsConfig: adk.ToolsConfig{
            ToolsNodeConfig: compose.ToolsNodeConfig{
                Tools: []tool.BaseTool{},
            },
        },
        MaxIterations: 10,
    })
    // 错误处理
    return a
}
```

### 工作方式
- 作为面试流程的最后一环，处理所有前面环节收集的信息
- 生成包含面试评估、优缺点分析和建议的完整报告
- 目前实现中没有配置特定工具，但预留了工具扩展接口

## Agent架构关系图

```
MainAgent
    └── RouterAgent
        ├── ResumAnalysisAgent
        ├── QuestionGeneratorAgent
        └── AnswerEvalAgent

InterView_Procrss
    └── SequentialAgent
        ├── ResumAnalysisAgent
        ├── LoopAgent
        │   ├── QuestionGeneratorAgent
        │   └── AnswerEvalAgent
        └── InterviewReportAgent
```

## 核心技术栈

- 基于 CloudWeGo Eino ADK (Agent Development Kit) 构建
- 采用 Chat Model Agent 作为基础架构
- 支持工具调用扩展功能
- 实现了循环和顺序组合Agent模式

## 使用说明

1. 系统提供两种入口模式：
   - 通过 `NewMainAgent()` 创建基于路由的灵活交互模式
   - 通过 `NewInterviewProcessAgent()` 创建完整的自动化面试流程

2. 每个Agent都可以独立使用，也可以根据需要组合使用

3. 所有Agent都支持配置最大迭代次数，防止无限循环

## 扩展建议

1. 为 InterviewReportAgent 添加专门的报告生成工具
2. 增强各个Agent的指令，使其更专业和详细
3. 考虑添加更多专业领域的子Agent，如特定技术栈专家
4. 实现更复杂的Agent协作机制，支持条件分支和错误恢复

## 总结

Agent模块是面试助手系统的智能核心，通过组合不同功能的专业Agent，实现了从简历分析到面试报告生成的完整面试流程自动化。系统采用了模块化设计，使得每个Agent职责明确，便于维护和扩展。