# Implementation Strategy: The MVP Approach

This document outlines the step-by-step implementation strategy, starting with a **Minimum Viable Product (MVP)** focused on Go-to-Go mapping.

---

## Why Start with Go?

1. **Simplicity:** Go's standard library includes `go/ast` for AST parsing
2. **No External Dependencies:** Tree-sitter can be added later
3. **Clear Use Case:** Many Go projects exist and need mapping
4. **Transferable Logic:** Once mastered for Go, scaling to other languages is straightforward

---

## MVP Implementation Phases

### Phase 1: Directory Walking with Gitignore Support

**Goal:** Walk the file system and respect `.gitignore` rules.

**Steps:**

1. Import `github.com/karrick/godirwalk` and `github.com/monochromegane/go-gitignore`
2. Create a walker function that:
   - Traverses directories using `godirwalk.Walk`
   - Loads `.gitignore` files from the root and subdirectories
   - Skips ignored files and binary files
   - Collects all `.go` files
3. Output a list of discovered Go files

**Deliverable:** A function that returns a slice of file paths

```go
func discoverGoFiles(root string) ([]string, error)
```

**Success Criteria:** Successfully ignore `node_modules`, `vendor`, `.git`, and other standard exclude patterns

---

### Phase 2: AST Parsing and Signature Extraction

**Goal:** Extract function and type signatures from Go files using `go/ast`.

**Steps:**

1. For each discovered `.go` file:
   - Read the file contents
   - Parse using `go/ast.Parse`
   - Traverse the AST to extract:
     - Function signatures: `func Name(params) return`
     - Method receivers: `func (receiver) Name(params)`
     - Type definitions: `type Name struct { ... }`
     - Interface definitions: `type Name interface { ... }`
2. Store extracted definitions in `FileNode.Definitions`

**Deliverable:** A function that extracts signatures from a single file

```go
func extractDefinitions(filePath string) ([]string, error)
```

**Example Output:**
```
func main()
func NewServer(port int) *Server
func (s *Server) Start()
type Server struct
type Config struct
```

**Success Criteria:** Correctly extract all top-level declarations without reading implementation details

---

### Phase 3: Import Graph Construction

**Goal:** Build a directed graph of import relationships.

**Steps:**

1. For each parsed file:
   - Extract `import` statements from the AST
   - Normalize import paths to relative file paths
   - Create edges in the dependency graph
2. Use Go's `gonum/graph` or similar library to represent the graph

**Deliverable:** A function that builds the import graph

```go
func buildImportGraph(nodes map[string]*FileNode) Graph
```

**Example Output:**
```
main.go → [utils.go, config.go]
server.go → [utils.go]
handler.go → [utils.go, auth.go]
```

**Success Criteria:** Correctly identify all import relationships without false positives

---

### Phase 4: File Ranking

**Goal:** Calculate importance scores based on import centrality.

**Steps:**

1. For each file in the graph:
   - Calculate in-degree (how many files import it)
   - Normalize to a rank score (0.0 - 1.0)
   - Assign `Rank` to the FileNode
2. Sort files by rank in descending order

**Deliverable:** A function that ranks files

```go
func rankFiles(nodes map[string]*FileNode, graph Graph) map[string]*FileNode
```

**Example Output:**
```
Rank 0.95: utils.go (in-degree: 5)
Rank 0.60: config.go (in-degree: 2)
Rank 0.30: logger.go (in-degree: 1)
Rank 0.10: handler.go (in-degree: 0)
```

**Success Criteria:** Files with more imports have higher ranks

---

### Phase 5: Token Counting and Output

**Goal:** Generate XML/JSON output respecting a token budget.

**Steps:**

1. Import `github.com/pkoukk/tiktoken-go` for token counting
2. Process files in rank order:
   - Calculate tokens needed for each file's output
   - Add to output if within budget
   - Stop when budget exhausted
3. Generate XML (or JSON) output

**Deliverable:** Output formatting functions

```go
func renderXML(nodes map[string]*FileNode, maxTokens int) (string, error)
func renderJSON(nodes map[string]*FileNode, maxTokens int) (string, error)
```

**Example Output:**
```xml
<repomap>
  <file path="utils.go" importance="high" rank="0.95">
    <definition>func Log(msg string)</definition>
    <definition>func Parse(data []byte) Config</definition>
  </file>
  <file path="config.go" importance="medium" rank="0.60">
    <definition>type Config struct</definition>
  </file>
</repomap>
```

**Success Criteria:** Output respects token budget and displays files in importance order

---

## MVP Architecture

```
repomap/
├── cmd/
│   └── repomap/
│       └── main.go              # CLI entry point
├── internal/
│   ├── discovery/
│   │   └── walker.go            # Phase 1: File discovery
│   ├── parsing/
│   │   └── extractor.go         # Phase 2: AST parsing
│   ├── graph/
│   │   └── builder.go           # Phase 3: Import graph
│   ├── ranking/
│   │   └── ranker.go            # Phase 4: Ranking
│   └── output/
│       ├── xml.go               # Phase 5: XML rendering
│       └── json.go              # Phase 5: JSON rendering
├── go.mod
└── README.md
```

---

## Development Timeline

| Phase | Tasks | Est. Time |
|-------|-------|-----------|
| 1 | File walking + gitignore | 1-2 days |
| 2 | AST parsing + extraction | 2-3 days |
| 3 | Graph construction | 1-2 days |
| 4 | Ranking algorithm | 1 day |
| 5 | Output rendering | 1-2 days |
| **Total** | **MVP Complete** | **~6-10 days** |

---

## Post-MVP: Multi-Language Support

After the MVP is complete, add support for other languages:

1. **Replace** `go/ast` with `go-tree-sitter` in the parsing phase
2. **Generalize** the extractor to handle any language supported by Tree-sitter
3. **Enhance** the CLI with `--include-lang`, `--exclude-lang` flags
4. **Test** with Python, JavaScript, and Rust projects

---

## Testing Strategy

### Unit Tests
- Test signature extraction for various Go constructs
- Test ranking algorithm with mock graphs
- Test token counting accuracy

### Integration Tests
- Run against real Go projects (small, medium, large)
- Verify import graph correctness
- Validate output format

### Performance Tests
- Benchmark file walking on 10K+ file repositories
- Measure memory usage with large import graphs
- Validate token counting speed

