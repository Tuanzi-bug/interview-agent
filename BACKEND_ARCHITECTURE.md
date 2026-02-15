# Backend Technology Stack & Architecture Diagram

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│                    (Next.js Frontend)                           │
└─────────────────────────────────┬───────────────────────────────┘
                                  │ HTTP/REST
                                  ▼
┌─────────────────────────────────────────────────────────────────┐
│                      HERTZ WEB FRAMEWORK                        │
│  (HTTP Server, Routing, Middleware, Connection Handling)       │
├─────────────────────────────────────────────────────────────────┤
│  Middleware Stack:                                              │
│  1. Recovery (Panic catching)                                  │
│  2. CORS Handler                                               │
│  3. JWT Authentication (with whitelist)                        │
│  4. Route Handler                                              │
└─────────────────────────┬───────────────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
   ┌─────────┐      ┌──────────┐     ┌──────────────┐
   │API      │      │Interview │    │Resume/User   │
   │Handlers │      │Handlers  │    │Handlers      │
   └────┬────┘      └────┬─────┘    └──────┬───────┘
        │                │                │
        └────────────────┼────────────────┘
                         ▼
        ┌────────────────────────────────────┐
        │   SERVICE LAYER                    │
        │ (Business Logic)                   │
        │                                    │
        │ ├─ Interview Service               │
        │ ├─ User Service                    │
        │ ├─ Resume Service                  │
        │ ├─ Prediction Service              │
        │ └─ Common Services                 │
        └────────┬──────────────────────────┘
                 │
        ┌────────┴────────────────────────┐
        ▼                                  ▼
┌─────────────────┐            ┌──────────────────────┐
│EINO AI AGENTS   │            │REPOSITORY LAYER      │
│                 │            │(Data Access)         │
│ ├─ Resume       │            │                      │
│ │  Parser       │            │ ├─ GORM + MySQL      │
│ ├─ Interview    │            │ ├─ Redis Client      │
│ │  (School)     │            │ └─ Milvus Client     │
│ ├─ Interview    │            │                      │
│ │  (Social)     │            │ Models (DAOs):       │
│ ├─ Interview    │            │ ├─ User              │
│ │  (Specialized)│            │ ├─ Interview         │
│ ├─ Prediction   │            │ ├─ Resume            │
│ ├─ Evaluation   │            │ └─ ...               │
│ │               │            │                      │
│ └─ Tools:       │            │                      │
│   ├─ PDF Parser │            │                      │
│   ├─ Retriever  │            │                      │
│   ├─ Get Info   │            │                      │
│   └─ Chat Model │            │                      │
└────────┬────────┘            └──────┬───────────────┘
         │                            │
         │         ┌──────────────────┴─────────────────┬─────────┐
         │         ▼                                     ▼         ▼
         │    ┌─────────────┐  ┌─────────────┐  ┌──────────────────┐
         │    │  MESSAGE    │  │   CONFIG    │  │  INITIALIZATION  │
         │    │   QUEUE     │  │ MANAGEMENT  │  │   & UTILITIES    │
         │    │ (Redis-based)│  │  (YAML)    │  │                  │
         │    └─────────────┘  └─────────────┘  └──────────────────┘
         │         │
         └─────────┼─────────────────────────────────┐
                   │                                 ▼
                   ▼                        ┌─────────────────┐
         ┌──────────────────┐               │  CONSUMER       │
         │  ASYNC TASKS     │               │ GOROUTINES      │
         │ (Async eval)     │               │                 │
         └──────────────────┘               │ (AI processing) │
                                            └─────────────────┘
         
         └─────────────────────────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        ▼                ▼                ▼
   ┌──────────┐    ┌──────────┐    ┌──────────────┐
   │ MYSQL    │    │ REDIS    │    │ MILVUS       │
   │          │    │          │    │ (Optional)   │
   │ Tables:  │    │ Cache    │    │              │
   │ ├─User   │    │ Sessions │    │ Vector DB    │
   │ ├─Resume │    │ Queue    │    │ (Embeddings) │
   │ ├─Interview         │    │ Semantic     │
   │ ├─Dialogue         │    │ Search       │
   │ ├─Evaluation       │    │              │
   │ └─...              │    └──────────────┘
   └──────────┘    └──────────┘
```

---

## 📦 Technology Stack Breakdown

### Backend Framework
```
Hertz (ByteDance Web Framework)
├─ HTTP/1.1 & HTTP/2 support
├─ High performance async I/O
├─ Built-in middleware system
└─ Tight Eino integration
    Version: v0.10.3
