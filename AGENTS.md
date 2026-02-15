# AGENTS.md

## OVERVIEW
`go-eino-interview-agent` is an AI-powered mock interview system designed to help job seekers prepare for technical interviews. The project leverages ByteDance's **Hertz** (HTTP framework) and **Eino** (AI Orchestration) to provide intelligent resume analysis, question generation, and answer evaluation.

### Tech Stack
- **Backend**: Go 1.25+, Hertz (Web), Eino (LLM Framework), GORM (ORM)
- **Frontend**: Next.js 14+ (TypeScript, Tailwind CSS)
- **Databases**: MySQL (Metadata), Redis (Cache/MQ), Milvus (Vector DB for RAG)
- **Protocols**: RESTful API, MCP (Model Context Protocol)

### Structure
The project follows a **split monorepo** structure:
- Root directory contains orchestration and common docs.
- `backend/`: Core logic, AI agents, and API service.
- `frontend/`: User interface.
- `mcpserver/`: Implementation of the Model Context Protocol.

---

## STRUCTURE
- `backend/`: Core logic, AI agents, and API service. ([Backend Guide](backend/AGENTS.md))
    - `api/`: Route definitions and HTTP handlers.
    - `chatApp/`: AI agent implementations. ([Eino Guide](backend/chatApp/AGENTS.md))
    - `internal/`: Business logic and services.
        - `eino/milvus/`: RAG logic. ([Milvus Guide](backend/internal/eino/milvus/AGENTS.md))
- `frontend/`: Next.js source code. ([Frontend Guide](frontend/AGENTS.md))
- `mcpserver/`: MCP server implementation.
- `doc/` & `wiki/`: Technical designs and guides.

---

## KNOWLEDGE BASES
For detailed information about specific parts of the system, refer to:
- [Backend Knowledge Base](backend/AGENTS.md)
- [Frontend Knowledge Base](frontend/AGENTS.md)
- [AI Orchestration (Eino) Guide](backend/chatApp/AGENTS.md)
- [Milvus RAG Guide](backend/internal/eino/milvus/AGENTS.md)

---

## WHERE TO LOOK
- **AI Agent Definitions**: `backend/chatApp/agent/` - Where Eino graphs and chains are defined.
- **RAG & Vector Logic**: `backend/internal/eino/milvus/` - Milvus integration for knowledge retrieval.
- **Business Logic**: `backend/internal/service/` - High-level interview flows.
- **API Handlers**: `backend/api/handler/` - HTTP request entry points.
- **Frontend Components**: `frontend/src/components/` - UI elements.

---

## CONVENTIONS
### 1. Dual `go.mod` Management
The project has multiple `go.mod` files (root, `backend/`, `mcpserver/`). 
- When working in the backend, use `backend/` as the working directory.
- Ensure dependencies are synced across modules if sharing logic.

### 2. Code Generation
- **Do not manually edit** files in `backend/idl/` or any files explicitly marked as generated (e.g., by GORM Gen or Eino plugins).
- Regenerate code using appropriate commands (found in `常用命令.md`) after modifying IDLs or schemas.

### 3. AI Patterns
- Follow the Eino orchestration patterns: use Graphs for complex flows and Chains for simple sequences.
- Always use `context.Context` for cancellation and timeouts in AI calls.

### 4. Naming & Style
- Follow standard Go idiomatic naming (CamelCase for exports).
- Keep the "Happy Path" left-aligned.
