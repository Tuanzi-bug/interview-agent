# 向量数据库 (Vector Database)

> 基于 Milvus 构建的文档知识库，为智能体提供特定领域的知识增强 (RAG)。

## 📚 为什么需要向量数据库？

通用大模型（如 GPT-4）拥有海量的通用知识，但缺乏：
1. **私有领域知识**: 如公司内部文档、特定技术栈的最新细节。
2. **实时信息**: 如最新的面试题库。

通过 RAG (Retrieval-Augmented Generation) 技术，我们将知识转化为向量存储，在生成回答前先检索相关内容，注入到 Prompt 中，从而让 Agent "懂行"。

## 🏗️ 架构概览

本项目封装了 `milvus` 模块，位于 `backend/internal/eino/milvus/`。

### 核心组件

1. **Embedding Service**: 使用 OpenAI Embedding 模型将文本转化为向量。
2. **Document Splitter**: 将长文档（如 Markdown）智能切分为适合检索的片段 (Chunk)。
3. **Milvus Manager**: 负责连接 Milvus 实例、管理 Collection 和索引。
4. **Retriever**: 供 Agent 调用的检索器接口。

```
[Markdown Docs] -> [Splitter] -> [Chunks] -> [Embedding] -> [Vector] -> [Milvus]
                                                                          ^
                                                                          |
[Agent Query] -------------------------------------------------> [Retriever]
```

## 🛠️ 实现细节

### 文档处理流水线

1. **导入 (Import)**: `importer.go` 读取 Markdown 文件。
2. **切分 (Split)**: `markdown.go` 根据标题层级（# ## ###）进行语义切分，保持上下文完整性。
3. **入库 (Index)**: 生成向量并存入 Milvus。

### 检索策略

- **相似度搜索**: 基于 Cosine Similarity 寻找最匹配的知识片段。
- **阈值过滤**: 过滤掉相关度过低的结果。
- **Metadata 过滤**: 支持按标签、来源等过滤。

## 💡 使用示例

在 Eino 中使用 Milvus 检索器：

```go
// 创建检索器
retriever := NewMilvusRetriever(ctx, &MilvusRetrieverConfig{
    Collection: "interview_knowledge",
    TopK:       3,
})

// 作为工具提供给 Agent
tools_config := adk.ToolsConfig{
    Tools: []componenttool.BaseTool{
        milvus_retriever_tool, // 封装后的 Tool
    },
}
```

这样，当用户问到某个特定知识点（如"字节跳动企业文化"）时，Agent 可以自动调用检索工具查阅资料。

---
