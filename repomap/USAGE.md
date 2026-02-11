# Repomap Usage Guide

`repomap` is a CLI tool designed to generate a token-optimized map of your repository (supporting Go, TypeScript, Python, and more), making it easier for AI agents to understand your codebase structure.

## Installation

### From Source
```bash
git clone https://github.com/spanexx/agents-cli.git
cd agents-cli/repomap/cmd/repomap
go install
```
Ensure `$GOPATH/bin` is in your `$PATH`.

Ensure `$GOPATH/bin` is in your `$PATH`.

### From NPM
Install globally via npm:
```bash
npm install -g @spanexx/repomap
```
or run directly with `npx`:
```bash
npx @spanexx/repomap --root .
```
*(No Go installation required)*

### Running Locally (Development)
If you prefer not to install it globally, you can build and run it locally:
```bash
# Build the binary
go build -o repomap cmd/repomap/main.go

# Run it
./repomap --version
```
Or run directly with `go run`:
```bash
go run cmd/repomap/main.go --root .
```

## Basic Usage

Run `repomap` in the root of your repository:

```bash
repomap
```
This will output an XML representation of your Go files to `stdout`.

### Common Options

-   **`--root <path>`**: Specify the root directory of the repository (default: `.`).
-   **`--output <format>`**: Choose output format: `xml` (default), `json`, or `text`.
-   **`--max-tokens <int>`**: Set a token budget. Repomap will prioritize important files to fit within this limit.
-   **`--plan <path>`**: Load an architectural plan (JSON) to merge with the map.
-   **`--analyze`**: Run static analysis to detect duplication, intent violations, and circular dependencies.

## Advanced Filtering

Repomap allows you to control which files are included in the map:

-   **`--include-ext <exts>`**: Comma-separated list of file extensions to include (default: `.go`).
    ```bash
    repomap --include-ext .go,.md
    ```
-   **`--exclude-ext <exts>`**: Comma-separated list of extensions to exclude.
-   **`--ignore-tests`**: Skip test files (`*_test.go`).
    ```bash
    repomap --ignore-tests
    ```

## Agent Mode & Visualizer

Repomap includes a built-in server for visualizing the plan and interacting with agents.

```bash
repomap --serve --port 9090
```
This starts a local server at `http://localhost:9090`.

### Provider Configuration

When running in `--serve` mode, you can configure which AI provider to use for the chat interface.
This is controlled via the `REPOMAP_PROVIDER` environment variable.

Supported providers:
-   `gemini-cli` (default)
-   `qwen-cli`
-   `qodercli`

```bash
# Use Qwen instead of Gemini
export REPOMAP_PROVIDER=qwen-cli
repomap --serve
```

You can also specify a fallback chain:
```bash
# Try Gemini first, fallback to Qwen
export REPOMAP_PROVIDER=gemini-cli,qwen-cli
repomap --serve
```

**Note:** Ensure you have the necessary authentication set up for the chosen provider (e.g., `gcloud auth login` for Gemini, or config files for others).

## Output Formats

Repomap supports multiple output formats to suit different needs:

-   **XML (`--output xml`)**: (Default) Best for LLM context. Structured, tag-based format containing file paths, definitions, and imports.
-   **JSON (`--output json`)**: Ideal for programmatic processing. Contains the same rich data as XML but in JSON format.
-   **Text (`--output text`)**: A simple, indented tree-like view of the repository structure. Good for quick human inspection.

## Configuration

You can configure `repomap` using a configuration file or environment variables. The precedence order is:
1.  **Command-line Flags** (Highest)
2.  **Environment Variables**
3.  **Configuration File**
4.  **Defaults** (Lowest)

### Configuration File
Create a `.repomaprc` file in your home directory (`$HOME/.repomaprc`) or the repository root:
```json
{
  "output": "json",
  "ignore-tests": true,
  "max-tokens": 8000,
  "include-ext": ".go,.ts,.tsx",
  "exclude-ext": "node_modules,dist"
}
```

### Environment Variables
All flags can be set via environment variables prefixed with `REPOMAP_`:

```bash
export REPOMAP_OUTPUT=json
export REPOMAP_MAX_TOKENS=8000
export REPOMAP_VERBOSE=true
repomap
```

## Troubleshooting

### Common Issues

**1. "Server failed: address already in use"**
The default port `8080` might be occupied. Try a different port:
```bash
repomap --serve --port 9091
```

**2. "Discovery failed: permission denied"**
Ensure you have read permissions for the directory you are scanning. You can exclude problematic directories:
```bash
repomap --exclude-ext node_modules,.git
```

**3. "Warning: failed to create gemini provider"**
This usually means authentication credentials are missing.
-   For **Gemini**: Run `gcloud auth application-default login` or set `GOOGLE_APPLICATION_CREDENTIALS`.
-   For **Qwen/Qoder**: Ensure their respective config files or environment variables are set.

**4. "Rendering failed"**
If using a custom output format, ensure the writer implementation supports the data structure. The default `xml` and `json` formats are robust for standard Go repositories.

## Examples

**Generate a JSON map of the current directory, excluding tests:**
```bash
repomap --output json --ignore-tests
```

**Generate a map for a different project with a strict token limit:**
```bash
repomap --root ../other-project --max-tokens 4000
```

**Start the Visualizer UI:**
```bash
repomap --serve
```
