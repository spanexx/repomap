# Milestone 3.2: Quick Reference

## Summary

API schema auto-detection and SDK generation. Analyze code to find REST, gRPC, and GraphQL APIs; generate client libraries for Go, Python, JavaScript.

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 22 of Phase 3

## Example Usage

```bash
# Detect REST API and generate OpenAPI schema
api-forge --source ./backend --api-type rest --output schema.json

# Generate Python SDK from detected API
api-forge --source ./api --output-sdk python --lang python

# Generate all SDKs
api-forge --source ./myapp --output-sdk all --output-dir ./sdks

# Extract gRPC schema
api-forge --source ./protos --api-type grpc --output proto-schema.json
```

## Key Components

| Component | Path | Purpose |
|-----------|------|---------|
| REST Detector | `pkg/apidetect/rest.go` | HTTP endpoint detection |
| gRPC Detector | `pkg/apidetect/grpc.go` | Service detection |
| GraphQL Detector | `pkg/apidetect/graphql.go` | Schema extraction |
| Schema Generators | `pkg/schema/*.go` | Format-specific generators |
| SDK Generators | `pkg/codegen/*.go` | Language-specific codegen |
| CLI Tool | `cmd/api-forge/main.go` | Command-line interface |

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| API Types Supported | 3 | — |
| SDK Languages | 3 | — |
| Detection Accuracy | >90% | — |
| Code Coverage | >80% | — |

## Acceptance Criteria

- ✅ Detects all major API patterns
- ✅ Generates valid OpenAPI/proto schemas
- ✅ SDKs are functional and tested
- ✅ >90% detection accuracy
- ✅ >80% code coverage
- ✅ Complete documentation

## Key Metrics

| Metric | Target |
|--------|--------|
| Detection Accuracy | >90% |
| SDK Correctness | 100% |
| Code Coverage | >80% |
