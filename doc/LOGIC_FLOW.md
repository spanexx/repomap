# Logic Flow: The Three Phases

The repomap process is divided into three distinct phases that work together to transform a massive codebase into a dense, useful summary.

---

## Phase A: Skeletons (Parsing Signatures)

### Goal
Extract function, type, and struct definitions from files **without** reading their full implementations.

### Why It Matters
A typical function might have:
- 1 line of signature
- 50+ lines of implementation logic

By extracting only the signature, we save ~95% of tokens while retaining the essential context of "what this file does."

### Example

**Input (Go file):**
```go
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // ... 50 lines of logic ...
    // validation, database calls, response formatting, etc.
}
```

**Repomap Extraction:**
```
func (s *Server) HandleRequest(w, r)
```

### Implementation Approach

1. Parse the file using the language-specific parser (e.g., `go/ast`)
2. Traverse the AST and collect only top-level declarations
3. Convert declarations to simplified signatures
4. Store in `FileNode.Definitions`

### Types of Definitions to Extract

- **Functions:** `func Name(params)`
- **Methods:** `func (receiver) Name(params)`
- **Structs/Types:** `type Name struct`
- **Interfaces:** `type Name interface`
- **Classes (in other languages):** `class Name`

---

## Phase B: Ranking (Centrality Calculation)

### Goal
Determine the importance of each file based on its role in the codebase structure.

### The Problem
If you have 1000 files but only 2000 tokens of budget, which files do you include?

### The Heuristic
Files that are imported by many others are likely more important for understanding the system.

**Example Importance Scores:**
- `utils.go` (imported by 15 other files): **High importance**
- `types.go` (imported by 8 files): **Medium importance**
- `logger.go` (imported by 2 files): **Low importance**
- `example_test.go` (imported by 0 files): **Very low importance**

### Implementation Steps

1. **Build the Import Graph**
   - Create a directed graph where each node is a file
   - For each import relationship, add an edge: `A → B` (A imports B)

2. **Calculate In-Degree**
   ```
   In-Degree(file) = number of other files that import this file
   ```

3. **Normalize Scores**
   - Convert in-degree to a ranked score (0.0 - 1.0)
   - Files with higher in-degree get higher ranks

4. **Sort Files**
   - Order `FileNode` instances by `Rank` descending
   - High-rank files appear first in output

### Graph Example

```
main.go
  └──> utils.go
  └──> config.go

server.go
  └──> utils.go

handler.go
  └──> utils.go
  └──> auth.go

logger.go
  (no imports from others)
```

**Calculated Ranks:**
- `utils.go`: In-Degree = 3 → Rank = 0.9 (Highest)
- `config.go`: In-Degree = 1 → Rank = 0.5
- `auth.go`: In-Degree = 1 → Rank = 0.5
- `logger.go`: In-Degree = 0 → Rank = 0.1 (Lowest)

### Ranking Algorithm (Pseudocode)

```pseudocode
function rankFiles(nodes, importGraph):
    for each file in importGraph:
        inDegree[file] = count of nodes pointing to file
    
    maxDegree = max(inDegree.values())
    
    for each file in nodes:
        nodes[file].Rank = inDegree[file] / maxDegree
    
    sort(nodes) by Rank descending
    return nodes
```

---

## Phase C: Output (XML/JSON Rendering)

### Goal
Generate a structured, token-budgeted representation of the codebase.

### Why Structured Output?
LLMs are excellent at parsing XML/JSON tags for structure. Providing explicit tags helps with:
- Accurate context extraction
- Clearer semantic boundaries
- Easier post-processing

### XML Output Format

```xml
<repomap>
  <file path="main.go" importance="high" rank="0.95">
    <definition>func main()</definition>
    <definition>func initDB()</definition>
  </file>
  
  <file path="pkg/auth/login.go" importance="medium" rank="0.60">
    <definition>struct LoginRequest</definition>
    <definition>func Authenticate(user, pass)</definition>
  </file>
  
  <file path="pkg/utils/logger.go" importance="low" rank="0.10">
    <definition>func Log(msg)</definition>
  </file>
</repomap>
```

### JSON Output Format

```json
{
  "repomap": {
    "files": [
      {
        "path": "main.go",
        "importance": "high",
        "rank": 0.95,
        "definitions": [
          "func main()",
          "func initDB()"
        ]
      },
      {
        "path": "pkg/auth/login.go",
        "importance": "medium",
        "rank": 0.60,
        "definitions": [
          "struct LoginRequest",
          "func Authenticate(user, pass)"
        ]
      }
    ]
  }
}
```

### Token Budget Enforcement

1. **Initialize** budget counter with `--max-tokens` value
2. **Process files** in rank order:
   - Calculate `TokenCount` for current file
   - If `remaining_budget >= TokenCount`: Include file, subtract from budget
   - Otherwise: Stop processing, include truncation notice
3. **Generate output** with only included files

### Importance Levels

Based on rank scores:
- **High:** Rank > 0.7
- **Medium:** Rank 0.3 - 0.7
- **Low:** Rank < 0.3

---

## Phase Integration Example

### Step-by-Step Flow

```
INPUT: /home/user/myproject (3,000 files, --max-tokens=5000)
  ↓
PHASE A (Discovery & Parsing)
  ├─ Walk filesystem: 3,000 files found
  ├─ Filter with .gitignore: 2,100 files remain
  ├─ Parse 2,100 files, extract signatures
  └─ Create FileNode for each
  ↓
PHASE B (Ranking)
  ├─ Build import graph from FileNodes
  ├─ Calculate in-degree for each file
  ├─ Normalize to rank scores (0.0 - 1.0)
  └─ Sort by rank descending
  ↓
PHASE C (Output Rendering)
  ├─ Process files in rank order
  ├─ Include high-rank files first: utils.go, config.go, types.go (token count: 800)
  ├─ Include medium-rank files: auth.go, handler.go (token count: 1,200)
  ├─ Include low-rank files until budget exhausted (token count: 3,000)
  ├─ Remaining budget: 5000 - 3000 = 2000 tokens
  └─ OUTPUT: XML with 47 files, token count: 3,000/5,000
```

