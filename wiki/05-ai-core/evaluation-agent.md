# 评估智能体 (Evaluation Agent)

> 分析面试全过程，生成详细的多维度评估报告。

## 🎯 职责

面试结束后，评估智能体负责回顾所有的问答记录，这相当于真实面试中的"面试官填写面评"环节。

**核心产出**:
- 总体评价与建议
- 多维度打分（如：基础知识、实战能力、沟通表达等）
- 具体的优缺点分析

## 🛠️ 工具集成

`RecordEvaluationAgent` 使用 `get_mianshi_info` 工具。

- **输入**: `user_id`, `report_id`
- **功能**: 从数据库拉取该场面试的所有 QA 对话记录。
- **目的**: 确保 LLM 拥有完整的上下文，而不是仅依赖有限的 Token 窗口或不可靠的记忆。

## 🧠 核心实现

代码位置: `backend/chatApp/agent/record_evaluation/record_evaluation_agent.go`

提示词中定义了严格的评分标准：

- **90-100分**: 优秀
- **80-89分**: 良好
- **70-79分**: 中等
- ...

## 📝 Prompt 设计策略

1. **维度评分**: 要求 LLM 生成 4-8 个维度的评分。
2. **证据支撑**: 评分必须基于实际的对话记录。
3. **结构化输出**:

```json
{
  "comment": "总体评价内容...",
  "dimensions": [
    {
      "dimension_name": "基础知识",
      "evaluation": "候选人在Golang并发模型方面理解深刻...",
      "score": 85
    },
    ...
  ]
}
```

## 🔄 工作流程

1. 触发评估（通常由用户在前端点击"结束面试并生成报告"）。
2. Agent 调用 `get_mianshi_info` 获取完整 Transcript。
3. Agent 分析每一轮问答的质量。
4. Agent 综合所有维度，生成最终 JSON 报告。
5. 后端解析 JSON 并存入数据库，供前端展示雷达图和详细报告。

---
