# Commands Reference

This is a quick reference for all repomap CLI commands with complete syntax, flags, and examples.

---

## Command: repomap

The main command for analyzing and mapping repositories.

### Complete Syntax

```bash
repomap [GLOBAL-FLAGS] [OPTIONS]
```

### Usage Patterns

#### Pattern 1: Analyze Current Directory

```bash
repomap
```

**What it does:**
- Analyzes the current directory
- Respects `.gitignore` rules
- Outputs XML format by default
- Uses 5000-token budget

**Expected output:**
```xml
<repomap>
  <file path="..." importance="..." rank="...">
    ...
  </file>
</repomap>
```

---

#### Pattern 2: Analyze Specific Directory

```bash
repomap --root /path/to/repo
```

**Example:**
```bash
repomap --root ./myapp
repomap --root /home/user/projects/agent-cli
```

**What it does:**
- Analyzes the specified directory
- Works with absolute or relative paths
- Respects `.gitignore` in that directory

---

#### Pattern 3: JSON Output

```bash
repomap --output json
```

**Expected output:**
```json
{
  "repomap": {
    "files": [
      {
        "path": "main.go",
        "importance": "high",
        "rank": 0.95,
        "definitions": [...]
      }
    ]
  }
}
```

**Common use cases:**
- Processing with `jq` or other JSON tools
- Integration with other tools
- Programmatic analysis

---

#### Pattern 4: XML Output

```bash
repomap --output xml
```

**Expected output:**
```xml
<repomap>
  <file path="main.go" importance="high" rank="0.95">
    <definition>func main()</definition>
  </file>
</repomap>
```

**Common use cases:**
- Human-readable format
- Integration with XML parsers
- Default output format

---

### All Available Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--root` | `<PATH>` | `.` | Repository root directory |
| `--output` | `xml\|json` | `xml` | Output format |
| `--max-tokens` | `<NUMBER>` | `5000` | Token budget |
| `--include-ext` | `<EXTS>` | all | File extensions to include |
| `--exclude-ext` | `<EXTS>` | none | File extensions to exclude |
| `--ignore-tests` | flag | false | Exclude test files |
| `--ignore-vendor` | flag | true | Exclude vendor dirs |
| `--verbose` | flag | false | Debug output |
| `--follow-symlinks` | flag | false | Follow symbolic links |
| `--gitignore-file` | `<PATH>` | `.gitignore` | Custom ignore file |
| `--include-lang` | `<LANGS>` | all | Languages to include |
| `--exclude-lang` | `<LANGS>` | none | Languages to exclude |
| `--max-files` | `<NUMBER>` | unlimited | Max files to process |
| `--parallel` | `<NUMBER>` | auto | Parallel workers |
| `--help` / `-h` | flag | - | Show help |
| `--version` / `-v` | flag | - | Show version |

---

## Real Command Examples

### Example 1: Basic Go Project Analysis

```bash
repomap --root ./mygoapp
```

**Output snippet:**
```xml
<repomap>
  <file path="main.go" importance="high" rank="0.95">
    <definition>func main()</definition>
    <definition>func init()</definition>
  </file>
  <file path="server/server.go" importance="high" rank="0.90">
    <definition>type Server struct</definition>
    <definition>func NewServer(config *Config) *Server</definition>
    <definition>func (s *Server) Start() error</definition>
  </file>
</repomap>
```

---

### Example 2: Strict Token Budget for LLM Context

```bash
repomap --root . --max-tokens 1000 --output json
```

**Behavior:**
- Includes only the highest-ranked files
- Typically 3-8 files with 1000-token budget
- Ideal for small context windows

**Output snippet:**
```json
{
  "repomap": {
    "files": [
      {
        "path": "utils.go",
        "importance": "high",
        "rank": 0.95,
        "definitions": ["func Helper()"]
      },
      {
        "path": "types.go",
        "importance": "high",
        "rank": 0.92,
        "definitions": ["type Config struct"]
      }
    ]
  }
}
```

---

### Example 3: Go Files Only

```bash
repomap --root . --include-ext .go --output xml
```

**Effect:**
- Excludes `.md`, `.txt`, `.json`, `.yaml` files
- Only analyzes `.go` source files
- Speeds up processing

**Output:**
Only `.go` files appear in the repomap.

---

### Example 4: Language Filtering (Python Project)

```bash
repomap --root ./pyproject --include-lang python --output json
```