```

### Database & ORM
```
GORM (Object-Relational Mapping)
├─ MySQL Driver
├─ Connection pooling
├─ Auto-migration
└─ Type-safe queries
    GORM: v1.25.11
    MySQL Driver: v1.5.7
```

### Caching & Message Queue
```
Redis
├─ In-memory caching
├─ Session storage
├─ Message queue (async tasks)
└─ Go client (go-redis)
    Version: v9.16.0
```

### Vector Database (Optional)
```
Milvus
├─ Vector similarity search
├─ Document embeddings
├─ Semantic search
└─ Filtering by metadata
    Version: v2.4.2
    Status: Currently disabled in main.go
```

### AI/LLM Framework
```
Eino (ByteDance AI Application Framework)
├─ Agent system (adk.Agent)
├─ Tool composition
├─ LLM orchestration
└─ Prompt management
    Core: v0.5.11
    
Extensions:
├─ PDF Parser: document/parser/pdf
├─ Doc Splitter: document/transformer/splitter/recursive
├─ Embedding: components/embedding/ark
├─ Indexer: components/indexer/milvus
├─ Retriever: components/retriever/milvus
├─ LLM Models: components/model/{ark,openai}
└─ Tools: components/tool/googlesearch
```

### LLM Models
```
1. OpenAI (Primary)
   └─ Version: gpt-4 (configurable)
   
2. ARK (ByteDance LLM)
   ├─ Embedding service
   └─ Alternative LLM provider
   
3. Coze Studio Integration
   └─ AI platform integration
```

### Authentication & Security
```
JWT (JSON Web Tokens)
├─ Stateless authentication
├─ Token validation middleware
├─ Whitelist-based auth skipping
└─ Version: v5.2.2

WeChat OAuth
├─ Social login
├─ OAuth integration
└─ LarkSuite SDK: v3.4.26
```

### Additional Libraries
```
├─ Sonic: Fast JSON serialization (v1.14.2)
├─ UUID: Unique ID generation (v1.6.0)
├─ godotenv: .env file loading (v1.5.1)
├─ YAML: Configuration parsing (v3.0.1)
├─ Protobuf: Data serialization (v1.36.10)
├─ OpenTelemetry: Observability (v1.36.0)
└─ MCP: Model Context Protocol (v0.43.0)
```

---

## 🎯 AI Agent Architecture

### Agent Types & Specialization

```
┌─────────────────────────────────────────────────────┐
│         EINO MULTI-AGENT SYSTEM                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─ INTERVIEW AGENTS (Interview Selection)          │
│  │                                                  │
│  ├─ Comprehensive Interviews                        │
│  │  ├─ SchoolComprehensiveAgent                     │
│  │  │  └─ For campus/grad recruitment               │
│  │  └─ SocialComprehensiveAgent                     │
│  │     └─ For professional recruitment              │
│  │                                                  │
│  ├─ Specialized Interviews (Technology-specific)    │
│  │  ├─ GoAgent (Golang specialist)                  │
│  │  ├─ JavaAgent                                    │
│  │  ├─ MySQLAgent                                   │
│  │  ├─ RedisAgent                                   │
│  │  └─ MQAgent (Message Queue specialist)           │
│  │                                                  │
│  ├─ RESUME AGENT                                    │
│  │  └─ ResumeParserAgent                            │
│  │     ├─ PDF parsing                               │
│  │     └─ Resume analysis & extraction              │
│  │                                                  │
│  ├─ PREDICTION AGENT                                │
│  │  └─ PredictionAgent                              │
│  │     └─ Resume-based prediction                   │
│  │                                                  │
│  └─ EVALUATION AGENT                                │
│     └─ RecordEvaluationAgent                        │
│        ├─ Answer evaluation                         │
│        └─ Scoring & feedback                        │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Agent Composition Pattern

```
Each Agent = Model + Tools + Instructions + Config

┌────────────────────────────┐
│ ChatModelAgent             │
├────────────────────────────┤
│ Name: "ResumeParserAgent"  │
│ Description: Resume parser │
│ Instruction: System prompt │
│ Model: OpenAI ChatModel    │
│ Tools:                     │
│  ├─ pdf_to_text            │
│  ├─ milvus_retriever       │
│  └─ get_resume_info        │
│ MaxIterations: 15          │
└────────────────────────────┘
      │
      ├─ Tool Loop ─────┐
      │                 ▼
      │            Tool Execution
      │                 │
      │                 ▼
      │            Tool Result
      │                 │
      └─────────────────┘
      
      Output: Structured JSON/Response
```

