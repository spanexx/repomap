# Tasks: Milestone 1.1 – Repomap Core Implementation

## Dependencies
- **Prerequisites:** None (Foundation Milestone)
- **Provides for:**
  - Milestone 1.2: CLI logic (1.1.13) and Output logic (1.1.10, 1.1.11) for framework extraction.
  - Milestone 1.3: Core logic for exhaustive testing and benchmarking.


## Phase A: Discovery & Filtering

### Task 1.1.1: Project Setup & Go Module Initialization
- [ ] Create `repomap/` directory structure
- [ ] Initialize `go.mod` with module path `github.com/spanexx/agents-cli/repomap`
- [ ] Create `cmd/repomap/main.go` with basic CLI skeleton
- [ ] Create `internal/discovery/`, `internal/parsing/`, etc. package directories
- [ ] Verify project builds: `go build -o repomap ./cmd/repomap`

**Acceptance Criteria:**
- Go modules resolve without errors
- Binary compiles and runs `repomap --help`

---

### Task 1.1.2: Implement File Discovery (Walker)
- [ ] Implement `discovery.Walk(root string) ([]string, error)` function
- [ ] Use `filepath.WalkDir` or `filepath.Walk` to traverse directories
- [ ] Filter out binary files (`.exe`, `.o`, `.a`, `.so`, `.dylib`)
- [ ] Filter by file extension (Go files by default in MVP)
- [ ] Write unit tests for Walk() with test fixtures

**Acceptance Criteria:**
- Correctly discovers all `.go` files in test repository
- Skips binary files and directories like `.git/`
- Unit tests pass with >90% code coverage

---

### Task 1.1.3: Implement .gitignore Parsing & Filtering
- [ ] Implement `gitignore.Parse(root string) (*Gitignore, error)`
- [ ] Read and parse `.gitignore` file from repository root
- [ ] Implement pattern matching (at least basic glob support)
- [ ] Integrate with `Walk()` to filter ignored files
- [ ] Handle nested `.gitignore` files (optional for MVP)
- [ ] Write unit tests for gitignore matching

**Acceptance Criteria:**
- Correctly excludes files/dirs matching `.gitignore` patterns
- Handles common patterns (`node_modules/`, `*.pyc`, `build/`)
- Unit tests pass for real-world `.gitignore` files

**Example:** Given `.gitignore` with `vendor/` and `*_test.go`:
- Files in `vendor/` are excluded
- `*_test.go` files are excluded (if filtering enabled)
- Other `.go` files are included

---

## Phase B: Parsing & Extraction

### Task 1.1.4: Go AST Parser – Extract Definitions
- [ ] Implement `parsing.ExtractDefinitions(filePath string) ([]string, error)`
- [ ] Use `go/ast.Parse()` to parse Go source files
- [ ] Extract function signatures: `func Name(params) return`
- [ ] Extract method signatures: `func (receiver) Name(params)`
- [ ] Extract type definitions: `type Name struct { ... }`
- [ ] Extract interface definitions: `type Name interface { ... }`
- [ ] Simplify signatures to remove unnecessary detail
- [ ] Write unit tests on real Go code samples

**Acceptance Criteria:**
- Correctly extracts all top-level declarations
- Signatures are simplified and token-efficient
- Handles syntax errors gracefully
- Unit tests cover 10+ real Go files

**Example Output:**
```
func main()
func NewServer(port int) *Server
func (s *Server) Start() error
type Server struct
type Config struct
```

---

### Task 1.1.5: Go AST Parser – Extract Imports
- [ ] Implement `parsing.ExtractImports(filePath string) ([]string, error)`
- [ ] Use `go/ast` to traverse import declarations
- [ ] Collect all imported package paths
- [ ] Normalize import paths to relative file paths
  - Example: `github.com/user/repo/pkg/utils` → `pkg/utils/...`
- [ ] Write unit tests on files with various import patterns

