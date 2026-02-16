# Resume Upload Optimization - Deployment Guide

This guide provides instructions for deploying and monitoring the Resume Upload Optimization system.

## 1. System Overview

The system consists of the following components:
- **Backend (Go/Hertz)**: Core API and background workers.
- **MySQL**: Persistent storage for metadata and job status.
- **Redis**: Job queue, progress broadcasting (Pub/Sub), and PDF text cache.
- **Frontend (Next.js)**: User interface for resume management and real-time progress.

## 2. Infrastructure Requirements

- **Docker**: Engine 20.10+
- **Docker Compose**: V2.0+
- **Resources**: 
  - Minimum: 2 CPU, 4GB RAM
  - Recommended: 4 CPU, 8GB RAM (especially if running AI models locally)

## 3. Configuration

### 3.1 Environment Variables

The following environment variables can be configured in your `.env` file or directly in the Docker Compose environment section.

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_HOST` | MySQL database host | `mysql` | `localhost` |
| `DB_PORT` | MySQL database port | `3306` | `3307` |
| `DB_USER` | MySQL username | `root` | `admin` |
| `DB_PASSWORD` | MySQL password | `root` | `secret_pass` |
| `DB_NAME` | MySQL database name | `interview_agent` | `resume_db` |
| `REDIS_HOST` | Redis host | `redis` | `localhost` |
| `REDIS_PORT` | Redis port | `6379` | `6379` |
| `REDIS_PASSWORD` | Redis password | `root` | `redis_pass` |
| `NEXT_PUBLIC_API_BASE_URL` | Frontend API base URL | `http://localhost:8888/api` | `https://api.example.com/api` |

### 3.2 System Constants and Timeouts

The following parameters are currently defined as constants within the system. To modify them, a code change and recompilation are required.

| Parameter | Location | Value | Description |
|-----------|----------|-------|-------------|
| `WorkerShutdownTimeout` | `worker/resume_worker_pool.go` | `180s` | Period to wait for active jobs during shutdown. |
| `DefaultJobTimeout` | `worker/resume_worker_pool.go` | `120s` | Maximum time for a single AI analysis job. |
| `DequeueTimeout` | `worker/resume_worker_pool.go` | `5s` | Redis BRPOP timeout for workers. |
| `QuickSyncTimeout` | `handler/interview/interviews_service.go` | `15s` | Sync wait time before falling back to async. |
| `PDFTextCacheTTL` | `service/resume_text_cache.go` | `30 days` | TTL for extracted PDF text in Redis. |
| `SSE Timeout` | `handler/interview/progress_handler.go` | `300s` | SSE connection timeout. |

### 3.3 Application Configuration (`backend/config.yaml`)

The backend uses `backend/config.yaml` for infrastructure connections:

```yaml
# backend/config.yaml
database:
  dsn: "root:root@tcp(mysql:3306)/interview_agent?charset=utf8mb4&parseTime=True&loc=Local"

redis:
  addr: "redis:6379"
  password: "root"

hertz:
  read_timeout: "300s"   # Supports large file uploads
  write_timeout: "300s"  # Supports long AI analysis
```

## 4. Docker Deployment

### 4.1 Production Deployment

Use `docker-compose-prod.yml` for production-ready setup.

```bash
# Start all services
docker compose -f docker-compose-prod.yml up -d
```

### 4.2 Health Checks

The system includes automated health checks for all critical services:

- **MySQL**: Checks availability using `mysqladmin ping`.
- **Redis**: Checks availability using `redis-cli ping`.
- **Backend**: Health check endpoint at `GET /health`.

```bash
# Check service health status
docker compose ps
```

## 5. Background Worker Pool

The backend automatically spawns a worker pool on startup.

- **Initialization**: Log message `Resume upload worker pool initialized with 3 workers`.
- **Concurrency**: Default 3 workers handle background resume extraction and AI analysis.
- **Resilience**: Workers automatically retry failed Redis connections and continue processing the next job if a specific job fails.
- **Graceful Shutdown**: The pool waits up to 180 seconds for active jobs to complete before shutting down.

## 6. Monitoring & Logging

### 6.1 Log Analysis

Logs are stored in `./backend/logs` and can be followed using:

```bash
# Monitor backend logs
docker logs -f app-backend
```

**Key Log Prefixes:**
- `[WorkerPool]`: Worker status updates, job start/completion, and errors.
- `[ResumeUploadQueue]`: Job enqueue/dequeue and progress broadcasting.
- `[PDFTextCache]`: Cache hits and misses.

### 6.2 Redis Monitoring

Monitor the job queue and progress updates:

```bash
# Monitor the job queue length
docker exec app-redis redis-cli -a root LLEN resume:upload:queue

# Monitor real-time progress broadcasts
docker exec app-redis redis-cli -a root psubscribe "resume:progress:*"
```

### 6.3 Database Monitoring

Track job status in the `resume_upload_status` table:

```bash
# Check table structure and indexes
docker exec -i app-mysql mysql -u root -proot interview_agent -e "DESC resume_upload_status; SHOW INDEX FROM resume_upload_status;"

# View recent job status
docker exec -i app-mysql mysql -u root -proot interview_agent -e "SELECT upload_id, status, progress, stage, error_msg FROM resume_upload_status ORDER BY created_at DESC LIMIT 10;"
```

## 7. Scaling Considerations

- **Horizontal Scaling**: Multiple backend instances can be deployed. They will compete for jobs in the single Redis queue (`resume:upload:queue`), providing automatic load balancing.
- **Worker Tuning**: If CPU usage is low but the queue is growing, increase the `worker_pool.size` in `config.yaml`.
- **Redis Pub/Sub**: Ensure your load balancer supports sticky sessions or long-lived connections for SSE (Server-Sent Events) progress streaming.

## 8. Troubleshooting

| Issue | Potential Cause | Solution |
|-------|-----------------|----------|
| **Always Async** | AI analysis consistently takes > 15s | Expected behavior. Ensure frontend handles the `is_async: true` response. |
| **Worker Stuck** | Network issues with AI provider | Check `job_timeout` settings. Workers should timeout after 120s. |
| **No Progress Events** | Redis Pub/Sub blocked or disconnected | Verify Redis health and network connectivity between Backend and Redis. |
| **Database Lock** | Multiple workers updating same record | The system uses `upload_id` unique constraints; check MySQL process list for long-running locks. |
| **413 Payload Too Large** | Resume file exceeds Nginx/Hertz limit | Increase `client_max_body_size` in Nginx and Hertz request limit. |

## 9. Rollback Procedures

In case of a failed deployment or critical bug:

1. **Stop Current Version**:
   ```bash
   docker compose -f docker-compose-prod.yml down
   ```

2. **Revert Code/Image**:
   Revert to the previous stable git tag or update the image tag in `docker-compose-prod.yml`.

3. **Restore Database (If Schema Changed)**:
   If the deployment included incompatible database migrations:
   ```bash
   # Restore from backup
   cat backup.sql | docker exec -i app-mysql mysql -u root -proot interview_agent
   ```

4. **Restart Services**:
   ```bash
   docker compose -f docker-compose-prod.yml up -d
   ```

## 10. Backup and Recovery

- **Database**: Regular backups of `mysql-data` volume.
- **Cache**: Redis data in `redis-data` volume is persistent via AOF (`appendonly yes`).
- **Resumes**: Uploaded PDF files should be backed up if not stored in a distributed object store (e.g., MinIO/S3).