### Tool System

```
┌──────────────────────────────────────────┐
│ EINO TOOL COMPOSITION                    │
├──────────────────────────────────────────┤
│                                          │
│ 1. PDFParserTool                         │
│    ├─ Input: Resume file path            │
│    ├─ Process: Extract text from PDF     │
│    └─ Output: Resume text content        │
│                                          │
│ 2. MilvusRetrieverTool                   │
│    ├─ Input: Query string                │
│    ├─ Process: Vector similarity search  │
│    └─ Output: Top-K relevant documents   │
│                                          │
│ 3. GetResumeInfoTool                     │
│    ├─ Input: User ID, Resume ID         │
│    ├─ Process: Database lookup           │
│    └─ Output: Resume data                │
│                                          │
│ 4. GetMianshiInfoTool                    │
│    ├─ Input: Interview ID                │
│    ├─ Process: Interview context lookup  │
│    └─ Output: Interview history          │
│                                          │
└──────────────────────────────────────────┘
```

---

## 📊 Request Processing Flow

### Interview Answer Submission Flow

```
1. USER SUBMITS ANSWER
   │
   └──→ POST /api/v1/interview/{id}/answer
       │
       ▼
2. HERTZ ROUTER
   ├─ Route matching
   └─ Middleware chain
       │
       ▼
3. JWT MIDDLEWARE
   ├─ Token validation
   └─ Extract user context
       │
       ▼
4. INTERVIEW HANDLER
   ├─ Parse request body
   ├─ Validate answer
   └─ Call service
       │
       ▼
5. INTERVIEW SERVICE
   ├─ Fetch interview record
   ├─ Validate state
   ├─ Store answer → InterviewDialogue table
   └─ Enqueue async task
       │
       ▼
6. REDIS MESSAGE QUEUE
   ├─ Publish task message
   └─ Return immediately (fast response)
       │
       ▼
7. HTTP RESPONSE
   ├─ Status: 200 OK
   └─ Body: {"answer_id": 123, "status": "submitted"}
   
   [Parallel Processing in Background]
   │
   ▼
8. CONSUMER GOROUTINE (Async)
   ├─ Dequeue task
   ├─ Fetch interview context
   └─ Call RecordEvaluationAgent
       │
       ├─ Get ChatModel (OpenAI)
       ├─ Get Tools (MianshiInfoTool)
       ├─ Execute agent loop
       │  ├─ Agent calls tool if needed
       │  ├─ Tool returns data
       │  └─ Agent generates evaluation
       └─ Format result as JSON
       │
       ▼
9. STORE EVALUATION
   └─ InterviewEvaluation table
       │
       ├─ scores
       ├─ feedback
       ├─ suggestions
       └─ timestamp
       │
       ▼
10. RESULT AVAILABLE
    └─ User polls/gets evaluation via
       GET /api/v1/interview/{id}/evaluation
```

---

## 🔄 Data Flow

### Interview Lifecycle

```
USER REGISTRATION
    │
    ▼
[User Table]
    │
    ├─→ Resume Upload
    │   │
    │   ▼
    │ [Resume Table]
    │   │
    │   └─→ Resume Analysis
    │       (ResumeParserAgent)
    │       │
    │       ▼
    │   [Resume Embeddings in Milvus] (optional)
    │
    ├─→ Interview Creation
    │   │
    │   ▼
    │ [InterviewRecord Table]
    │   │
    │   └─→ Start Interview
    │       │
    │       ▼
    │   Questions + Answers
    │   (Real-time conversation)
    │       │
    │       ├─ Question generated by Interview Agent
    │       │  (SchoolComprehensiveAgent, GoAgent, etc.)
    │       │
    │       └─ Answer from User
    │           │
    │           ▼
    │       [InterviewDialogue Table]
    │           │
    │           └─→ Async Evaluation
    │               │
    │               ▼
    │           RecordEvaluationAgent
    │           (with MianshiInfoTool)
    │               │
    │               ▼
    │           [InterviewEvaluation Table]
    │
    └─→ Interview Result
        │
        ├─ Prediction (PredictionAgent)
        │ (ResumeAgent + InterviewResults)
        │
        └─ [PredictionRecord Table]
```

