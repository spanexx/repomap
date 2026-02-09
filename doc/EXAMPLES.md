# Real-World Examples

This document showcases repomap with different scenarios and expected outputs.

---

## Example 1: Simple Web Server (Go)

### Repository Structure

```
webserver/
├── main.go
├── go.mod
├── config/
│   └── config.go
├── server/
│   ├── server.go
│   └── handler.go
├── utils/
│   └── logger.go
└── tests/
    ├── server_test.go
    └── handler_test.go
```

### Command

```bash
repomap --root ./webserver --max-tokens 2000 --output xml
```

### Expected Output

```xml
<repomap>
  <file path="server/server.go" importance="high" rank="0.90">
    <definition>type Server struct</definition>
    <definition>func NewServer(config *Config) *Server</definition>
    <definition>func (s *Server) Start(port int) error</definition>
    <definition>func (s *Server) Shutdown() error</definition>
  </file>
  
  <file path="server/handler.go" importance="high" rank="0.85">
    <definition>func HandleRequest(w http.ResponseWriter, r *http.Request)</definition>
    <definition>func HandleHealth(w http.ResponseWriter, r *http.Request)</definition>
  </file>
  
  <file path="config/config.go" importance="medium" rank="0.70">
    <definition>type Config struct</definition>
    <definition>func LoadConfig(path string) (*Config, error)</definition>
  </file>
  
  <file path="main.go" importance="medium" rank="0.65">
    <definition>func main()</definition>
    <definition>func init()</definition>
  </file>
  
  <file path="utils/logger.go" importance="low" rank="0.25">
    <definition>func Log(msg string)</definition>
  </file>
</repomap>
```

### Analysis

- **server.go** has the highest rank because it's imported by handler.go and main.go
- **handler.go** is also highly ranked (imported by main.go)
- **config.go** is moderately important (imported by main.go and server.go)
- **logger.go** has low importance (imported by few files)
- Test files are excluded by default

---

## Example 2: Microservices Architecture

### Repository Structure

```
microservices/
├── services/
│   ├── auth/
│   │   ├── main.go
│   │   ├── handler.go
│   │   └── db.go
│   ├── user/
│   │   ├── main.go
│   │   └── handler.go
│   └── payment/
│       ├── main.go
│       └── handler.go
├── shared/
│   ├── models.go
│   ├── errors.go
│   └── middleware.go
└── go.mod
```

### Command

```bash
repomap --root ./microservices --max-tokens 3500 --output json
```

### Expected Output (JSON)

```json
{
  "repomap": {
    "files": [
      {
        "path": "shared/models.go",
        "importance": "high",
        "rank": 0.95,
        "definitions": [
          "type User struct",
          "type Payment struct",
          "type AuthToken struct"
        ]
      },
      {
        "path": "shared/middleware.go",
        "importance": "high",
        "rank": 0.90,
        "definitions": [
          "func AuthMiddleware(next http.Handler) http.Handler",
          "func LoggingMiddleware(next http.Handler) http.Handler"
        ]
      },
      {
        "path": "services/auth/handler.go",
        "importance": "medium",
        "rank": 0.75,
        "definitions": [
          "func HandleLogin(w http.ResponseWriter, r *http.Request)",
          "func HandleRegister(w http.ResponseWriter, r *http.Request)"
        ]
      },
      {
        "path": "services/auth/db.go",
        "importance": "medium",
        "rank": 0.70,
        "definitions": [
          "type AuthDB struct",
          "func (db *AuthDB) GetUser(id string) (*User, error)"
        ]
      },
      {
        "path": "services/user/handler.go",
        "importance": "medium",
        "rank": 0.65,
        "definitions": [
          "func HandleGetUser(w http.ResponseWriter, r *http.Request)"
        ]
      }
    ]
  }
}
```

### Analysis

- **shared/models.go** is highest priority (imported by all services)
- **shared/middleware.go** is also essential (shared utilities)
- Service-specific handlers follow in importance
- Database layers and main.go are excluded due to token budget

---

## Example 3: Token Budget Comparison

### Command Set

```bash
# Tight budget: 500 tokens
repomap --root . --max-tokens 500 --output xml > repomap-tight.xml

# Medium budget: 2000 tokens
repomap --root . --max-tokens 2000 --output xml > repomap-medium.xml

# Generous budget: 5000 tokens
repomap --root . --max-tokens 5000 --output xml > repomap-generous.xml
```

### Output Comparison

| Budget | Files Included | Key Characteristics |
|--------|----------------|-------------------|
| 500 tokens | 3-5 files | Only core dependencies (models, shared utilities) |
| 2000 tokens | 8-12 files | Core + important services/handlers |
| 5000 tokens | 20-30 files | Nearly complete picture with utilities and tests |

### 500-Token Output Example

```xml
<repomap>
  <file path="shared/models.go" importance="high" rank="0.95">
    <definition>type User struct</definition>
    <definition>type Config struct</definition>
  </file>
  
  <file path="shared/errors.go" importance="high" rank="0.90">
    <definition>type Error struct</definition>
    <definition>func NewError(code int) *Error</definition>
  </file>
  
  <file path="services/auth/handler.go" importance="medium" rank="0.70">
    <definition>func HandleLogin(w, r)</definition>
  </file>
</repomap>
<!-- Budget exhausted: 485/500 tokens used. -->
```

---

## Example 4: Language Filtering

### Command

```bash
repomap --root . --include-ext .go,.mod --output xml
```

### Before Filtering

```
Total files: 150
├── Go files: 120
├── Test files: 15
├── Config files (JSON, YAML): 10
└── Documentation: 5
```

### After Filtering

```
Total files: 122 (120 Go + 2 go.mod)
├── Go files: 120
└── Module files: 2
```

### Output

Only Go source files and module definitions are included in the repomap.

---

## Example 5: Excluding Test Files

### Command

```bash
repomap --root . --ignore-tests --output json
```

### Effect

All files matching these patterns are excluded:
- `*_test.go`
- `test_*.go`
- `*_suite_test.go`
- Any file in a `test/` or `tests/` directory

### Size Reduction

Typically reduces output by 20-30% for projects with extensive test suites.

---

## Example 6: Verbose Output for Debugging

### Command

```bash
repomap --root . --verbose --output xml
```

### Additional Output (to stderr)

```
[INFO] Starting discovery phase...
[INFO] Found 256 files
[INFO] Applying .gitignore filters...
[INFO] 178 files after filtering
[INFO] Parsing Go files...
[INFO] Extracted definitions from 178 files
[INFO] Building import graph...
[INFO] Graph contains 178 nodes and 342 edges
[INFO] Calculating rankings...
[INFO] Top files by rank:
  - utils.go: 0.95
  - config.go: 0.88
  - types.go: 0.82
[INFO] Rendering output...
[INFO] Output: 125 files in 1980 tokens
[INFO] Done!
```

---

## Example 7: Piping to Other Tools

### Count Files in Output

```bash
repomap --root . --output json | jq '.repomap.files | length'
```

Output: `47`

### Extract File Paths

```bash
repomap --root . --output json | jq -r '.repomap.files[].path'
```

Output:
```
main.go
config/config.go
server/server.go
utils/logger.go
...
```

### Filter by Importance Level

```bash
repomap --root . --output json | jq '.repomap.files[] | select(.importance == "high") | .path'
```

Output:
```
utils.go
config.go
types.go
```

