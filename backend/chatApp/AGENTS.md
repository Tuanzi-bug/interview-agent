# AI Orchestration Knowledge Base (Eino)

**Parent:** [Backend Knowledge Base](../AGENTS.md)  
**See also:** [Root Knowledge Base](../../AGENTS.md)

## OVERVIEW
The `chatApp/` directory is the "brain" of the application. It leverages ByteDance's **Eino** framework to define how LLMs are prompted, how they use tools, and how they transition between states in an interview flow.

## CORE CONCEPTS

### 1. Eino Graphs vs. Chains
- **Graphs (`compose.Graph`)**: Used for complex, stateful workflows that require branching or iterative logic (e.g., the interview flow which evolves based on user answers).
- **Chains (`compose.Chain`)**: Used for simple, linear sequences of LLM calls (e.g., direct question generation).

### 2. Agents (`adk.Agent`)
Most AI logic is encapsulated in `adk.Agent` implementations. These combine:
- **Instruction**: The system prompt.
- **Tools**: Capabilities the agent can invoke (defined in `tool/`).
- **Model**: The LLM configuration (defined in `chat/`).

### 3. Agent Definition vs. Execution
We maintain a strict separation between *what* an agent is and *how* it is used:
- **`agent/` (Definition)**: Where the agent's prompts, tools, and Eino graphs are defined.
- **`agent_service/` (Execution)**: Service wrappers that manage the lifecycle of an agent session, handle data persistence, and provide a clean API for the rest of the backend.

## DIRECTORY STRUCTURE

| Path | Purpose | Key Files |
|------|---------|-----------|
| `agent/` | AI Agent definitions | `resume/`, `interview/`, `record_evaluation/` |
| `agent_service/` | Execution and lifecycle management | `interview_agent_service.go`, `record_evaluation_service.go` |
| `tool/` | LLM Tool implementations | `pdfParserTool.go`, `milvus_retriever_tool.go` |
| `chat/` | LLM client initialization | `openAi.go` |
| `config/` | Chat-specific configurations | `config.go` |

## KEY PATTERNS

### Adding a New Agent
1. Define the agent struct and its Eino graph/chain in `agent/{purpose}/`.
2. Implement any required tools in `tool/`.
3. Create a corresponding service in `agent_service/` to handle invocation and business logic.

### Invocation Flow
`api/handler` -> `internal/service` -> `agent_service` -> `adk.Runner` -> `adk.Agent` (Eino Graph)
