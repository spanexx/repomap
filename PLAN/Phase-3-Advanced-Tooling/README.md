# Phase 3: Advanced Tooling

## Objective

Build the "Hands" and "Brain" tools that extend agent capabilities to execute code safely, interact with APIs intelligently, manage version control, and maintain long-term memory.

## Scope

### In-Scope for Phase 3

1. **Sandboxed Execution** (Milestone 3.1)
   - `sandbox-run` – execute shell commands in isolated environments
   - Docker container or Firecracker VM support
   - Structured JSON output with stdout/stderr/exit_code
   - Dry-run mode with impact prediction
   - Resource limits and timeouts

2. **API Integration** (Milestone 3.2)
   - `api-forge` – intelligent API interaction
   - Schema reverse-engineering from errors
   - Fuzzy endpoint exploration
   - Request generation and validation
   - Error recovery and retry strategies

3. **Git Operations** (Milestone 3.3)
   - `git-surgeon` – merge conflict resolution
   - Parse and understand conflict markers
   - Suggest resolutions using AST analysis
   - JSON API for conflict manipulation
   - Support for various conflict patterns

4. **Long-Term Memory** (Milestone 3.4)
   - `mem-kv` – persistent key-value store
   - Store facts, preferences, summaries
   - Cross-session access and retrieval
   - Encryption at rest
   - TTL and expiration support

5. **Math Co-Processor** (Milestone 3.5)
   - `calc-bridge` – symbolic computation
   - SymPy integration for complex math
   - Natural language math query parsing
   - Data transformation and aggregation
   - Deterministic results

### Out-of-Scope for Phase 3

- Distributed memory (deferred to Phase 4)
- Advanced cloud integrations (deferred to Phase 4)
- Performance optimization (to Phase 4)

## Milestones

### [Milestone 3.1: Sandbox-Run](Milestone-3.1-Sandbox-Run/)
Safe code execution in isolated environments.

### [Milestone 3.2: API-Forge](Milestone-3.2-API-Forge/)
Intelligent API interaction and schema detection.

### [Milestone 3.3: Git-Surgeon](Milestone-3.3-Git-Surgeon/)
Merge conflict resolution and git manipulation.

### [Milestone 3.4: Mem-KV](Milestone-3.4-Mem-KV/)
Persistent memory store for cross-session state.

### [Milestone 3.5: Calc-Bridge](Milestone-3.5-Calc-Bridge/)
Math and symbolic computation co-processor.

## Success Criteria

- ✅ sandbox-run safely executes commands with resource limits
- ✅ api-forge can reverse-engineer schemas from API responses
- ✅ git-surgeon resolves common merge conflicts correctly
- ✅ mem-kv persists and retrieves data across sessions
- ✅ calc-bridge handles complex mathematical operations
- ✅ All tools use CLI framework
- ✅ >80% code coverage on all tools
- ✅ Comprehensive documentation and examples

## Dependencies

- **Internal:** Phase 1 & 2 must be complete
- **External:** Docker/Firecracker, SymPy, SQLite or equivalent

## Timeline Estimate

- **Milestone 3.1:** 3–4 weeks (Sandbox implementation)
- **Milestone 3.2:** 2–3 weeks (API schema detection)
- **Milestone 3.3:** 2–3 weeks (Git operations)
- **Milestone 3.4:** 2 weeks (Memory store)
- **Milestone 3.5:** 2 weeks (Math processor)
- **Total Phase 3:** ~11–15 weeks

## Handoff to Phase 4

Upon completion of Phase 3, Phase 4 will:
- Integrate all tools into unified server
- Add authentication and authorization
- Deploy to cloud platforms
- Add monitoring and observability

---

## Navigation

- [Milestone 3.1: Sandbox-Run](Milestone-3.1-Sandbox-Run/README.md)
- [Milestone 3.2: API-Forge](Milestone-3.2-API-Forge/README.md)
- [Milestone 3.3: Git-Surgeon](Milestone-3.3-Git-Surgeon/README.md)
- [Milestone 3.4: Mem-KV](Milestone-3.4-Mem-KV/README.md)
- [Milestone 3.5: Calc-Bridge](Milestone-3.5-Calc-Bridge/README.md)
