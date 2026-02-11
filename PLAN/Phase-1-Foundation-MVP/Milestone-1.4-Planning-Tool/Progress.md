# Tasks: Milestone 1.4 - Planning & Intent Analysis

This milestone introduces features to transform `repomap` into a proactive planning tool for AI agents.

## Schema & Core Planning Model

### Task 1.4.1: Update Data Schema
**Status:** ✅ Complete
- **2026-02-10**: Updated `FileNode` struct in `internal/output/types.go` with `Status`, `Intent`, `Issues`, `Comments`. Added `planning_test.go` to verify JSON serialization.

### Task 1.4.2: Planning Logic Package
**Status:** ✅ Complete
- **2026-02-10**: Implemented `internal/planning` package with `Planner` struct, `LoadPlan` schema parsing, and `ApplyPlan` logic to merge changes into the repomap. Verified with unit tests.

### Task 1.4.3: External Gemini CLI Wrapper
**Status:** ✅ Complete
- **2026-02-10**: Replaced internal provider with CLI wrapper to resolve auth issues. Implemented chat history persistence and "New Session" feature.

### Task 1.4.4: DRY Detection
**Status:** ✅ Complete
- **2026-02-10**: Implemented `DuplicationDetector` using rolling hash to identify duplicate code blocks. Verified with unit tests and self-analysis.

### Task 1.4.5: Intent Validation
**Status:** ✅ Complete
- **2026-02-10**: Implemented `IntentValidator` to enforce Clean Architecture boundaries and `DependencyGraph` to detect circular dependencies.

### Task 1.4.6 & 1.4.7: CLI Integration (--plan & --analyze)
**Status:** ✅ Complete
- **2026-02-10**: Added `--plan` to load architectural plans and `--analyze` to trigger static analysis. Integrated into `repomap` CLI.

### Task 1.4.8: Visualizer Compatibility
**Status:** ✅ Complete
- **2026-02-10**: Verified JSON output contains `issues` field. Visualizer UI updated to render Markdown in chat.

### Task 1.4.9: Documentation
**Status:** ✅ Complete
- **2026-02-10**: Updated `USAGE.md` and created `walkthrough.md` covering new features.

## Documentation & Testing

| Status | Count |
|---|---|
| ✅ Complete | 9 |
| 🟨 In Progress | 0 |
| 🔴 Blocked | 0 |
| ⬜ Not Started | 0 |
| **Total** | **9** |
