# Resume Upload System - Performance Analysis

## Executive Summary

The resume upload optimization system implements a hybrid synchronous/asynchronous processing architecture designed to maximize user experience (UX) while ensuring system reliability. By leveraging Redis-based text caching, an asynchronous worker pool, and real-time Server-Sent Events (SSE), the system achieves sub-second response times for cached resumes and maintains responsiveness for complex AI-intensive processing.

**Key Performance Highlights:**
- **Cache Hit Efficiency**: Reduces processing time from seconds to milliseconds (~1-5ms), achieving a 99%+ reduction in latency for previously processed documents.
- **Responsive Hybrid Flow**: Provides immediate results for small documents (< 15s) while gracefully falling back to background processing with real-time feedback for larger files.
- **Scalable Processing**: A fixed worker pool (3 workers) handles concurrent background tasks, preventing resource exhaustion during peak loads.
- **Real-time Feedback**: SSE connection establishment in < 100ms ensures users are immediately informed of their upload status.

## Performance Metrics

### Overview Table

| Metric | Target | Expected | Production | Notes |
|--------|--------|----------|------------|-------|
| Small PDF (< 500KB) Sync Time | < 5s | 2-4s | TBD | Depends on AI provider latency |
| Cache Hit Retrieval Time | < 1s | 1-10ms | TBD | Redis GET + Reconstruction |
| Quick Sync Timeout | 15s | 15s | 15s | Fixed system constant |
| Worker Job Processing | < 120s | 30-60s | TBD | Max limit 120s |
| SSE Connection Establishment | < 100ms | 20-50ms | TBD | Handshake + initial event |
| SSE Event Latency | < 50ms | 5-15ms | TBD | Redis Pub/Sub delivery |
| Database Update Latency | < 10ms | 1-5ms | TBD | Standard MySQL write |
| Redis Queue Throughput | > 100 j/s | > 500 j/s | TBD | Limited by network/CPU |
| Concurrent Upload Support | 50+ | 100+ | TBD | Scalable via horizontal pods |

### Cache Performance

The system utilizes a content-based SHA256 hashing strategy to identify resumes.

- **Hit Rate Expectations**: Expected 20-40% in production as users often refine and re-upload similar versions of their resumes.
- **Miss Rate Impact**: A cache miss adds approximately 2-10ms to the total processing time (Redis SET operation), which is negligible compared to PDF parsing (500ms-5s).
- **TTL Effectiveness**: The 30-day TTL (2,592,000s) balances storage efficiency with the likelihood of re-uploads within a job application cycle.

### Upload Mode Distribution

| File Size / Complexity | Mode | Probability | Rationale |
|-----------------------|------|-------------|-----------|
| < 500KB (Simple) | Sync | 80% | Usually parses and analyzes within 15s |
| 500KB - 2MB | Mixed | 50% | Dependent on AI analysis depth |
| > 2MB or Scanned | Async | 95%+ | Extraction and AI analysis typically > 15s |
| Any (Cache Hit) | Sync | 100% | Retrieval is near-instant (< 10ms) |

### Processing Time Breakdown

| Phase | Duration (Miss) | Duration (Hit) | Resource Usage |
|-------|-----------------|----------------|----------------|
| PDF Text Extraction | 500ms - 5s | 1ms - 5ms | High CPU (Miss) / Low (Hit) |
| AI Analysis | 2s - 60s | N/A (Cached) | Low (I/O Bound) |
| DB / Redis Overhead | 10ms - 50ms | 5ms - 10ms | Low |
| **Total End-to-End** | **3s - 70s** | **< 100ms** | |

## Worker Pool Performance

- **Throughput**: With 3 workers and an average job time of 30s, the system handles 6 jobs per minute per instance.
- **Concurrency**: 3 concurrent jobs are processed without context switching overhead, matching typical dual-core or quad-core cloud environments.
- **Queue Depth**: The Redis list (`resume:upload:queue`) can store thousands of jobs with minimal memory footprint (~200 bytes per job).

## Redis & Database Performance

### Redis Operations
- **LPUSH/BRPOP**: Atomic operations ensure no job is lost or double-processed. Latency is typically < 1ms.
- **Pub/Sub**: Used for progress broadcasting. Subscription establishment takes ~50ms; message delivery is near-instant.
- **Cache**: SHA256 hashing takes ~10-50ms for a 5MB PDF, significantly faster than OCR or complex parsing.

### Database Operations
- **Status Updates**: Occur at 0%, 25%, 50%, 75%, and 100%. Each update is a simple indexed write by `upload_id`.
- **Indexing**: `upload_id` has a unique index, ensuring `O(1)` lookups for status checks and updates.

## SSE Performance

