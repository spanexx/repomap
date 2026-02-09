# Phase 2: Multi-Language Support & Semantic Tools

## Objective

Extend agents-cli beyond Go to support all major programming languages, and add semantic understanding capabilities through embedding-based code search.

## Scope

### In-Scope for Phase 2

1. **Multi-Language Parsing** (Milestone 2.1)
   - Tree-sitter integration for Go, Python, JavaScript, TypeScript, Rust, Java, C++
   - Update repomap to support all languages
   - Language-specific ranking strategies
   - Extend CLI framework to handle language selection

2. **Semantic Search Tool** (Milestone 2.2)
   - New tool: `semgrep-cli` – semantic grep using embeddings
   - Local lightweight embedding model (all-MiniLM-L6-v2)
   - Vector similarity search for code snippets
   - Natural language queries: "where is authentication handled?"
   - Integration with repomap output

3. **Web Navigation Foundation** (Milestone 2.3)
   - New tool: `web-ray` – structured web automation
   - Headless browser with Playwright
   - Accessibility tree extraction
   - Element coordinate detection
   - JavaScript execution support

### Out-of-Scope for Phase 2

- Cloud deployment (deferred to Phase 4)
- Authentication/authorization (deferred to Phase 4)
- Advanced memory management (deferred to Phase 3)
- Sandboxed execution (deferred to Phase 3)

## Milestones

### [Milestone 2.1: Tree-Sitter Integration](Milestone-2.1-Tree-Sitter-Integration/)
Integrate Tree-sitter for multi-language parsing and update repomap to support all major languages.

### [Milestone 2.2: Semgrep-CLI](Milestone-2.2-Semgrep-CLI/)
Build semantic code search using embeddings and natural language queries.

### [Milestone 2.3: Web-Ray Foundation](Milestone-2.3-Web-Ray-Foundation/)
Create structured web automation tool for agents to navigate and extract data from websites.

## Success Criteria

- ✅ Repomap supports Go, Python, JavaScript, TypeScript, Rust, Java, C++
- ✅ Multi-language ranking strategies work correctly
- ✅ semgrep-cli finds relevant code using semantic search
- ✅ web-ray provides structured access tree for any website
- ✅ All tools use CLI framework without deviation
- ✅ >80% code coverage on new packages
- ✅ Performance: Embedding search <500ms per query
- ✅ Documentation complete for all tools

## Dependencies

- **Internal:** Phase 1 (Foundation & MVP) – must be complete
- **External:** Tree-sitter, Playwright, sentence-transformers or similar

## Timeline Estimate

- **Milestone 2.1:** 3–4 weeks (Tree-sitter integration + repomap upgrade)
- **Milestone 2.2:** 2–3 weeks (Semantic search implementation)
- **Milestone 2.3:** 2–3 weeks (Web automation)
- **Total Phase 2:** ~7–10 weeks

## Handoff to Phase 3

Upon completion of Phase 2, Phase 3 will:
- Build sandboxed execution (sandbox-run)
- Create API schema detection (api-forge)
- Implement git merge conflict solver (git-surgeon)
- Add long-term memory store (mem-kv)
- Build math co-processor (calc-bridge)

---

## Navigation

- [Milestone 2.1: Tree-Sitter Integration](Milestone-2.1-Tree-Sitter-Integration/README.md)
- [Milestone 2.2: Semgrep-CLI](Milestone-2.2-Semgrep-CLI/README.md)
- [Milestone 2.3: Web-Ray Foundation](Milestone-2.3-Web-Ray-Foundation/README.md)
