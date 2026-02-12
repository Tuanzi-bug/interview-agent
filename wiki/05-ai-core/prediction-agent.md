# 预测智能体 (Prediction Agent)

> 基于简历预测可能问到的面试题，帮助候选人备战。

## 🎯 职责

在正式面试前，候选人通常希望知道"针对我的简历，面试官可能会问什么"。预测智能体就是为了满足这一需求。

**核心产出**:
- 5 道高概率面试题
- 每道题的考察重点
- 推荐的回答思路 (Thinking Path)
- 参考答案
- 可能的追问 (Follow up)

## 🧠 核心实现

代码位置: `backend/chatApp/agent/prediction/prediction_agent.go`

`PredictionAgent` 也是一个标准 `ChatModelAgent`，但它是一个 **One-Shot Agent**，即通常只进行一次交互（"给你简历，给我题目"），不需要通过工具多轮获取信息，而是直接将简历内容作为 Prompt 的上下文传入。

*(注：视具体实现，简历内容可能通过 System Prompt 注入，或作为 User Message 传入)*

## 📝 Prompt 设计策略

1. **数量约束**: "必须严格生成 5 道题目"。
2. **内容约束**: "必须结合简历中的项目经历和技能点"。
3. **格式约束**: 标准 JSON，无 Markdown 标记。

### JSON 结构

```json
{
  "questions": [
    {
      "question": "请介绍一下你在XXX项目中是如何解决缓存一致性问题的？",
      "content": "【重点考察】分布式缓存一致性",
      "focus": "项目经历深度、解决方案的合理性",
      "thinking_path": "先说场景，再说方案（如延迟双删、Binlog订阅），最后说优缺点",
      "reference_answer": "...",
      "follow_up": "如果数据库主从延迟很大怎么办？"
    },
    ...
  ]
}
```

## 💡 应用场景

- **面试前突击**: 用户上传简历后，一键生成"押题"。
- **模拟演练**: 配合面试智能体，可以选择这 5 道题作为"开卷考试"的题目。

---
