# Backend Analysis Summary

## Quick Reference

### Core Stack
| Component | Technology | Version |
|-----------|-----------|---------|
| **Web Framework** | Hertz | v0.10.3 |
| **ORM** | GORM | v1.25.11 |
| **Database** | MySQL | via gorm.io/driver/mysql v1.5.7 |
| **Cache/Queue** | Redis | v9.16.0 |
| **Vector DB** | Milvus | v2.4.2 (optional, disabled) |
| **LLM Framework** | Eino | v0.5.11 |
| **LLM Models** | OpenAI + ARK | via eino-ext |
| **Auth** | JWT + WeChat OAuth | jwt/v5 v5.2.2 |

### Architecture Pattern
**Layered Architecture (Clean Hybrid)**
```
API Layer (Hertz handlers)
  ↓
Service Layer (business logic, AI agents)
  ↓
Repository Layer (data access)
  ↓
Data Layer (MySQL, Redis, Milvus)
```

### Key Integration Points

**1. Eino AI Framework**
- Multi-agent system for different interview types
- Tool composition for complex tasks
- LLM orchestration with OpenAI/ARK models

**2. Interview Agent Types**
- ResumeParserAgent (PDF parsing, resume analysis)
- ComprehensiveInterviewAgents (campus/professional recruitment)
- SpecializedAgents (Go, Java, MySQL, Redis, MQ)
- PredictionAgent (resume-based prediction)
- RecordEvaluationAgent (answer evaluation)

**3. Async Processing**
- Redis-based message queue
- Consumer goroutines for AI evaluation
- Decoupled HTTP responses from AI processing

**4. Document Vector Retrieval (Optional)**
- Milvus for semantic search
- Document embedding via ARK
- Filtering by language/category

### Directory Structure Highlights

```
backend/
├── api/                      # HTTP layer (Hertz routers, handlers, DTOs)
├── chatApp/                  # AI agents & LLM integration
│   ├── agent/               # Multi-agent definitions
│   ├── agent_service/       # Service layer for agents
│   ├── tool/                # Tool definitions
│   └── chat/                # LLM client wrappers
├── internal/                # Core application logic
│   ├── service/             # Business logic services
│   ├── repository/          # Data access layer (GORM, Redis)
│   ├── model/               # Domain models (GORM entities)
│   ├── middleware/          # JWT, CORS, etc.
│   ├── config/              # Configuration management
│   ├── eino/                # Eino framework integration
│   │   └── milvus/          # Vector database module
│   ├── mq/                  # Message queue
│   └── utils/               # Helpers
├── main.go                  # Entry point (initialization, server startup)
├── config.yaml              # Configuration file
└── go.mod/go.sum            # Dependencies
```

### Design Patterns in Use

1. **Repository Pattern** - Data access abstraction
2. **Dependency Injection** - Constructor-based DI for agents
3. **Singleton Pattern** - Global DB, Redis, Milvus instances
4. **Factory Pattern** - Agent creation functions
5. **Middleware Pipeline** - Request processing chain
6. **DAO Pattern** - Database access objects in models
7. **Tool Composition** - Eino's tool system for agent capabilities

### Initialization Flow (main.go)

1. Load `.env` file
2. Parse `config.yaml`
3. Expand environment variables
4. Connect to MySQL (GORM)
5. Connect to Redis
6. (Optional) Connect to Milvus
7. Initialize message queue consumer
8. Start Hertz HTTP server
9. Wait for shutdown signal (graceful cleanup)

### Critical Files

| File | Purpose |
|------|---------|
| `main.go` | Server initialization, middleware setup |
| `internal/repository/database.go` | MySQL connection, GORM setup |
| `internal/repository/redis.go` | Redis client management |
| `chatApp/agent/*` | Agent definitions (resume, interview, evaluation) |
| `chatApp/tool/*` | Agent tools (PDF parsing, retrieval, data fetching) |
| `internal/service/` | Business logic (user, interview, prediction) |
| `internal/eino/milvus/` | Vector DB integration (currently disabled) |
| `api/handler/` | HTTP request handlers |
| `api/router/` | Route registration |

### API Endpoints Pattern
```
/api/v1/user/          - User management (login, register, profile)
/api/v1/interview/     - Interview management
/api/v1/resume/        - Resume management
/api/v1/prediction/    - Prediction services
/health                - Health check
```

### Configuration Management

**Format**: YAML + Environment Variables
**Key Sections**:
- Database connection (MySQL)
- Redis configuration
- Hertz settings (timeouts, ports)
- OpenAI API keys
- Milvus configuration (if enabled)
- WeChat OAuth configuration

**Security**: Secrets via `${VAR_NAME}` placeholders expanded at runtime

### Strengths

✅ Clear separation of concerns (layered architecture)  
✅ Professional Go idioms and patterns  
✅ Flexible multi-agent system for extensibility  
✅ Async processing with message queue  
✅ Configuration externalization (no hardcoded secrets)  
✅ Graceful initialization and shutdown  
✅ Comprehensive test coverage in core modules  

### Areas for Improvement

⚠️ No structured logging (using `log` package)  
⚠️ No E2E tests for HTTP endpoints  
⚠️ No observability/monitoring infrastructure  
⚠️ No API rate limiting  
⚠️ Auto-migration vs explicit database migrations  
⚠️ CORS allows all origins (may need restriction)  

### Deployment

- **Docker Support**: Includes docker-compose.yml
- **Services**: Backend (Go), MySQL, Redis, Milvus, Frontend (Next.js), Nginx
- **Environment-Aware**: Separate config for dev/test/prod
- **Graceful Shutdown**: Proper cleanup of resources

---

## For More Details

See `BACKEND_ANALYSIS.md` for:
- Full technology stack breakdown
- Detailed architecture pattern explanation
- Eino integration deep dive
- Design patterns analysis
- Security considerations
- Request flow examples
- Configuration management details
- Extensibility roadmap
- Code quality assessment
