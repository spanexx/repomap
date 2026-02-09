# Milestone 1.1: Repomap Core Implementation

## Objective

Implement the complete **repomap** tool—the MVP that generates token-optimized, high-density repository maps. This milestone delivers the core functionality: file discovery, AST parsing, import graph construction, ranking, and structured output (XML/JSON).

## Scope

### In-Scope

1. **Phase A: Discovery**
   - Fast directory walking using `godirwalk` (or standard `filepath.WalkDir`)
   - `.gitignore` parsing and file filtering
   - Support for nested `.gitignore` files

2. **Phase B: Parsing**
   - Go AST parsing using `go/ast` (standard library)
   - Extract function, method, struct, and interface signatures
   - Simplify signatures to token-efficient strings

3. **Phase C: Import Graph**
   - Build directed graph from import relationships
   - Normalize import paths to relative file paths
   - Prepare for ranking in Phase D

4. **Phase D: Ranking**
   - Calculate in-degree (import frequency) for each file
   - Normalize to rank scores (0.0 - 1.0)
   - Sort files by importance

5. **Phase E: Output**
   - XML format rendering
   - JSON format rendering
   - Token counting using `tiktoken-go`
   - Respect `--max-tokens` budget
   - Include importance labels and truncation notices

6. **CLI Interface**
   - `--root` flag for repository path
   - `--output` flag (xml/json)
   - `--max-tokens` flag for budget
   - `--include-ext`, `--exclude-ext` for filtering
   - `--ignore-tests` flag
   - `--verbose` flag for debugging

### Out-of-Scope

- Multi-language support (deferred to Phase 2)
- Tree-sitter integration (deferred to Phase 2)
- Server mode or distributed features (deferred to Phase 4+)
- Integration with external agents (deferred to Phase 3+)

## Deliverables

1. **Binary:** `repomap` executable (single Go binary)
2. **Documentation:** Complete user guide and API documentation (in `repomap/doc/`)
3. **Source Code:** Well-structured Go modules in `repomap/` directory
   - `cmd/repomap/main.go` – CLI entry point
   - `internal/discovery/` – File walking
   - `internal/parsing/` – AST extraction
   - `internal/graph/` – Import graph construction
   - `internal/ranking/` – Importance scoring
   - `internal/output/` – XML/JSON rendering

## Success Criteria

- ✅ `repomap` compiles without errors on Go 1.19+
- ✅ Correctly discovers and filters Go files using `.gitignore`
- ✅ Parses Go files and extracts all top-level definitions
- ✅ Builds accurate import graphs from Go `import` statements
- ✅ Ranks files by in-degree and importance
- ✅ Generates valid XML and JSON output
- ✅ Respects `--max-tokens` budget and truncates gracefully
- ✅ Processes 10K-file repositories in under 30 seconds
- ✅ All functionality covered by unit tests
- ✅ `--help` and `--version` flags work
- ✅ Documentation is complete and accurate

## Key Implementation Details

### Technologies & Libraries

- **Go 1.19+** – Standard library only for MVP (no external dependencies initially)
- **Parsing:** `go/ast` (standard library)
- **File Walking:** `filepath.WalkDir` or `filepath.Walk`
- **Token Counting:** `tiktoken-go` (add later if needed for strict budgets)

### Architecture

```
repomap/
├── cmd/repomap/
│   └── main.go                  # CLI entry point
├── internal/
│   ├── discovery/
│   │   ├── walker.go            # File discovery + gitignore
│   │   └── gitignore.go         # Gitignore parsing
│   ├── parsing/
│   │   ├── extractor.go         # Go AST extraction
│   │   └── types.go             # Definition types
│   ├── graph/
│   │   ├── builder.go           # Import graph construction
│   │   └── graph.go             # Graph data structure
│   ├── ranking/
│   │   ├── ranker.go            # Importance scoring
│   │   └── score.go             # Score calculation
│   └── output/
│       ├── xml.go               # XML rendering
│       ├── json.go              # JSON rendering
│       └── types.go             # Output data structures
├── repomap/                     # Exported public types
│   └── repomap.go               # Main API
├── go.mod
├── go.sum
└── README.md
```

### Implementation Phases (Sequential)

1. **Week 1:** Discovery (Phase A) + Parsing (Phase B)
   - File walking with `.gitignore` support
   - Go AST extraction
   - Unit tests for both

2. **Week 2:** Graph + Ranking (Phases C & D)
   - Import graph construction
   - Ranking algorithm
   - Integration tests

3. **Week 3:** Output + CLI (Phase E + interface)
   - XML/JSON rendering
   - Command-line interface
   - Full integration tests
   - Performance testing

## Testing Strategy

### Unit Tests
- Test `.gitignore` filtering with various patterns
- Test Go AST extraction on real code samples
- Test ranking algorithm with mock graphs
- Test output formatting (XML/JSON validity)

### Integration Tests
- Run against small, medium, and large Go repositories
- Verify correct file counts and rankings
- Validate output structure

### Performance Tests
- Benchmark on 1K, 10K, 50K, and 100K file counts
- Measure memory usage
- Target: 100K files in < 30 seconds

## Definition of Done

- All code merged to `main` branch
- All tests passing
- Code reviewed and approved
- Documentation complete
- Performance benchmarks recorded
- Binary releases created

## Dependencies

- Go 1.19+ (compiler)
- Standard library modules: `go/ast`, `flag`, `os`, `path/filepath`, `io/ioutil`

## Open Questions

1. Should we use `monochromegane/go-gitignore` or parse gitignore manually?
   - **Decision:** Start with simple `.gitignore` parsing; switch to library if complexity grows

2. Should token counting be accurate or approximate in MVP?
   - **Decision:** Approximate (character count / 4) in MVP; integrate `tiktoken-go` if LLM integration requires it

3. What's the maximum repository size we should support in Phase 1?
   - **Decision:** 100K files; larger repos deferred to Phase 2 with optimization

## Next Steps

1. Create Go project structure (go.mod, initial packages)
2. Implement Phase A (discovery) with tests
3. Implement Phase B (parsing) with tests
4. Implement Phase C (graph) with tests
5. Implement Phase D (ranking) with tests
6. Implement Phase E (output) with CLI
7. Full integration testing and performance tuning
8. Release and documentation finalization
