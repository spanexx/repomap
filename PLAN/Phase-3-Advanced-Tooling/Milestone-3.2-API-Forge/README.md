# Milestone 3.2: API-Forge

## Objective

Automatically detect API schemas and patterns from code, generating comprehensive API documentation and client SDKs.

## Scope

### In-Scope

1. **API Detection**
   - HTTP endpoint discovery (REST APIs)
   - gRPC service detection
   - GraphQL schema extraction
   - Function signatures as APIs

2. **Schema Generation**
   - OpenAPI/Swagger generation
   - Protocol Buffers schema extraction
   - JSON schema creation
   - Type definitions

3. **Documentation Generation**
   - Endpoint descriptions with examples
   - Parameter and response documentation
   - Error code mapping
   - Authentication guide

4. **SDK Generation**
   - Client library generation (Python, JavaScript, Go)
   - Type-safe API calls
   - Error handling utilities

### Out-of-Scope

- Mobile SDK generation
- Custom language SDKs
- API testing framework

## Deliverables

1. **API Detection Engine** (`pkg/apidetect/`)
2. **Schema Generator** (`pkg/schema/`)
3. **SDK Generator** (`pkg/codegen/`)
4. **api-forge Tool** (`cmd/api-forge/`)
5. **Documentation & Examples**

## Success Criteria

- ✅ Detects all major API patterns
- ✅ Generates accurate schemas
- ✅ Creates functional client SDKs
- ✅ >80% code coverage
- ✅ Generated SDKs work correctly

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** ~25–30

---

## Supported API Types

| API Type | Format | Priority |
|----------|--------|----------|
| REST (HTTP) | OpenAPI | High |
| gRPC | protobuf | Medium |
| GraphQL | GraphQL schema | Medium |
| Function APIs | Custom format | Low |

---

## Dependencies

- Phase 1 completion
- Phase 2.1 (multi-language parsing)
