# Milestone 3.1: Quick Reference

## Summary

Sandboxed code execution environment. Execute Go, Python, JavaScript, Rust, and Java code safely with resource limits and comprehensive output capture.

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 18 of Phase 3

## Example Usage

```bash
# Simple code execution
sandbox-run --lang python --code "print('Hello')"

# File execution with resource limits
sandbox-run --lang go --file main.go --timeout 5s --memory 512m

# With output capture
sandbox-run --lang javascript --code "console.log('test')" --output json

# Complex example
sandbox-run --lang python --code "
import requests
response = requests.get('http://localhost')
print(response.status_code)
" --timeout 10s --network localhost
```

## Key Components

| Component | Path | Purpose |
|-----------|------|---------|
| Container Manager | `pkg/sandbox/container.go` | Container lifecycle |
| Resource Control | `pkg/sandbox/resources.go` | CPU/memory limits |
| Executor | `pkg/executor/executor.go` | Code compilation & execution |
| Monitor | `pkg/monitor/monitor.go` | Execution metrics |
| CLI Tool | `cmd/sandbox-run/main.go` | Command-line interface |

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| Languages Supported | 5 | — |
| Execution Success | >95% | — |
| Code Coverage | >80% | — |
| Startup Time | <2s | — |

## Acceptance Criteria

- ✅ Executes code in isolated container
- ✅ Enforces resource limits
- ✅ Captures stdout/stderr
- ✅ Provides execution metrics
- ✅ >95% execution success rate
- ✅ >80% code coverage
- ✅ Handles errors gracefully

## Key Metrics

| Metric | Target |
|--------|--------|
| Container Startup | <2s |
| Execution Success | >95% |
| CPU Overhead | <10% |
| Memory Overhead | <50MB |
