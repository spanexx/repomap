# Milestone 4.1: Server Mode & API

## Objective

Deploy agents-cli tools as a unified REST API server enabling agent orchestration and integration.

## Scope

### In-Scope

1. **REST API Gateway**
   - Unified API for all tools
   - Request routing and handling
   - Response formatting
   - Error handling and logging

2. **Tool Integration**
   - Expose all Phase 1-3 tools as API endpoints
   - Request/response schemas
   - Streaming output support
   - Job queuing and async execution

3. **Server Infrastructure**
   - HTTP server (Go net/http)
   - Request validation
   - Rate limiting
   - Request logging and tracing

4. **agents-server Tool** (`cmd/agents-server/`)
   - CLI for running server
   - Configuration management
   - Health checks
   - Graceful shutdown

### Out-of-Scope

- Web UI (Phase 4.2+)
- Advanced caching beyond tool-specific
- Custom middleware

## Deliverables

1. **API Gateway** (`pkg/gateway/`)
2. **Tool API Adapters** (`pkg/adapters/`)
3. **Job Queue** (`pkg/queue/`)
4. **agents-server Tool** (`cmd/agents-server/`)
5. **API Documentation** (OpenAPI)
6. **Integration & Deployment Guides**

## Success Criteria

- ✅ All tools accessible via REST API
- ✅ API response <500ms for most operations
- ✅ Async execution with job tracking
- ✅ >80% code coverage

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** ~20–25

---

## API Endpoints Overview

```
POST   /api/v1/repomap           - Analyze repository
POST   /api/v1/semgrep           - Semantic search
POST   /api/v1/web-ray           - Web navigation
POST   /api/v1/sandbox-run       - Code execution
POST   /api/v1/api-forge         - API detection
POST   /api/v1/git-surgeon       - Merge resolution
POST   /api/v1/mem-kv            - Memory operations
POST   /api/v1/calc-bridge       - Math computation
GET    /api/v1/health            - Health check
GET    /api/v1/jobs/{id}         - Job status
```

## Dependencies

- All Phase 1-3 tools completion
