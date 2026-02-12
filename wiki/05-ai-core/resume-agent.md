# 简历分析智能体 (Resume Agent)

> 负责解析 PDF 简历，提取结构化信息，为面试提供上下文。

## 🎯 职责

简历分析是智能面试的第一步。该智能体将非结构化的 PDF 文档转化为结构化的 JSON 数据，以便面试官智能体和预测智能体使用。

**核心产出**:
- 个人基本信息
- 教育背景
- 工作/实习经历
- 项目经验
- 技术栈摘要
- 核心竞争力分析

## 🛠️ 工具集成

`ResumeParserAgent` 必须配合 **PDF 解析工具** (`pdf_to_text`) 使用。

1. **PDF 解析**: 将二进制 PDF 转换为纯文本。
2. **文本提取**: LLM 从海量文本中提取关键字段。

## 🧠 核心实现

代码位置: `backend/chatApp/agent/resume/resume.go`

```go
baseAgent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "ResumeParserAgent",
    Instruction: `你是一个专业的简历分析专家...
    任务步骤（必须按顺序执行）：
    1. 【必须】使用 pdf_to_text 工具解析提供的简历文件路径...
    2. 从解析的简历文本中提取...
    ...
    必须返回的JSON格式...`,
    Model:       model,
    // ...
})
```

## 📝 Prompt 设计策略

1. **CoT (Chain of Thought) 引导**: 明确指示 "步骤1: 调用工具", "步骤2: 提取信息"。
2. **强制 JSON 输出**: 提供了详细的 JSON 模板，要求 "所有字段都必须填充实际数据"。
3. **分析增强**: 不仅仅是复制文本，还要求 LLM 进行归纳总结，如 "分析候选人的背景特点" 和 "生成面试建议"。

### 输出 JSON 结构示例

```json
{
  "basic_info": { "name": "张三", ... },
  "education": [ ... ],
  "work_experience": [ ... ],
  "tech_stack": ["Go", "Kubernetes"],
  "projects": [ ... ],
  "interview_focus_areas": ["微服务架构", "高并发处理"]
}
```

---
