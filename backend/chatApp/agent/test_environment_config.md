# Go-Eino Interview Agent Test Environment Configuration

## 1. Required Components

### 1.1 Software Requirements
- **Go version**: 1.20 or higher
- **MySQL**: 8.0+
  - Required database: `eino_interview`
  - Test database: `eino_interview_test`
- **Redis**: 6.0+
  - Default port: 6379
  - Test instance should be separate from production
- **Milvus**: 2.0+
  - Vector database for resume embedding storage
- **Hertz framework**: Latest stable version
- **Volcano Ark AI services**: API access and credentials

### 1.2 Development Tools
- **GoLand** or **VS Code** with Go extensions
- **Git** for version control
- **Docker** and **Docker Compose** for containerized testing
- **Postman** or **Insomnia** for API testing

## 2. Environment Configuration

### 2.1 Environment Variables

Create a `.env.test` file with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=test_user
DB_PASSWORD=test_password
DB_NAME=eino_interview_test

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=1

# Milvus Configuration
MILVUS_HOST=localhost
MILVUS_PORT=19530

# AI Service Configuration
AI_SERVICE_URL=http://localhost:8000/v1
AI_API_KEY=test_api_key

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8888

# JWT Configuration
JWT_SECRET=test_jwt_secret
JWT_EXPIRATION=3600

# Logging Configuration
LOG_LEVEL=debug
LOG_FILE=./logs/test.log

# Test Configuration
TEST_MODE=true
DISABLE_EXTERNAL_SERVICES=false
```

### 2.2 Test Database Schema

Ensure the test database has the following tables:
- `users`
- `interviews`
- `questions`
- `answers`
- `evaluations`
- `sessions`

Use the same schema as production but with test data.

### 2.3 Docker Compose for Test Environment

Create a `docker-compose.test.yml` file:

```yaml
version: '3'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: eino_interview_test
      MYSQL_USER: test_user
      MYSQL_PASSWORD: test_password
    ports:
      - "3307:3306"
    volumes:
      - mysql_test_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:6.0
    ports:
      - "6380:6379"
    volumes:
      - redis_test_data:/data

  milvus:
    image: milvusdb/milvus:v2.2.8
    ports:
      - "19530:19530"
      - "19121:19121"
    volumes:
      - milvus_test_data:/var/lib/milvus

  # Mock AI Service
  mock-ai-service:
    build: ./test/mock-ai
    ports:
      - "8000:8000"

volumes:
  mysql_test_data:
  redis_test_data:
  milvus_test_data:
```

## 3. Test Dependencies

### 3.1 Go Test Dependencies

Add to go.mod:

```
require (
    github.com/golang/mock v1.6.0
    github.com/bytedance/sonic v1.8.0
    github.com/go-playground/assert/v2 v2.2.0
    github.com/go-playground/validator/v10 v10.11.2
    github.com/go-redis/redis/v8 v8.11.5
    github.com/golang-jwt/jwt/v4 v4.4.3
    github.com/stretchr/testify v1.8.0
    github.com/tidwall/gjson v1.14.0
    github.com/tidwall/sjson v1.2.5
    gorm.io/driver/mysql v1.4.6
    gorm.io/gorm v1.24.3
    github.com/bytedance/hertz v0.5.0
)
```

### 3.2 Mock Services

Create mock implementations for:
- AI evaluation service
- Redis message broker
- Milvus vector database
- User authentication service

### 3.3 Test Data

Prepare test data including:
- Test user accounts
- Sample resume content
- Mock interview sessions
- Expected AI responses
- Test code snippets (both correct and incorrect)

## 4. Setup Instructions

### 4.1 Local Development Setup

1. Install required software components
2. Create test database and user
3. Configure environment variables
4. Initialize Go modules
5. Start required services (MySQL, Redis, Milvus)
6. Run database migrations for test schema

### 4.2 Docker-Based Setup

1. Build test images:
   ```bash
   docker-compose -f docker-compose.test.yml build
   ```

2. Start test environment:
   ```bash
   docker-compose -f docker-compose.test.yml up -d
   ```

3. Run database migrations:
   ```bash
   go run ./cmd/migrate/main.go --config .env.test
   ```

4. Seed test data:
   ```bash
   go run ./cmd/seed/main.go --config .env.test --data test
   ```

### 4.3 CI/CD Integration

For CI/CD pipelines, configure:
- Test database instances
- Service containers
- Environment variables
- Test execution steps
- Code coverage reporting

## 5. Cleanup Procedures

After testing, ensure proper cleanup:
- Stop test services
- Remove test containers
- Clear test data
- Cleanup temporary files
- Close database connections

## 6. Verification Steps

To verify test environment setup:
1. Check database connectivity
2. Verify Redis connection
3. Test AI service mock responses
4. Validate API endpoints are reachable
5. Run a sample test case to confirm everything works