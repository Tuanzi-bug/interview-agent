# 构建系统

**最后更新**: 2026-02-01  
**维护者**: DevOps 团队

## 概述

本文档描述面试吧平台的构建、打包和部署流程。

## Docker 构建

### 后端 Dockerfile

```dockerfile
# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./main"]
```

### 前端 Dockerfile

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ .
RUN npm run build

# Runtime stage
FROM nginx:alpine

COPY --from=builder /app/.next /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## Docker Compose

### 开发环境 (docker-compose.yml)

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: interview
    ports:
      - "3306:3306"
    volumes:
      - ./backend/db_schema.sql:/docker-entrypoint-initdb.d/schema.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

  backend:
    build:
      context: .
      dockerfile: backend/Dockerfile
    ports:
      - "8080:8080"
    environment:
      DB_HOST: mysql
      REDIS_HOST: redis
    depends_on:
      - mysql
      - redis

  frontend:
    build:
      context: .
      dockerfile: frontend/Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - backend
```

### 生产环境 (docker-compose-prod.yml)

```yaml
version: '3.8'

services:
  backend:
    image: interview-backend:latest
    restart: always
    environment:
      ENV: production
      DB_HOST: ${DB_HOST}
      REDIS_HOST: ${REDIS_HOST}
    # 更多生产配置...

  frontend:
    image: interview-frontend:latest
    restart: always
    # 更多生产配置...
```

## Makefile

```makefile
# Makefile
.PHONY: help build run test clean deploy

help:
	@echo "Available commands:"
	@echo "  make build       - Build Docker images"
	@echo "  make run         - Run development environment"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deploy      - Deploy to production"

build:
	docker-compose build

run:
	docker-compose up

test:
	go test ./... -v -cover

stop:
	docker-compose down

clean:
	docker-compose down -v
	rm -rf ./backend/build ./frontend/.next

deploy:
	docker-compose -f docker-compose-prod.yml up -d

logs:
	docker-compose logs -f

db-migrate:
	mysql -h localhost -u root -proot interview < backend/db_schema.sql
```

## 常用命令

```bash
# 开发环境
make run              # 启动开发环境
make test             # 运行测试
make stop             # 停止服务
make clean            # 清理

# 构建
make build            # 构建 Docker 镜像

# 部署
make deploy           # 部署到生产环境

# 查看日志
make logs             # 查看日志
```

---

**维护者**: DevOps 团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
