# Milestone 3.2: API-Forge – Tasks

## Task 1: Design API Detection Strategy
- **Objective:** Define patterns for REST, gRPC, GraphQL detection
- **Deliverable:** Detection algorithm design document
- **Acceptance Criteria:** Clear patterns defined for each API type

## Task 2: Build REST Endpoint Detector
- **Objective:** Implement detection of HTTP endpoints (routes, handlers)
- **Patterns:** Express, FastAPI, Gin route definitions
- **Deliverable:** REST detector for multiple frameworks
- **Acceptance Criteria:** Detects endpoints in test apps

## Task 3: Build gRPC Service Detector
- **Objective:** Extract gRPC services from proto files
- **Features:** Service and RPC definition extraction
- **Deliverable:** gRPC detector
- **Acceptance Criteria:** Extracts all services and methods

## Task 4: Build GraphQL Schema Detector
- **Objective:** Extract GraphQL schema from code
- **Features:** Query, mutation, subscription detection
- **Deliverable:** GraphQL detector
- **Acceptance Criteria:** Extracts complete schema

## Task 5: Implement OpenAPI Generator
- **Objective:** Create `pkg/schema/openapi.go` for OpenAPI 3.0 generation
- **Features:** Endpoint specs, parameters, responses
- **Deliverable:** OpenAPI generator
- **Acceptance Criteria:** Generates valid OpenAPI schemas

## Task 6: Implement gRPC Schema Generator
- **Objective:** Generate protobuf definitions from code
- **Deliverable:** protobuf generator
- **Acceptance Criteria:** Generates valid .proto files

## Task 7: Implement GraphQL Schema Generator
- **Objective:** Generate GraphQL schema definitions
- **Deliverable:** GraphQL schema generator
- **Acceptance Criteria:** Generates valid GraphQL schema files

## Task 8: Build Documentation Generator
- **Objective:** Create `pkg/docs/generator.go` for API documentation
- **Content:** Descriptions, examples, error codes
- **Deliverable:** Documentation generator
- **Acceptance Criteria:** Generates readable docs

## Task 9: Implement Example Generator
- **Objective:** Generate request/response examples
- **Examples:** cURL, Python, JavaScript examples
- **Deliverable:** Example generator
- **Acceptance Criteria:** Examples are executable

## Task 10: Build Python SDK Generator
- **Objective:** Generate Python client libraries
- **Features:** Type hints, async support, error handling
- **Deliverable:** Python SDK generator
- **Acceptance Criteria:** Generated SDKs work correctly

## Task 11: Build JavaScript SDK Generator
- **Objective:** Generate JavaScript/TypeScript client libraries
- **Features:** Promise/async-await, type definitions
- **Deliverable:** JavaScript SDK generator
- **Acceptance Criteria:** Generated SDKs work correctly

## Task 12: Build Go SDK Generator
- **Objective:** Generate Go client libraries
- **Features:** Idiomatic Go, error handling
- **Deliverable:** Go SDK generator
- **Acceptance Criteria:** Generated SDKs work correctly

## Task 13: Implement Type Detection
- **Objective:** Extract and map parameter/response types
- **Features:** Primitive types, complex types, generics
- **Deliverable:** Type detector
- **Acceptance Criteria:** Types detected accurately

## Task 14: Implement Error Code Mapping
- **Objective:** Detect and document error codes and exceptions
- **Mapping:** HTTP status → error description
- **Deliverable:** Error mapper
- **Acceptance Criteria:** Errors documented

## Task 15: Build Authentication Detector
- **Objective:** Detect authentication methods (API key, OAuth, JWT)
- **Deliverable:** Auth detector
- **Acceptance Criteria:** Detects common auth patterns

## Task 16: Create api-forge CLI
- **Objective:** Build `cmd/api-forge/main.go` with CLI interface
- **Flags:** --source, --api-type, --output, --sdk-lang
- **Deliverable:** CLI tool executable
- **Acceptance Criteria:** Tool accepts input and generates output

## Task 17: Implement Schema Export
- **Objective:** Export schemas in standard formats
- **Formats:** JSON, YAML, protobuf
- **Deliverable:** Schema exporters
- **Acceptance Criteria:** Exports valid formats

## Task 18: Build SDK Template System
- **Objective:** Create customizable SDK generation templates
- **Features:** Per-language templates, customization
- **Deliverable:** Template system
- **Acceptance Criteria:** Templates work for all languages

## Task 19: Implement Incremental Generation
- **Objective:** Update existing SDKs instead of regenerating
- **Features:** Merge new methods, preserve custom code
- **Deliverable:** Incremental generator
- **Acceptance Criteria:** Updates work correctly

## Task 20: Create Integration Tests
- **Objective:** Test on real APIs (OpenAPI.json files, sample projects)
- **Test Cases:** REST, gRPC, GraphQL samples
- **Deliverable:** Integration tests
- **Acceptance Criteria:** All tests passing

## Task 21: Build Documentation Generator
- **Objective:** Generate markdown documentation from schemas
- **Format:** Human-readable API docs
- **Deliverable:** Doc generator
- **Acceptance Criteria:** Docs readable and complete

## Task 22: Implement Validation
- **Objective:** Validate generated schemas and SDKs
- **Checks:** Schema validity, SDK compilability
- **Deliverable:** Validator
- **Acceptance Criteria:** Catches invalid output

## Task 23: Create Performance Benchmarks
- **Objective:** Measure detection and generation speed
- **Benchmarks:** Time for large APIs
- **Deliverable:** Benchmark report
- **Acceptance Criteria:** Performance acceptable

## Task 24: Create User Documentation
- **Objective:** Write guide for API detection and SDK generation
- **Content:** Setup, examples, advanced usage
- **Deliverable:** docs/API_FORGE.md
- **Acceptance Criteria:** Clear guide with examples

## Task 25: Final Testing and Optimization
- **Objective:** End-to-end testing and performance optimization
- **Deliverable:** All tests passing, optimized
- **Acceptance Criteria:** Ready for production use
