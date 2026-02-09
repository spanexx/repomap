# Recommended Go Libraries

This guide covers the industry-standard libraries used in the repomap project. Don't reinvent the wheel—these libraries are battle-tested and optimized.

## 1. godirwalk - File System Walking

**Package:** `github.com/karrick/godirwalk`

### Why It's Needed
Walking large directory trees efficiently is critical for performance. The standard library's `filepath.WalkDir` is adequate but slower than specialized solutions for massive repositories.

### Benefits
- Significantly faster than `filepath.WalkDir` for large file systems
- Lower memory footprint
- Built-in support for symlink handling
- Optimized for real-world codebases

### Basic Usage
```go
import "github.com/karrick/godirwalk"

err := godirwalk.Walk("./src", &godirwalk.Options{
    Callback: func(path string, entry *godirwalk.Dirent) error {
        // Process each file
        return nil
    },
})
```

---

## 2. go-gitignore - Git Ignore Parsing

**Package:** `github.com/monochromegane/go-gitignore`

### Why It's Needed
Crucial for excluding patterns defined in `.gitignore`. Without this, your repomap would include `node_modules`, `vendor`, build artifacts, and other irrelevant files.

### Benefits
- Reads and parses standard `.gitignore` files
- Efficiently matches paths against ignore patterns
- Supports nested `.gitignore` files in subdirectories
- Respects standard gitignore semantics

### Basic Usage
```go
import "github.com/monochromegane/go-gitignore"

gitignore := gitignore.New(excludeFile)
if gitignore.Match(filePath, false) {
    // This file should be ignored
}
```

---

## 3. go-tree-sitter - Parsing and AST Extraction

**Package:** `github.com/smacker/go-tree-sitter`

### Why It's Needed
The **core** of the parsing phase. Tree-sitter provides robust parsing for multiple languages without writing language-specific regex patterns. This is the key to supporting Python, JavaScript, Go, Rust, and other languages.

### Benefits
- Language-agnostic parsing through unified interface
- Robust to syntax errors (incremental parsing)
- Extracts complete Abstract Syntax Trees (ASTs)
- Significantly more reliable than regex-based parsing
- Community-maintained parsers for most popular languages

### Supported Languages
- Go, Python, JavaScript, TypeScript, Rust, C, C++, Java, and more

### Basic Usage
```go
import "github.com/smacker/go-tree-sitter/golang"
import "github.com/smacker/go-tree-sitter"

parser := sitter.NewParser()
parser.SetLanguage(golang.GetLanguage())

tree, _ := parser.ParseCtx(ctx, nil, sourceCode)
root := tree.RootNode()

// Traverse AST to extract definitions
```

### Note for MVP
For the initial implementation, you can use Go's standard library `go/ast` package instead. Switch to Tree-sitter when you need to support multiple languages.

---

## 4. tiktoken-go - Token Counting

**Package:** `github.com/pkoukk/tiktoken-go`

### Why It's Needed
To accurately respect the `--max-tokens` budget. Since the output needs to fit within LLM context windows, precise token counting is essential.

### Benefits
- Accurate token counting matching OpenAI's tokenization
- Supports multiple encoding formats (cl100k_base, p50k_base, etc.)
- Lightweight and fast
- Eliminates guesswork about token budgets

### Basic Usage
```go
import "github.com/pkoukk/tiktoken-go"

enc, _ := tiktoken.GetEncoding("cl100k_base")
tokens := enc.Encode(text, nil, nil)
tokenCount := len(tokens)
```

---

## Installation

Add these libraries to your `go.mod` file:

```bash
go get github.com/karrick/godirwalk
go get github.com/monochromegane/go-gitignore
go get github.com/smacker/go-tree-sitter
go get github.com/pkoukk/tiktoken-go
```

---

## Library Decision Matrix

| Task | Library | Why | Alternative |
|------|---------|-----|-------------|
| File Walking | godirwalk | Fast, efficient | filepath.WalkDir (slower) |
| Gitignore | go-gitignore | Standard parsing | Manual pattern matching |
| Parsing | go-tree-sitter | Language-agnostic | go/ast (Go-only) |
| Token Counting | tiktoken-go | Accurate | Approximation/estimation |