- **Connection Establishment**: < 100ms including JWT authentication and initial "connected" event.
- **Latency**: Progress events are streamed immediately upon Redis Pub/Sub trigger, bypassing database polling.
- **Concurrency**: Each SSE connection consumes one goroutine (~2KB memory). A single instance can comfortably support 1,000+ concurrent listeners.

## Benchmarking Methodology

### Test Environment Setup
- **Hardware**: 4 CPU Cores, 8GB RAM.
- **Services**: MySQL 8.0, Redis 7.0 (Running in Docker).
- **Network**: Localhost (to eliminate external AI provider variance, use mock AI parser).

### Test Scenarios

1. **Quick Sync Benchmark**:
   - Upload 100KB PDF.
   - Target: Response < 5s, `is_async: false`.
2. **Async Fallback Benchmark**:
   - Upload 5MB PDF (or use mock with 20s delay).
   - Target: Immediate response < 1s, `is_async: true`.
3. **Cache Hit Benchmark**:
   - Upload same file twice.
   - Target: Second response < 100ms.
4. **Concurrent Load Test**:
   - 50 simultaneous uploads.
   - Target: No 5xx errors, all jobs eventually reach 100%.

### Measurement Tools
- `ab` (Apache Bench) or `wrk` for API throughput.
- `time curl` for individual request latency.
- `docker stats` for resource utilization.
- `redis-cli --latency` for infrastructure health.

## Optimization Recommendations

### Short-term (Quick Wins)
- **Pre-calculating Hash**: Calculate SHA256 during the initial file read to avoid secondary I/O.
- **Gzip SSE Events**: Enable compression for SSE streams if event frequency increases.
- **Redis Connection Pooling**: Ensure `max_idle` and `max_active` are tuned for 100+ concurrent connections.

### Medium-term (Architectural)
- **Dynamic Worker Scaling**: Adjust worker pool size based on Redis queue length.
- **Multi-stage Caching**: Cache AI analysis results separately from raw text to allow partial re-use.
- **Local OCR**: For scanned PDFs, integrate a lightweight local OCR (e.g., Tesseract) to avoid AI extraction costs.

### Long-term (Scaling)
- **Dedicated Worker Nodes**: Move the `WorkerPool` to a separate microservice to isolate CPU-intensive parsing from the API.
- **Global CDN Caching**: For frequently accessed resumes, cache results at the edge.
- **Vector DB Pre-warming**: Simultaneously push extracted text to Milvus during the AI analysis phase.

## Capacity Planning

### Single Server Capacity (Estimated)
- **Concurrent Users**: 500 active sessions.
- **Upload Throughput**: 10 uploads/sec (sustained).
- **Processing Throughput**: 6-10 jobs/minute (with 3 workers).
- **Storage**: ~100KB per cache entry; 1 million entries ≈ 100GB.

### Resource Requirements
- **CPU**: 2 Cores (API) + 2 Cores (Workers).
- **Memory**: 2GB (OS/Runtime) + 1GB (Redis Cache) + 1GB (Buffer).
- **Network**: 100Mbps (supports 10x 1MB uploads/sec).

## Performance Testing Procedures

### Manual Testing
```bash
# Test Cache Hit
time curl -X POST http://localhost:8888/api/resume/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@test.pdf"

# Test SSE Latency
curl -N "http://localhost:8888/api/resume/upload/progress/$UPLOAD_ID?token=$TOKEN"
```

### Automated Load Testing
```bash
# Using wrk to test upload throughput
wrk -t4 -c50 -d30s --timeout 60s http://localhost:8888/api/resume/upload
```

## Monitoring Recommendations

### Key Metrics to Track
- `resume_upload_cache_hit_ratio`: Percentage of uploads served from cache.
- `resume_upload_sync_success_rate`: Percentage of uploads completing within 15s.
- `resume_worker_queue_length`: Number of jobs waiting in Redis.
- `resume_processing_duration_seconds`: Histogram of total processing time.

### Alerting Thresholds
- **Queue Length > 50**: Indicates workers cannot keep up with load.
- **Sync Success Rate < 50%**: Suggests AI provider latency or complex document increase.
- **Error Rate > 1%**: Investigate worker logs for processing failures.

## Appendix

### Performance Troubleshooting
- **Slow Sync**: Check AI provider API latency; consider reducing `QuickSyncTimeout` if AI consistently takes > 15s.
- **High Redis Memory**: Monitor `used_memory_human`; adjust `PDFTextCacheTTL` or implementation eviction policy (LRU).
- **SSE Drops**: Check Nginx `proxy_read_timeout` (should be >= 300s).

### Related Documentation
- [Resume Upload Deployment Guide](../deployment/resume-upload-deployment.md)
- [API Specification](../api/resume-api.md)
