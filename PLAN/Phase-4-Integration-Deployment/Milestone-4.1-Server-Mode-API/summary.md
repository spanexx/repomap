# Milestone 4.1: Quick Reference

## Summary

REST API server exposing all agents-cli tools with job queue, monitoring, and deployment readiness.

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** 25 sequential

## Key Metrics

| Metric | Target |
|--------|--------|
| Response Latency | <500ms |
| Availability | >99.9% |
| Throughput | 100+ req/s |
| Code Coverage | >80% |

## API Examples

```bash
# Analyze repository
curl -X POST http://localhost:8080/api/v1/repomap \
  -H "Content-Type: application/json" \
  -d '{"root": "/path/to/repo"}'

# Semantic search
curl -X POST http://localhost:8080/api/v1/semgrep \
  -d '{"query": "find authentication logic", "root": "/path"}'

# Execute code
curl -X POST http://localhost:8080/api/v1/sandbox-run \
  -d '{"lang": "python", "code": "print(42)"}'

# Check job status
curl http://localhost:8080/api/v1/jobs/{job_id}
```

## Acceptance Criteria

- ✅ All tools accessible via API
- ✅ Request validation working
- ✅ Async execution with job tracking
- ✅ <500ms response latency
- ✅ >99.9% availability
- ✅ >80% code coverage
