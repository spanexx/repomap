# Repomap

Repomap is a CLI tool that generates a token-optimized “map” of a repository for use as LLM context. It focuses on **high-signal structure** (top-level definitions + imports), and then ranks files so the output stays useful even with a strict token budget.

Start here:

- [`README.md`](./README.md) for install + quick usage
- `repomap/USAGE.md` for detailed flags and examples

## What Repomap outputs

Repomap emits a structured representation of your repo:

- **XML** (default): best for pasting into prompts
- **JSON**: best for programmatic consumption
- **Text**: quick human skim

The output includes:

- file paths
- extracted definitions (functions/types/classes)
- extracted imports
- importance/rank metadata

## Why this is useful

Large repos don’t fit into a prompt. Repomap acts as a funnel:

- it **reduces** the codebase into a compact, structured summary
- it preserves the “shape” of the system (entry points, core packages, dependencies)
- it allows an agent to navigate and ask targeted follow-up questions

## High-level pipeline

1. **Discovery**
   - walk the directory tree
2. **Filtering**
   - respect `.gitignore`
   - drop binaries / irrelevant files
3. **Parsing**
   - extract top-level definitions + imports
4. **Ranking**
   - score files by centrality (dependency structure)
5. **Rendering**
   - emit XML/JSON/text
   - apply `--max-tokens` budget

## Typical workflow

1. Run Repomap and save output:

```bash
repomap --root . --output xml --max-tokens 8000 > repomap.xml
```

2. Paste `repomap.xml` into your agent prompt, and then ask for specific files to be opened.

3. (Optional) Run in server mode for a UI/API:

```bash
REPOMAP_PROVIDER=qodercli repomap --serve --port 8080
REPOMAP_PROVIDER=gemini-shell repomap --serve --port 8080
REPOMAP_PROVIDER=qwen-cli repomap --serve --port 8080
REPOMAP_PROVIDER=claude-cli repomap --serve --port 8080

```

Open `http://localhost:8080`.