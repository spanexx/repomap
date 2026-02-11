# Repomap

Repomap is a CLI tool that generates a token-optimized “map” of a repository (XML/JSON/text) for use as LLM context. It extracts top-level definitions and imports, ranks files by importance, and renders a compact summary that fits within a token budget.

## Features

- **Token-budgeted output** via `--max-tokens`
- **Multiple output formats**: `xml` (default), `json`, `text`
- **Language support**
  - Go: AST-based definitions + imports extraction
  - Other languages: basic line-based fallback extractor (best-effort)
- **Visualizer / server mode**: `--serve` to start a local UI/API

For a higher-level overview and architecture notes, see [`INDEX.md`](./INDEX.md).

## Install

### From source (Go)

```bash
git clone https://github.com/spanexx/agents-cli.git
cd agents-cli/repomap/repomap

go build -o repomap ./cmd/repomap
./repomap --version
```

### From npm

```bash
npm install -g @spanexx/repomap
repomap --root .
```

Or:

```bash
npx @spanexx/repomap --root .
```

## Usage

Basic:

```bash
repomap --root .
```

Common options:

```bash
repomap --root . --output xml
repomap --root . --output json --max-tokens 4000
repomap --include-ext .go,.md --ignore-tests
```

Server mode:

```bash
repomap --serve --port 8080
```

For a fuller CLI guide, see `repomap/USAGE.md`.

## Using Repomap in a new codebase (public install)

If you installed Repomap via npm, you can run it in any repository without cloning this repo.

1. **Go to the repo you want to map**

```bash
cd /path/to/your/project
```

2. **Run Repomap**

Using `npx` (no global install required):

```bash
npx @spanexx/repomap --root .
```

Or if installed globally:

```bash
repomap --root .
```

3. **Choose an output format and save it**

```bash
repomap --root . --output xml --max-tokens 8000 > repomap.xml
repomap --root . --output json --max-tokens 8000 > repomap.json
```

4. **(Optional) Start the Visualizer UI**

```bash
REPOMAP_PROVIDER=qodercli repomap --serve --port 8080
```

Then open `http://localhost:8080`.

## Configuration

Repomap supports configuration via:

- **CLI flags** (highest priority)
- **Environment variables** prefixed with `REPOMAP_`
- **Config file** (`.repomaprc` as JSON)

Example `.repomaprc`:

```json
{
  "output": "json",
  "ignore-tests": true,
  "max-tokens": 8000,
  "include-ext": ".go,.ts,.tsx"
}
```

## Provider credentials (important)

This project uses environment variables for provider credentials.

- **Never commit secrets** (API keys, OAuth client secrets) to git.
- If you previously committed a secret, rotate it and rewrite git history before making the repo public.

## License

GPL-3.0. See `LICENSE`.
