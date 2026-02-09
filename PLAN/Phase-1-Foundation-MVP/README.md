# Phase 1: Foundation & MVP

## Objective

Establish the foundational infrastructure, tooling, and architecture patterns for agents-cli. This phase implements the core `repomap` tool (the primary MVP) and a reusable CLI framework that future tools can build upon.

## Scope

### In-Scope for Phase 1

1. **Core `repomap` Tool**
   - File system discovery with `.gitignore` support
   - AST parsing for Go (using `go/ast` standard library)
   - Import graph construction and ranking
   - XML/JSON output with token budgeting
   - Configurable command-line interface

2. **CLI Framework & Shared Infrastructure**
   - Standard CLI flag/option parsing
   - Output formatting (XML, JSON)
   - Configuration and defaults management
   - Shared utilities for future tools

3. **Testing & Quality Assurance**
   - Unit tests for core parsing and ranking logic
   - Integration tests on real Go repositories
   - Performance benchmarks
   - CI/CD pipeline setup

4. **Documentation & Deployment**
   - Architecture documentation (in `repomap/doc/`)
   - User guides and quick-start
   - Build and installation instructions
   - Single-binary distribution

### Out-of-Scope for Phase 1

1. **Multi-language Support** (defer to Phase 2)
   - Tree-sitter integration
   - Python, JavaScript, Rust, etc. parsing
   - Language-specific ranking strategies

2. **Advanced Features** (defer to Phase 2+)
   - Semantic search (`semgrep-cli`)
   - Web navigation (`web-ray`)
   - Sandboxed execution (`sandbox-run`)
   - Other "Hands" and "Brain" tools

3. **Agent Integration** (defer to Phase 3+)
   - Integration with specific LLM agents
   - Custom MCP server protocols
   - Advanced memory/context management

4. **Distributed/Cloud Infrastructure** (defer to Phase 4+)
   - Server mode for remote access
   - Horizontal scaling
   - Cloud deployment templates

## Milestones & Dependencies

The milestones in Phase 1 follow a logical progression, where each builds on the technical debt or experimental code of the previous one to harden it into a production-grade system.

### [Milestone 1.1: Repomap Core Implementation](Milestone-1.1-Repomap-Core/)
**Status:** Primary MVP (Foundation)
Deliver the core logic for repository mapping. This milestone focuses on the "what" (parsing, graphing, ranking). It includes an MVP CLI that serves as the source for extraction in the next milestone.

### [Milestone 1.2: CLI Framework & Integration](Milestone-1.2-CLI-Framework/)
**Status:** Extraction & Hardening
**Dependency:** Milestone 1.1 Complete
Extract reusable patterns (CLI flags, output formats, config) from the Repomap MVP into a robust `pkg/` library. Refactor Repomap to use this framework, ensuring common plumbing for all future Phase 2 tools.

### [Milestone 1.3: Testing Infrastructure & CI/CD](Milestone-1.3-Testing-Infrastructure/)
**Status:** Quality Gate & Automation
**Dependency:** Milestones 1.1 and 1.2 Code Complete
Take the baseline tests from the previous milestones and scale them into an exhaustive test suite. Implement performance benchmarks and CI/CD automation to prepare for Phase 2 development.

---

## Success Criteria

- ✅ `repomap` binary compiles and runs on Go 1.19+
- ✅ Analyzes Go repositories of 1K–100K files within 30 seconds
- ✅ Respects `.gitignore` and produces correct import graphs
- ✅ Output respects token budget and prioritizes files by centrality
- ✅ All core functionality covered by exhaustive test suite (Milestone 1.3)
- ✅ CI/CD pipeline automated for unit/integration/benchmarks
- ✅ Reusable CLI framework (`pkg/`) validated by Refactored Repomap

## Phase 1 Dependency Path

```text
[M1.1: Core Logic] --> [M1.2: Library Extraction] --> [M1.3: Automation & QA]
```

## Timeline Estimate

- **Milestone 1.1:** 2–3 weeks (core development)
- **Milestone 1.2:** 1 week (Refactoring to Framework) -- *Requires M1.1*
- **Milestone 1.3:** 1 week (Testing + CI/CD setup) -- *Requires M1.1 & M1.2*
- **Total Phase 1:** ~4–5 weeks

## Handoff to Phase 2

Upon completion of Phase 1, Phase 2 will:
- Add multi-language support (Tree-sitter integration)
- Implement `semgrep-cli` using the Phase 1 CLI framework
- Extend ranking and filtering capabilities

---

## Navigation

- [Milestone 1.1: Repomap Core](Milestone-1.1-Repomap-Core/README.md)
- [Milestone 1.2: CLI Framework](Milestone-1.2-CLI-Framework/README.md)
- [Milestone 1.3: Testing Infrastructure](Milestone-1.3-Testing-Infrastructure/README.md)