**Acceptance Criteria:**
- Correctly extracts all imports from Go files
- Handles `import "x"`, `import (...)`, aliased imports
- Normalizes paths correctly for internal imports
- Unit tests cover various import patterns

---

## Phase C: Import Graph Construction

### Task 1.1.6: Build Import Graph
- [ ] Implement `graph.Builder` struct with methods:
  - `AddFile(path string, imports []string)`
  - `Build() *Graph`
- [ ] Represent graph as adjacency list or matrix
- [ ] Validate graph integrity (no cycles, all imports resolved)
- [ ] Write unit tests on mock file sets

**Acceptance Criteria:**
- Graph correctly represents import relationships
- In-degree calculation is accurate
- Can handle 10K+ nodes without memory issues
- Unit tests pass with mock graphs

**Example:** Given files with imports:
- `main.go` imports `utils.go`, `config.go`
- `server.go` imports `utils.go`

Expected graph:
```
main.go ──┐
          ├──> utils.go
server.go─┘
config.go
```

---

## Phase D: Ranking & Importance Scoring

### Task 1.1.7: File Ranking Algorithm
- [ ] Implement `ranking.Rank(graph *Graph) map[string]float64`
- [ ] Calculate in-degree for each node
- [ ] Normalize to rank scores (0.0 - 1.0)
  - Formula: `rank[file] = in_degree[file] / max_in_degree`
- [ ] Sort files by rank descending
- [ ] Write unit tests on mock graphs

**Acceptance Criteria:**
- Ranking scores are normalized to [0.0, 1.0]
- Files with more dependencies rank higher
- Scores sum correctly
- Unit tests verify ranking on multiple graphs

**Example:** Given in-degrees [3, 2, 0, 1]:
- Ranks: [1.0, 0.67, 0.0, 0.33]

---

### Task 1.1.8: Importance Level Assignment
- [ ] Implement importance labeling: `high` (>0.7), `medium` (0.3-0.7), `low` (<0.3)
- [ ] Assign to each FileNode
- [ ] Write unit tests

**Acceptance Criteria:**
- Importance levels assigned correctly based on rank
- All FileNodes have an importance level

---

## Phase E: Output & CLI

### Task 1.1.9: Output Data Structures
- [ ] Define `FileNode` struct:
  ```go
  type FileNode struct {
    Path        string
    Language    string
    Importance  string   // high, medium, low
    Rank        float64
    Definitions []string
    TokenCount  int
  }
  ```
- [ ] Define `RepoMap` struct with file list and metadata
- [ ] Write validation methods

**Acceptance Criteria:**
- Structs are well-defined and exported
- All fields populated correctly

---

### Task 1.1.10: XML Output Rendering
- [ ] Implement `output.RenderXML(nodes []*FileNode) (string, error)`
- [ ] Generate XML with structure:
  ```xml
  <repomap>
    <file path="..." importance="..." rank="...">
      <definition>...</definition>
    </file>
  </repomap>
  ```
- [ ] Respect token budget: stop adding files when budget exhausted
- [ ] Include truncation notice if budget exceeded
- [ ] Write unit tests

**Acceptance Criteria:**
- Valid XML output
- Respects token budget
- Includes all required fields
- Unit tests verify XML structure

---

### Task 1.1.11: JSON Output Rendering
- [ ] Implement `output.RenderJSON(nodes []*FileNode) (string, error)`
- [ ] Generate JSON with structure:
  ```json
  {
    "repomap": {
      "files": [
        {
          "path": "...",
          "importance": "...",
          "rank": ...,
          "definitions": [...]
        }
      ]
    }
  }
  ```
- [ ] Respect token budget (same as XML)
- [ ] Include truncation notice
- [ ] Write unit tests

**Acceptance Criteria:**
- Valid JSON output
- Matches expected schema
- Respects token budget
- Unit tests verify JSON validity

---

