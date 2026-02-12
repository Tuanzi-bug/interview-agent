# Eino 框架理解

> Eino 是字节跳动开源的大语言模型应用开发框架，旨在简化 LLM 应用的构建过程。

## 🔍 框架概览

Eino 提供了从模型调用、提示词工程、工具调用到记忆管理的完整解决方案。在本项目中，我们使用 Eino 来构建所有的智能体。

### 核心组件

1. **Model (模型)**: 封装了底层 LLM 的调用，支持多种模型提供商（如 OpenAI, Claude 等）。
2. **Prompt (提示词)**: 管理和模板化提示词，支持动态参数注入。
3. **Tools (工具)**: 允许 LLM 调用外部函数或 API，扩展模型能力。
4. **Memory (记忆)**: 管理对话上下文，支持多轮对话。
5. **Agent (智能体)**: 将上述组件组合，实现复杂的业务逻辑。

## 💡 在项目中的应用

本项目主要使用了 Eino 的 `ChatModelAgent` 来构建各种角色的面试官和辅助智能体。

### 主要使用模式

```go
// 1. 创建 ChatModel
model, err := chat.CreatOpenAiChatModel(ctx, userId)

// 2. 配置 Agent
config := &adk.ChatModelAgentConfig{
    Name:        "AgentName",
    Description: "Agent Description",
    Instruction: "System Prompt",
    Model:       model,
    ToolsConfig: toolsConfig, // 可选：工具配置
}

// 3. 创建 Agent 实例
agent, err := adk.NewChatModelAgent(ctx, config)
```

### 关键特性应用

- **工具调用 (Function Calling)**:
    - 简历解析 (`pdf_to_text`)
    - 数据库检索 (`GetMianshiInfoTool`)
    - 简历信息获取 (`GetResumeInfoTool`)

- **流式输出 (Streaming)**:
    - 支持 SSE (Server-Sent Events) 输出，提供更好的用户体验。

- **多模态支持**:
    - 支持文本和文件输入（如简历 PDF）。

## 📚 学习资源

- [Eino 官方文档](https://github.com/cloudwego/eino) (注：需替换为实际链接)
- [CloudWeGo 官网](https://www.cloudwego.io/)

---
