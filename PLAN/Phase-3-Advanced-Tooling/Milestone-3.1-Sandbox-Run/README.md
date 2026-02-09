# Milestone 3.1: Sandbox-Run

## Objective

Build a sandboxed code execution environment that safely executes arbitrary code with resource controls, output capture, and execution metrics.

## Scope

### In-Scope

1. **Sandbox Technology**
   - Container-based isolation (Docker/OCI)
   - Resource limits (CPU, memory, disk, time)
   - Network isolation options
   - Process tree management

2. **Execution Engine**
   - Support Go, Python, JavaScript, Rust, Java
   - Compile and execute support
   - Real-time output streaming
   - Execution metrics (time, memory, CPU)

3. **Safety Features**
   - Timeout enforcement
   - Signal handling
   - File system access control
   - Network access control

4. **sandbox-run Tool** (`cmd/sandbox-run/`)
   - CLI for running code in sandbox
   - Result reporting (stdout, stderr, exit code)
   - Performance metrics

### Out-of-Scope

- Custom kernel modules
- GPU isolation
- Advanced network routing
- Persistent state between runs

## Deliverables

1. **Sandbox Container Manager** (`pkg/sandbox/`)
2. **Execution Engine** (`pkg/executor/`)
3. **Resource Monitor** (`pkg/monitor/`)
4. **sandbox-run Tool** (`cmd/sandbox-run/`)
5. **Documentation & Examples**

## Success Criteria

- ✅ Executes code safely in isolated environment
- ✅ Enforces resource limits effectively
- ✅ Captures all output accurately
- ✅ Provides execution metrics
- ✅ Handles errors gracefully
- ✅ >80% code coverage
- ✅ All tests passing

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** ~25–30

---

## Supported Languages

| Language | Compiler | Interpreter | Priority |
|----------|----------|-------------|----------|
| Go | ✓ (native) | — | High |
| Python | — | ✓ (3.9+) | High |
| JavaScript | — | ✓ (Node.js) | High |
| Rust | ✓ | — | Medium |
| Java | ✓ | ✓ | Medium |

---

## Architecture

```
Execution Flow:
1. Code received → 2. Container started → 3. Code compiled/run → 
4. Output captured → 5. Metrics collected → 6. Container cleaned
```

## Dependencies

- Phase 1 completion
- Docker/OCI runtime
- Language toolchains
