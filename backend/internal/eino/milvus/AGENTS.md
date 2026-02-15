# Milvus RAG Knowledge Base

**Parent:** [Backend Knowledge Base](../../AGENTS.md)  
**See also:** [Root Knowledge Base](../../../../AGENTS.md)

## OVERVIEW
This directory contains the RAG (Retrieval-Augmented Generation) logic, specifically the integration with the **Milvus** vector database. It handles document indexing, vector storage, and similarity retrieval.

---

## RAG PIPELINE
The system follows a standard RAG workflow:
1. **Document**: Markdown files are loaded from `data/` or external sources (e.g., Feishu).
2. **Split**: Documents are partitioned into smaller chunks using `splitter/`.
3. **Embed**: Text chunks are converted into vector embeddings via `storage/embedding.go`.
4. **Vector**: Embeddings and metadata are stored/indexed in Milvus using `storage/indexer.go`.
5. **Search**: High-relevance chunks are retrieved based on query similarity via `retrieval/`.

---

## KEY COMPONENTS
- **`MilvusManager` (`init.go`)**: Entry point for initializing connections and services.
- **`importer.go`**: High-level service to import documents into the pipeline.
- **`retrieval/`**: Logic for searching the vector store with optional filters (Language, Category).
- **`storage/`**: Handles the underlying vectorization and indexing.

---

## MANAGEMENT & CLI
For administrative tasks (creating collections, checking status), use the **milvusctl** tool:
- **Source**: `cmd/milvusctl/main.go`
- **Usage**: Reference the documentation in `cmd/milvusctl/cmd.md`.

---

## CONVENTIONS
- Always verify Milvus connection health before starting indexing/retrieval operations.
- Use `context.Context` for all external database and API calls.
- Follow the metadata schema defined in `document_metadata.go` for consistent retrieval.
