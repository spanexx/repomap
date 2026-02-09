# Milestone 2.2: Quick Reference

## Summary

Semantic code search tool using embeddings. Natural language queries find relevant code across multi-language repositories.

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 10 of Phase 2

## Example Usage

```bash
# Find authentication logic
semgrep-cli --root ./app "where is user authentication handled?"

# Find database connections
semgrep-cli --root ./backend "find database connection logic"

# Filter by language
semgrep-cli --root ./src "find error handling" --lang python

# Interactive mode
semgrep-cli --interactive
```

## Key Components

| Component | Path | Purpose |
|-----------|------|---------|
| Model Loader | `pkg/embeddings/model.go` | Load & manage embedding model |
| Vectorizer | `pkg/embeddings/vectorizer.go` | Convert text to vectors |
| Vector Index | `pkg/search/index.go` | Store & search vectors |
| Similarity | `pkg/search/cosine.go` | Vector similarity computation |
| Ranker | `pkg/search/ranker.go` | Multi-factor ranking |
| CLI Tool | `cmd/semgrep-cli/main.go` | Command-line interface |

## Execution Strategy

1. **Model & Infrastructure (Days 1–5)**
   - Evaluate embedding models
   - Set up vector infrastructure
   - Implement vectorization

2. **Search Engine (Days 6–12)**
   - Create index structure
   - Implement similarity search
   - Build ranking logic

3. **Interface & Polish (Days 13–21)**
   - CLI tool implementation
   - Output formatting
   - Testing & documentation

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| Query Latency | <500ms | — |
| Result Relevance | >90% | — |
| Code Coverage | >85% | — |
| Test Pass Rate | 100% | — |

## Acceptance Criteria

- ✅ Semantic queries work on diverse codebases
- ✅ Results ranked by relevance
- ✅ <500ms response time
- ✅ >90% relevance on benchmark queries
- ✅ >85% code coverage
- ✅ Interactive REPL mode
- ✅ Multi-repository search support

## Key Metrics

| Metric | Target |
|--------|--------|
| Query Response Time | <500ms |
| Index Build Time | <5s per 10k LOC |
| Index Memory | <100MB per 50k LOC |
| Relevance Score | >90% |
