# Data Structures

## FileNode

Represents the "skeleton" of a single file in the codebase.

```go
type FileNode struct {
    Path        string      // File path relative to root
    Language    string      // Programming language (go, python, js, etc.)
    Imports     []string    // List of other files this file depends on
    Definitions []string    // e.g., "func NewServer()", "type Config struct"
    Rank        float64     // Calculated importance score (0.0 - 1.0)
    TokenCount  int         // How many tokens needed to represent this node
}
```

### Fields Explained

- **Path:** The relative file path from the repository root (e.g., `pkg/auth/login.go`)
- **Language:** The detected programming language for specialized parsing
- **Imports:** A slice of import paths extracted from the file, representing dependencies
- **Definitions:** Simplified signatures of functions, types, structs, and classes defined in the file
- **Rank:** A calculated importance score based on the file's centrality in the import graph
- **TokenCount:** The number of tokens required to represent this node's information in the output

## RepoMap

Represents the complete analyzed repository structure and relationships.

```go
type RepoMap struct {
    Nodes map[string]*FileNode       // Map of file path -> FileNode
    Graph *simple.DirectedGraph      // Import dependency graph for ranking
}
```

### Fields Explained

- **Nodes:** A map keyed by file path, containing all analyzed FileNode instances
- **Graph:** A directed graph where edges represent import relationships:
  - An edge from `A` to `B` means "A imports B"
  - Used to calculate in-degree and centrality scores

## Relationships

The data structures work together as follows:

1. **Discovery Phase:** Populate the `Nodes` map with FileNode entries for each file
2. **Parsing Phase:** Fill in the `Imports` and `Definitions` fields for each FileNode
3. **Ranking Phase:** Build the `Graph` from the import relationships and calculate `Rank` scores
4. **Rendering Phase:** Traverse `Nodes` in rank order, respecting the `TokenCount` budget

### Import Graph Example

Given these relationships:
- `main.go` imports `utils.go` and `config.go`
- `server.go` imports `utils.go`
- `logger.go` imports nothing

The graph would look like:
```
main.go ──┐
          ├──> utils.go
server.go─┘      │
                └──> (dependency)

config.go

logger.go
```

In this example, `utils.go` has the highest in-degree (2), making it rank as more important.
