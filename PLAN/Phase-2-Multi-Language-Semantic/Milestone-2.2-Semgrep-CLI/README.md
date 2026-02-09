# Milestone 2.2: Semgrep-CLI

## Objective

Build a semantic code search tool that uses embeddings and natural language queries to find relevant code snippets across repositories.

## Scope

### In-Scope

1. **Embedding Model**
   - Integrate lightweight embedding model (all-MiniLM-L6-v2)
   - Code snippet vectorization
   - Local inference (no external API calls)

2. **Semantic Search**
   - Index code definitions and snippets
   - Vector similarity search
   - Natural language query interface
   - Rank results by relevance

3. **Query Interface**
   - CLI: `semgrep-cli "where is user authentication?"`
   - Support for various query types
   - Result formatting (JSON/XML/text)

4. **Integration**
   - Works with repomap output
   - Extends CLI framework
   - Cross-language support

### Out-of-Scope

- Real-time indexing
- Advanced ML model fine-tuning
- Web UI (Phase 4+)

## Deliverables

1. **Embedding Service** (`pkg/embeddings/`)
2. **Search Index** (`pkg/search/`)
3. **semgrep-cli Tool** (`cmd/semgrep-cli/`)
4. **Integration Tests**
5. **Documentation & Examples**

## Success Criteria

- ✅ Finds relevant code for natural language queries
- ✅ Query response <500ms per search
- ✅ >90% relevance on test queries
- ✅ Works with all Phase 2.1 languages
- ✅ >80% code coverage
- ✅ Clear documentation with examples

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** ~20–25

---

## Example Usage

```bash
semgrep-cli --root ./myapp "where is user authentication handled?"
# Output: JSON with relevant code snippets, files, line numbers

semgrep-cli --root ./backend "find database connection logic"
# Output: Ranked results with vector similarity scores
```

## Dependencies

- Phase 1 completion
- Phase 2.1 (Tree-sitter) completion
- Embedding model library
