# Backend Knowledge Base

**Parent:** [../AGENTS.md](../AGENTS.md)

## OVERVIEW
The `backend/` directory is the core of the `ai-eino-interview-agent` system, focusing on business logic, API orchestration, and AI agent implementation using the **Eino** framework.

- **Go Module**: `ai-eino-interview-agent`
- **Primary Frameworks**: Hertz (Web), Eino (AI), GORM (ORM).

## DIRECTORY STRUCTURE
- `api/`: Entry point for external communication. Contains **Hertz** handlers and router definitions.
- `chatApp/`: Core AI logic. This is where **Eino** agents, graphs, and tools are implemented.
- `internal/`: Private business logic. Includes services, repositories, and specialized integrations like Milvus for RAG.
- `idl/`: Interface Definition Language (Thrift/Proto) files. **DO NOT EDIT** these files manually; they are used for code generation.

## KEY CONCEPTS

For global conventions (Code Generation, Dependency Management, etc.), refer to the [Root AGENTS.md](../AGENTS.md).

### 1. Brain (`chatApp`) vs. Mouth (`api`)
To maintain a clean architecture, we distinguish between the AI's "thought" process and its "communication":
- **`chatApp/`**: Implements the AI orchestration (agents, chains, tools). This is the "Brain".
- **`api/`**: Handles HTTP requests/responses. This is the "Mouth".
Handlers in `api/` should delegate complex logic to services in `internal/` or agent services in `chatApp/`.

### 2. Internal Package Discipline
The `internal/` directory follows Go's visibility rules. Packages here are only accessible to other packages within the `backend/` module. 
- Use `internal/service` for high-level business flows.
- Use `internal/repository` for data persistence.
- Do not expose internal logic to the root or other modules (like `mcpserver`).
