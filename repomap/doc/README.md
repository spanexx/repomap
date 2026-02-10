# Repomap Documentation

Repomap is a CLI tool designed to compress Go repositories into a token-optimized format (XML or JSON) for Large Language Models (LLMs). It analyzes your codebase to understand file importance based on dependency structure and generates a concise map of definitions.

## Installation

### From Source

```bash
git clone https://github.com/spanexx/agents-cli.git
cd agents-cli/repomap
go build -o repomap ./cmd/repomap
```

## Usage

Run the `repomap` binary from your terminal:

```bash
./repomap --root /path/to/your/repo
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `--root` | Path to the repository root directory. | `.` (Current directory) |
| `--output` | Output format (`xml` or `json`). | `xml` |
| `--max-tokens` | Maximum token budget. Output will truncated if exceeded. | `0` (Unlimited) |
| `--include-ext` | Comma-separated list of file extensions to include. | `.go` |
| `--exclude-ext` | Comma-separated list of file extensions to exclude. | (None) |
| `--ignore-tests` | If set, ignores `*_test.go` files. | `false` |
| `--verbose` | Enable verbose logging to stderr. | `false` |
| `--version` | Show version information. | `false` |

## Examples

**Generate a map for the current directory:**
```bash
repomap
```

**Generate a JSON map for a specific project:**
```bash
repomap --root ~/projects/my-app --output json
```

**Generate a map with a 2000 token limit:**
```bash
repomap --max-tokens 2000
```

**Ignore test files and verbose logging:**
```bash
repomap --ignore-tests --verbose
```

## How it Works

1.  **Discovery**: Traverses the directory, respecting `.gitignore` rules.
2.  **Parsing**: Extracts top-level definitions (functions, types, interfaces) and imports from Go files using AST parsing.
3.  **Graph Construction**: Builds a dependency graph based on imports.
4.  **Ranking**: Ranks files using PageRank-like algorithm (In-Degree Centrality) to determine importance.
5.  **Output**: Renders the most important files first, truncating less important ones if the token budget is exceeded.
