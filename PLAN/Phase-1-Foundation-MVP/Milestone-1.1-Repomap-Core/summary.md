# Summary: Milestone 1.1 – Repomap Core Implementation

## Quick Overview

**Milestone 1.1** delivers the **repomap** tool—the core MVP that generates token-optimized, high-density repository maps. It encompasses the complete pipeline from file discovery through ranking and structured output.

## What Gets Built

A production-ready **repomap** binary that analyzes Go repositories and outputs a concise, ranked summary of all files, their definitions, and import relationships.

### Key Deliverables

| Deliverable | Description |
|------------|-------------|
| **repomap binary** | Single Go executable (Linux, macOS, Windows) |
| **Source code** | Well-structured Go packages: discovery, parsing, graph, ranking, output |
| **Unit tests** | >85% code coverage on all packages |
| **Integration tests** | Full pipeline tests on real Go repositories |
| **Performance benchmarks** | Recorded metrics for discovery, parsing, ranking |
| **User documentation** | Guides in `repomap/doc/` (QUICK_START.md, CLI_REFERENCE.md, etc.) |
| **Code documentation** | godoc comments on all exported APIs |

## Phases & Workflow

### Phase A: Discovery & Filtering (2 tasks)
- Walk file system using `filepath.WalkDir`
- Parse and apply `.gitignore` rules
- Filter by file type and exclude binary files

### Phase B: Parsing & Extraction (2 tasks)
- Use Go's `go/ast` to parse source files
- Extract function, method, struct, and interface signatures
- Extract import statements

### Phase C: Import Graph (1 task)
- Build directed graph from import relationships
- Normalize import paths

### Phase D: Ranking (2 tasks)
- Calculate in-degree for each file
- Compute importance scores (0.0 - 1.0)
- Assign importance labels (high/medium/low)

### Phase E: Output & CLI (6 tasks)
- Define output data structures
- Implement XML rendering
- Implement JSON rendering
- Integrate token counting (MVP: simple approximation)
- Implement CLI flag parsing
- Create main integration layer

### Testing & Tooling (3 tasks)
- Unit tests across all packages
- Integration tests on real repositories
- Performance benchmarking

### Documentation & Release (2 tasks)
- Add godoc comments to all exported functions
- Create build scripts and release binaries

## Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Build Status | No errors on Go 1.19+ | TBD |
| Test Coverage | >85% across all packages | TBD |
| Integration Tests | Pass on 5+ real Go repos | TBD |
| Performance | <5s on 1K files, <30s on 10K | TBD |
| Discovery Accuracy | 100% of Go files discovered, gitignore respected | TBD |
| Parsing Accuracy | All top-level definitions extracted | TBD |
| Graph Accuracy | Import relationships correct | TBD |
| Output Validity | Valid XML/JSON, respects token budget | TBD |
| Documentation | Complete user guide and API docs | TBD |

## Timeline

| Week | Phase | Key Milestones |
|------|-------|----------------|
| Week 1 | Discovery & Filtering | File walker complete, gitignore parsing complete |
| Week 2 | Parsing & Graph | Go AST extraction complete, import graph complete |
| Week 2–3 | Ranking & Output | Ranking algorithm complete, XML/JSON rendering complete |
| Week 3 | Testing & Release | All tests passing, binaries released |

**Total Duration:** 3 weeks (2–3 weeks estimated)

## Dependencies

### External
- Go 1.19+ compiler
- Standard library modules only (no external packages)

### Internal
- Documentation in `repomap/doc/` (ARCHITECTURE.md, CLI_REFERENCE.md, etc.)

### Blocking
- None at project start

## Risks & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|-----------|
| Large repos (100K+) timeout | Medium | High | Limit Phase 1 to 100K files; optimize in Phase 2 |
| Complex .gitignore patterns | Low | Medium | Use simple patterns first; switch to library if needed |
| Token counting inaccuracy | Low | Medium | Use approximation in MVP; integrate `tiktoken-go` later |

## Acceptance Criteria

✅ **Milestone 1.1 is COMPLETE when:**

1. `repomap` binary compiles without errors
2. File discovery correctly finds and filters Go files per `.gitignore`
3. Go AST parsing extracts all top-level definitions
4. Import graph is accurately constructed
5. File ranking is based on in-degree and normalized to [0.0, 1.0]
6. XML output is valid and respects token budget
7. JSON output is valid and respects token budget
8. CLI interface supports all required flags
9. >85% unit test coverage across all packages
10. Integration tests pass on 5+ real Go repositories
11. Performance: <5s on 1K files, <30s on 10K files
12. User documentation is complete and accurate
13. All code is documented with godoc comments
14. Binary releases created for Linux, macOS, Windows

## Next Phase: Milestone 1.2

Upon completion, **Milestone 1.2 (CLI Framework & Integration)** will:
- Extract reusable CLI patterns from repomap
- Create a base CLI framework for future tools
- Establish standard output formatting conventions
- Document CLI best practices for the project

## Who Should Read This

- **Project Leads:** Understand scope and timeline
- **Developers:** Use as implementation guide
- **QA:** Use as testing checklist
- **Product:** Track against success metrics
