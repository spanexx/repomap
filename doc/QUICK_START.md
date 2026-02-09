# Quick Start Guide

Get up and running with repomap in minutes.

---

## Prerequisites

- **Go 1.19 or later** (for Go standard library `go/ast` support)
- **Git** (to manage the project and handle `.gitignore`)

---

## Installation

### Clone the Repository

```bash
git clone https://github.com/yourusername/agents-cli.git
cd agents-cli/repomap
```

### Install Dependencies

```bash
go mod download
```

### Build the Binary

```bash
go build -o repomap ./cmd/repomap
```

### Verify Installation

```bash
./repomap --version
```

Expected output:
```
repomap version 1.0.0
```

---

## Your First Run

### Basic Usage

Map a Go project in the current directory:

```bash
./repomap --root . --output xml
```

This generates an XML representation of your codebase and prints it to stdout.

### Save Output to File

```bash
./repomap --root . --output xml > repomap.xml
```

### Limit Output Size

Map with a 2000-token budget:

```bash
./repomap --root . --max-tokens 2000 --output json
```

---

## Common Use Cases

### Map Your Current Project

```bash
repomap --root . --output xml
```

### Map a Specific Subdirectory

```bash
repomap --root ./src --output json
```

### Use JSON Format for Processing

```bash
repomap --root . --output json | jq '.repomap.files | length'
```

### Strict Budget for Small Context

```bash
repomap --root . --max-tokens 500 --output xml
```

### Focus on Specific File Types

```bash
repomap --root . --include-ext .go,.mod --output xml
```

### Exclude Test Files

```bash
repomap --root . --ignore-tests --output json
```

---

## Output Examples

### XML Format

```xml
<repomap>
  <file path="main.go" importance="high" rank="0.95">
    <definition>func main()</definition>
    <definition>func init()</definition>
  </file>
  <file path="pkg/utils/helper.go" importance="medium" rank="0.65">
    <definition>func Helper(input string) string</definition>
    <definition>type Config struct</definition>
  </file>
</repomap>
```

### JSON Format

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
          "func init()"
        ]
      },
      {
        "path": "pkg/utils/helper.go",
        "importance": "medium",
        "rank": 0.65,
        "definitions": [
          "func Helper(input string) string",
          "type Config struct"
        ]
      }
    ]
  }
}
```

---

## Understanding the Output

### Importance Levels

- **high:** Rank > 0.7 (frequently imported, core files)
- **medium:** Rank 0.3 - 0.7 (moderately important)
- **low:** Rank < 0.3 (utility or leaf files)

### Rank Score

The rank (0.0 - 1.0) is based on how many other files import this file:
- Higher rank = more dependencies on this file
- Lower rank = fewer dependencies (leaf node)

---

## Troubleshooting

### "No files found"

- Verify the `--root` path exists and contains Go files
- Check that files aren't being filtered by `.gitignore`
- Try `repomap --root . --verbose` for debug output

### "Token budget exceeded"

- Increase `--max-tokens` value
- Use `--ignore-tests` to reduce output size

### Command Not Found

- Verify the binary is in your PATH or use `./repomap` from the build directory
- Rebuild with `go build -o repomap ./cmd/repomap`

---

## Next Steps

For more detailed information, see:
- [ARCHITECTURE.md](ARCHITECTURE.md) - Design overview
- [CLI_REFERENCE.md](CLI_REFERENCE.md) - Complete command documentation
- [EXAMPLES.md](EXAMPLES.md) - Real-world examples
