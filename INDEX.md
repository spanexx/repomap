# Repomap

### 1. The Architecture

The tool needs to act as a funnel: taking a massive file system and reducing it to a high-density "map."

**Pipeline:**

1. **Discovery:** Walk the directory tree (fast).
2. **Filtering:** Apply `.gitignore` and exclude binary/irrelevant files.
3. **Parsing:** Use **Tree-sitter** (via Go bindings) to extract definitions (functions, classes, structs) without reading the whole file body.
4. **Ranking:** Score files based on "centrality" (how many other files import them).
5. **rendering:** Generate the output tree, cutting off details once the `max-tokens` budget is reached.

---

### 2. Recommended Go Libraries

Don't reinvent the wheel. These libraries are industry standards for this stack:

* **File Walking:** `github.com/karrick/godirwalk` (Much faster than `filepath.WalkDir` for massive repos).
* **Git Ignore:** `github.com/monochromegane/go-gitignore` (Crucial for not mapping `node_modules` or `vendor`).
* **Parsing (The Core):** `github.com/smacker/go-tree-sitter` (Go bindings for Tree-sitter. This is how you get robust parsing for Python, JS, Go, Rust, etc., without writing regexes).
* **Token Counting:** `github.com/pkoukk/tiktoken-go` (To accurately stick to the `--max-tokens` budget).

---

### 3. The Data Structure (Go)

You need a struct that represents the "Skeleton" of a file.

```go
type FileNode struct {
	Path        string
	Language    string
	Imports     []string // List of other files this file depends on
	Definitions []string // e.g., "func NewServer()", "type Config struct"
	Rank        float64  // Calculated importance score
	TokenCount  int      // How expensive this node is to print
}

type RepoMap struct {
	Nodes map[string]*FileNode
	Graph *simple.DirectedGraph // To calculate PageRank/Centrality
}

```

---

### 4. The Logic Flow

#### Phase A: The "Skeletons" (Parsing)

Instead of reading the file content, you only want the *signatures*.

* **Input (Go file):**
```go
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // ... 50 lines of logic ...
}

```


* **Repomap Extraction:**
```go
"func (s *Server) HandleRequest(w, r)"

```


* *Why:* This saves ~95% of tokens but keeps the context of "what this file does."

#### Phase B: The "Ranking" (Context Optimization)

If you have 1000 files but only 2000 tokens of budget, which files do you show?

* **Heuristic:** Files that are imported by many others (e.g., `utils.go`, `types.go`) or lie at the root of the AST are likely more important for context.
* **Implementation:** Build a simple directed graph where `A imports B` creates an edge `A -> B`. Calculate the "In-Degree" (number of incoming edges) for each file.
* **Sort:** Sort `FileNodes` by `Rank` descending.

#### Phase C: The Output (XML/JSON)

LLMs are excellent at parsing XML tags for structure.

```xml
<repomap>
  <file path="main.go" importance="high">
    func main()
    func initDB()
  </file>
  <file path="pkg/auth/login.go" importance="medium">
    struct LoginRequest
    func Authenticate(user, pass)
  </file>
  <file path="pkg/utils/logger.go" importance="low">
    </file>
</repomap>

```

---

### 5. CLI Usage Design

Design the flags for the Agent, not for a human.

```bash
# Basic usage
repomap --root ./src --json

# Strict budget for a small context window
repomap --max-tokens 1000 --output xml

# Focus on specific languages
repomap --include-ext .go,.js --ignore-tests

```

### 6. Implementation Strategy: "The MVP"

Don't try to support every language immediately. Start with **Go** mapping **Go**.

1. **Step 1:** Create a walker that lists all `.go` files and respects `.gitignore`.
2. **Step 2:** Use `go/ast` (standard library) instead of Tree-sitter for the MVP. It's built-in and easier to start with for Go files.
* *Note:* Switch to Tree-sitter later when you want to support Python/JS.


3. **Step 3:** Print a tree showing only `struct` names and `func` signatures.

### Next Step

Would you like me to write the **Go code for "Step 2"** (using `go/ast` to extract function signatures from a file and return a simplified string)?