**Effect:**
- Only Python files are analyzed
- Excludes requirements.txt, setup.py (non-code files)
- Optimized for Python-specific parsing

---

### Example 5: Exclude Tests and Vendors

```bash
repomap --root . --ignore-tests --ignore-vendor --output xml
```

**Effect:**
- Removes all `*_test.go`, `test_*.go` files
- Excludes `vendor/` and `node_modules/` directories
- Typical reduction: 20-40% fewer files

**Before:**
```
178 files analyzed
```

**After:**
```
120 files analyzed (tests and vendors removed)
```

---

### Example 6: Large Budget for Comprehensive View

```bash
repomap --root . --max-tokens 10000 --output json
```

**Effect:**
- Includes almost all files (unless very large repo)
- Provides complete picture of codebase
- Useful for thorough analysis

---

### Example 7: Verbose Mode for Debugging

```bash
repomap --root . --verbose --output xml 2> debug.log
```

**Stderr output (debug.log):**
```
[INFO] Starting discovery phase...
[INFO] Found 256 files
[INFO] Applying .gitignore filters...
[INFO] 178 files after filtering
[INFO] Parsing Go files...
[INFO] Extracted definitions from 178 files
[INFO] Building import graph with 178 nodes
[INFO] Calculating rankings...
[INFO] Top 5 files by rank:
  [0.95] utils.go (imported by 15 files)
  [0.92] types.go (imported by 12 files)
  [0.88] config.go (imported by 9 files)
  [0.75] handler.go (imported by 5 files)
  [0.65] logger.go (imported by 3 files)
[INFO] Rendering output...
[INFO] Output: 47 files in 2150 tokens
[INFO] Done!
```

---

### Example 8: Save to File

```bash
repomap --root . --output json --max-tokens 5000 > repomap.json
```

**Result:** Creates `repomap.json` with full analysis.

---

### Example 9: Compare Multiple Budgets

```bash
# Tight budget
repomap --root . --max-tokens 500 --output json > tight.json
echo "Tight budget: $(jq '.repomap.files | length' tight.json) files"

# Medium budget
repomap --root . --max-tokens 2000 --output json > medium.json
echo "Medium budget: $(jq '.repomap.files | length' medium.json) files"

# Generous budget
repomap --root . --max-tokens 5000 --output json > generous.json
echo "Generous budget: $(jq '.repomap.files | length' generous.json) files"
```

**Output:**
```
Tight budget: 4 files
Medium budget: 12 files
Generous budget: 35 files
```

---

### Example 10: Integration with jq for Analysis

```bash
# Count total files
repomap --root . --output json | jq '.repomap.files | length'
# Output: 47

# Extract all file paths
repomap --root . --output json | jq -r '.repomap.files[].path'
# Output:
# main.go
# server/server.go
# utils/helper.go
# ...

# List only high-importance files
repomap --root . --output json | jq '.repomap.files[] | select(.importance == "high") | .path'
# Output:
# utils.go
# types.go
# config.go

# Get total rank score
repomap --root . --output json | jq '[.repomap.files[].rank] | add'
# Output: 34.5
```

---

### Example 11: Custom .gitignore File

```bash
repomap --root /external/repo --gitignore-file /external/repo/.custom-ignore --output xml
```

**Effect:**
- Uses `/external/repo/.custom-ignore` instead of `.gitignore`
- Useful for external repositories or custom filters

---

### Example 12: Processing with xargs

```bash
repomap --root . --output json | \
  jq -r '.repomap.files[] | select(.importance == "high") | .path' | \
  xargs -I {} echo "Important file: {}"
```

**Output:**
```
Important file: main.go
Important file: server/server.go
Important file: utils/helper.go
```

---

## Expected Outputs by Command

### Basic Command

```bash
$ repomap
```

```xml
<repomap>
  <!-- Generated repomap in XML format -->
</repomap>
```

Exit code: `0`

---

### With Options

```bash
$ repomap --max-tokens 2000 --output json --verbose
```

Stdout:
```json
{ "repomap": { "files": [...] } }
```

Stderr:
```
[INFO] Processing complete: 45 files analyzed
```

Exit code: `0`

---

### Error Cases

```bash
$ repomap --root /nonexistent
```

Stderr:
```
Error: root directory not found: /nonexistent
```

Exit code: `1`

---

```bash
$ repomap --max-tokens invalid
```

Stderr:
```
Error: invalid token budget: expected number, got "invalid"
```

Exit code: `1`

