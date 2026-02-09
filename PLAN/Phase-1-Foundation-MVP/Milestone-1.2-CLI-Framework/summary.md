# Summary: Milestone 1.2 – CLI Framework & Integration

## Quick Overview

**Milestone 1.2** extracts reusable patterns from repomap to create a shared CLI framework that all future agent tools can leverage. This establishes consistent conventions for flags, output, configuration, and error handling across the agents-cli project.

## What Gets Built

A modular CLI framework (`pkg/cli/`, `pkg/output/`, `pkg/config/`, etc.) that future tools like `semgrep-cli` and `web-ray` can use to eliminate boilerplate and maintain consistency.

### Key Deliverables

| Deliverable | Description |
|------------|-------------|
| **CLI Framework** | Reusable flag parsing, help text, app structure |
| **Output Formatters** | XML, JSON, and text/table formatting |
| **Configuration Management** | Config files, env vars, CLI flag hierarchy |
| **Error Handling** | Standardized error types and messages |
| **Shared Utilities** | Path helpers, token counting, filtering, logging |
| **Documentation** | Framework guide, examples, best practices |
| **Refactored Repomap** | Uses framework without changing CLI interface |

## Phases & Workflow

### Framework Foundation (3 tasks)
- Build core CLI app structure
- Implement flag parsing abstraction
- Generate help text automatically

### Output Formatting (5 tasks)
- Define output writer interface
- Extract and generalize XML formatter
- Extract and generalize JSON formatter
- Implement text/table formatter
- Create formatter factory

### Configuration Management (3 tasks)
- Load configuration files
- Support environment variables
- Implement override hierarchy (defaults → env → config → flags)

### Error Handling (2 tasks)
- Define error types and exit codes
- Create standardized error messages

### Shared Utilities (4 tasks)
- Path normalization and operations
- Token counting
- File filtering
- Logging abstraction

### Integration & Documentation (2 tasks)
- Refactor repomap to use the framework
- Comprehensive framework documentation

### Testing & Validation (3 tasks)
- Unit tests (>80% coverage)
- Integration tests
- Compatibility tests with repomap

## Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Code Coverage | >80% on framework packages | TBD |
| Reusability | Used by 2+ tools in Phase 2 | TBD |
| Breaking Changes | None to repomap CLI | TBD |
| Performance | <10% startup overhead | TBD |
| Documentation | Complete with examples | TBD |

## Timeline

| Week | Phase | Key Milestones |
|------|-------|----------------|
| Week 4 | Framework + Output | CLI framework, flag parsing, output formatters |
| Week 5 | Config + Testing | Config management, error handling, refactoring, tests |

**Total Duration:** 1 week (after Milestone 1.1)

## Dependencies

### External
- Go 1.19+ (same as Phase 1)
- Standard library only

### Internal
- **Blocks:** Milestone 1.1 must be complete before starting 1.2
- **Enables:** Phase 2 tools (semgrep-cli, web-ray) can be built on this framework

## Architecture Overview

```
pkg/
├── cli/           # App structure, flags, help
├── output/        # Writer interface, formatters (XML, JSON, text)
├── config/        # Config loading, environment vars, hierarchy
├── errors/        # Error types, codes, messages
└── util/          # Paths, tokens, filtering, logging
```

## Key Design Decisions

1. **Simple, extensible framework** – Don't over-engineer; make it easy to add formatters
2. **Configuration hierarchy** – CLI flags override environment variables, which override config files
3. **Pluggable output** – Tools can add custom formatters without modifying core
4. **Reuse from repomap** – Extract tested patterns rather than rewriting
5. **Minimal breaking changes** – Repomap CLI interface stays the same

## Acceptance Criteria

✅ **Milestone 1.2 is COMPLETE when:**

1. CLI framework package is well-structured and documented
2. Output formatters work with any data type (not just repomap)
3. Configuration hierarchy works correctly (tested)
4. Error handling is consistent and user-friendly
5. Utilities are reusable across multiple tools
6. Repomap is refactored to use the framework
7. All repomap tests still pass (no regressions)
8. >80% code coverage on all framework packages
9. Integration tests pass on repomap + mock tools
10. Documentation includes working examples
11. Performance overhead is <10% startup time
12. Code is reviewed and approved

## Next Phase: Milestone 1.3

Upon completion, **Milestone 1.3 (Testing Infrastructure & CI/CD)** will:
- Set up comprehensive test suite for both Repomap and Framework
- Configure CI/CD pipeline for automated testing
- Establish performance benchmarking
- Document testing procedures and standards

## Future Tools Built on This Framework

The framework is designed to support:

- **semgrep-cli** (Phase 2) – Semantic grep using embeddings
- **web-ray** (Phase 2) – Headless browser with structured output
- **sandbox-run** (Phase 3) – Sandboxed command execution
- **api-forge** (Phase 3) – API schema detection and correction
- **git-surgeon** (Phase 3) – Git merge conflict resolution
- **mem-kv** (Phase 3) – Persistent key-value memory store
- **calc-bridge** (Phase 3) – Math and symbolic computation

Each tool will use `pkg/cli/`, `pkg/output/`, `pkg/config/` without reimplementing common logic.

## Who Should Read This

- **Framework Developers:** Use as implementation guide
- **Tool Developers (Phase 2+):** Learn how to build tools using the framework
- **QA/Testing:** Understand what gets tested and how
- **Project Leads:** Track progress and dependencies