### Task 1.1.12: Token Counting
- [ ] Implement simple token counter (MVP): `TokenCount(text) int`
  - **MVP Formula:** `len(text) / 4` (rough approximation)
- [ ] Integrate with output rendering
- [ ] Track cumulative tokens while building output
- [ ] Write unit tests

**Acceptance Criteria:**
- Token counts are reasonable approximations
- Total output respects `--max-tokens` budget
- Gracefully handles budget exhaustion

---

### Task 1.1.13: CLI Interface & Flags
- [ ] Implement command-line flag parsing:
  - `--root <path>` (repository root)
  - `--output <xml|json>` (output format)
  - `--max-tokens <number>` (budget)
  - `--include-ext <exts>` (comma-separated extensions)
  - `--exclude-ext <exts>` (comma-separated)
  - `--ignore-tests` (boolean flag)
  - `--verbose` (debug output)
  - `--help` / `-h` (help message)
  - `--version` / `-v` (version)
- [ ] Implement flag validation
- [ ] Implement help text
- [ ] Write CLI tests

**Acceptance Criteria:**
- All flags parse correctly
- Help text is clear and complete
- Invalid flags produce error messages
- CLI tests pass

---

### Task 1.1.14: Main CLI Integration
- [ ] Implement `main()` function to:
  1. Parse flags
  2. Call `discovery.Walk()`
  3. Call `parsing.Extract*()` for each file
  4. Call `graph.Build()`
  5. Call `ranking.Rank()`
  6. Call `output.Render*()`
  7. Print to stdout
- [ ] Handle errors gracefully
- [ ] Add timing/performance logging (with `--verbose`)
- [ ] Write integration tests

**Acceptance Criteria:**
- CLI runs end-to-end without errors
- Output is correct and respects all flags
- Integration tests pass on real repositories

---

## Testing & Quality (Initial Baseline)

> **Note:** These tasks provide initial confidence. Milestone 1.3 will consolidate and expand these into a full test suite with CI/CD integration.

### Task 1.1.15: Essential Unit Tests
- [ ] Achieve >70% coverage for core logic (Parsing, Graph, Ranking)
- [ ] Verify error handling for basic file operations
- [ ] **Dependency:** Required before proceeding to Milestone 1.2 refactoring.

**Acceptance Criteria:**
- Core logic has unit tests
- Coverage >70%
- All core tests pass: `go test ./...`

---

### Task 1.1.16: MVP Integration Tests
- [ ] Create one small Go repo fixture
- [ ] Verify `repomap` runs end-to-end on this fixture
- [ ] **Dependency:** Milestone 1.3 Task 1.3.10 will expand this to multi-scale fixtures.

**Acceptance Criteria:**
- Integration tests pass on small fixture
- Output is structurally correct

---

### Task 1.1.17: Performance Benchmarking
- [ ] Create benchmarks for key functions
- [ ] Run: `go test -bench=. -benchmem ./...`
- [ ] Document results

**Acceptance Criteria:**
- Benchmarks recorded
- Discovery: <100ms for 1K files
- Parsing: <1s for 1K files
- Overall: <5s for 10K files

---

## Documentation & Release

### Task 1.1.18: Code Documentation
- [ ] Add godoc comments to all exported functions
- [ ] Document package purposes
- [ ] Add code examples to critical functions

**Acceptance Criteria:**
- All exported items documented
- `go doc` output is clear and helpful

---

### Task 1.1.19: Build & Release
- [ ] Create `Makefile` or build script
- [ ] Generate releases for Linux, macOS, Windows
- [ ] Create checksums
- [ ] Document build process

**Acceptance Criteria:**
- Binary builds without errors
- All platforms work
- Releases are downloadable

---

## Summary

**Total Tasks:** 19
**Estimated Effort:** 2–3 weeks
**Dependencies:** Go 1.19+
**Success:** All tasks completed + tests passing + documentation complete