---

## 🔐 Security Architecture

```
┌────────────────────────────────────────┐
│ SECURITY LAYERS                        │
├────────────────────────────────────────┤
│                                        │
│ 1. HTTPS/TLS (Transport)               │
│    └─ CORS headers configured          │
│                                        │
│ 2. JWT AUTHENTICATION                  │
│    ├─ Token validation                 │
│    ├─ Signature verification           │
│    └─ Whitelist for public endpoints   │
│                                        │
│ 3. INPUT VALIDATION                    │
│    ├─ DTO validation (struct tags)     │
│    ├─ Database constraints             │
│    └─ Type safety (Go generics)        │
│                                        │
│ 4. DATA PROTECTION                     │
│    ├─ Password hashing (bcrypt likely) │
│    ├─ Environment variables for secrets│
│    └─ No hardcoded credentials         │
│                                        │
│ 5. ERROR HANDLING                      │
│    ├─ Recovery middleware              │
│    ├─ No sensitive data in errors      │
│    └─ Structured error responses       │
│                                        │
└────────────────────────────────────────┘
```

---

## 📈 Scalability Considerations

### Horizontal Scaling
```
Load Balancer
    │
    ├─ Backend Instance 1 ─┐
    ├─ Backend Instance 2  ├─ Shared Databases
    └─ Backend Instance N ─┘
         │       │       │
         └───────┼───────┘
                 │
        ┌────────┼────────┐
        ▼        ▼        ▼
      MySQL    Redis   Milvus
```

### Message Queue Scaling
```
Multiple Consumer Goroutines
    │
    ├─ Consumer 1
    ├─ Consumer 2
    ├─ Consumer 3
    └─ Consumer N
         │
         └─→ All process tasks from Redis Queue
             (Parallel AI evaluation)
```

### Caching Strategy
```
User Session → Redis Cache (with TTL)
Resume Data → Redis Cache (on access)
Frequent Queries → Prepared statements + DB connection pooling
```

---

## 🚀 Deployment Architecture

```
┌────────────────────────────────────────────────┐
│ DOCKER COMPOSE (Development/Production)        │
├────────────────────────────────────────────────┤
│                                                │
│ Services:                                      │
│ ├─ backend (Go service) port 8080             │
│ ├─ mysql port 3306                            │
│ ├─ redis port 6379                            │
│ ├─ milvus port 19530 (optional)               │
│ ├─ frontend (Next.js) port 3000               │
│ └─ nginx (reverse proxy) port 80/443          │
│                                                │
│ Volumes:                                       │
│ ├─ mysql_data (persistent)                    │
│ └─ milvus_data (persistent)                   │
│                                                │
│ Networks:                                      │
│ └─ Shared network for service communication   │
│                                                │
└────────────────────────────────────────────────┘
```

---

## 📋 Configuration Management

```
config.yaml (Primary)
    ├─ server (host, port)
    ├─ database (MySQL DSN, pooling)
    ├─ redis (host, port, auth)
    ├─ milvus (address, credentials)
    ├─ openai (api_key, model)
    ├─ embedding (model, batch_size)
    ├─ hertz (timeouts, log_level)
    └─ (environment variables interpolation)
        │
        └─ Expanded at runtime
            ├─ ${OPENAI_API_KEY}
            ├─ ${DB_PASSWORD}
            └─ ${MILVUS_PASSWORD}
```

---

## Summary Table

| Layer | Technologies | Purpose |
|-------|-------------|---------|
| **Transport** | Hertz, HTTP/2 | Web server, routing |
| **Authentication** | JWT, WeChat OAuth | User auth, session management |
| **API** | REST, GORM entities | Request/response handling |
| **Business Logic** | Services, Eino Agents | Interview, resume, evaluation |
| **AI** | Eino, OpenAI/ARK, Tools | LLM orchestration, agent loop |
| **Data Access** | GORM, Repository Pattern | Database abstraction |
| **Persistence** | MySQL, Redis, Milvus | Data storage, caching, vectors |
| **Async** | Redis Queue, Goroutines | Background task processing |
| **Configuration** | YAML, Env Vars | Runtime configuration |
| **Deployment** | Docker Compose | Containerized services |

---

**Generated**: February 15, 2026  
**For detailed analysis, see**: `BACKEND_ANALYSIS.md`  
**For quick reference, see**: `BACKEND_SUMMARY.md`
