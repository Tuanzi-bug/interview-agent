# 面试智能体 (Interview Agent)

> 模拟真实面试官，进行多轮技术或综合面试。

## 🎯 职责与分类

面试智能体是系统的核心，负责与用户进行直接的面试交互。为了应对不同的面试场景，我们设计了多种类型的面试官。

### 1. 综合面试官 (Comprehensive)

- **校招 (School)**: `SchoolComprehensiveAgent`
    - **目标**: 评估应届生的基础知识、学习潜力、职业素养。
    - **特点**: 题目涉及面广，包含基础算法、数据结构、计算机基础，以及软技能问题。
- **社招 (Social)**: `SocialComprehensiveAgent`
    - **目标**: 评估候选人的实战经验、架构设计能力、领导力。
    - **特点**: 深入挖掘项目经验，关注系统设计、故障排查、技术选型。

### 2. 专项面试官 (Specialized)

- **语言专项**: `GoAgent`, `JavaAgent`
    - **目标**: 深度考察特定编程语言的特性、原理及生态。
- **中间件专项**: `RedisAgent`, `MysqlAgent`, `MqAgent`
    - **目标**: 考察特定中间件的原理、调优及实战应用。

## 🧠 核心实现

代码位置: `backend/chatApp/agent/interview/`

所有面试智能体都基于 Eino 的 `ChatModelAgent` 实现。

### 构造函数示例 (Go)

```go
func NewGoSpecializedAgent(userId uint, needResumeTool bool) (adk.Agent, error) {
    // ... 模型初始化 ...
    
    // 如果需要，注入简历获取工具
    var toolsConfig adk.ToolsConfig
    if needResumeTool {
        toolsConfig = adk.ToolsConfig{
            ToolsNodeConfig: compose.ToolsNodeConfig{
                Tools: []componenttool.BaseTool{
                    tool2.GetResumeInfoTool(),
                },
            },
        }
    }

    baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
        Name:          "GoSpecializedAgent",
        Instruction:   GoSpecializedAgentInstruction, // 核心 Instruction
        Model:         model,
        ToolsConfig:   toolsConfig,
        MaxIterations: 15,
    })
    return baseAgent, nil
}
```

## 📝 Prompt 设计策略

Prompt 是智能体的灵魂。以下是`SchoolComprehensiveAgent`的核心 Prompt 策略：

1. **角色定义**: "你是一个经验丰富的校招综合面试官..."
2. **核心职责**:
    - "每次调用只生成一道问题"
    - "通过递进式的问题深入了解候选人的真实水平"
3. **面试策略**:
    - 第一题从候选人背景出发。
    - 结合实际场景和代码示例。
    - 难度循序渐进。
4. **输出约束**:
    - 必须返回 JSON: `{"question_text": "..."}`。

## 🔄 工作流程

1. **初始化**: 根据用户选择的面试类型（如 Go 专项），初始化对应的 Agent 实例。
2. **加载上下文**: 如果关联了简历，Agent 会调用 `get_resume_info` 获取简历摘要。
3. **首轮提问**: Agent 生成第一个问题。
4. **交互循环**:
    - 用户回答。
    - Agent 接收回答，结合上下文（Memory）进行分析。
    - Agent 生成下一个问题（追问或新话题）。
5. **结束**: 达到预定轮数或用户手动结束。

---
