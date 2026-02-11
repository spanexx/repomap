# Tasks: Milestone 1.4 - Planning & Intent Analysis

This milestone introduces features to transform `repomap` into a proactive planning tool for AI agents.

## Schema & Core Planning Model

- [x] **Task 1.4.1**: Update `FileNode` Schema @[repomap/visualizer-ui/src/types.ts]
  - Add optional fields: `Status`, `Intent`, `Issues` (array of objects), `Comments` (array of objects)
  - Ensure Go struct (`internal/output/types.go`) matches TypeScript interface.
  - Test serialization/deserialization.

- [x] **Task 1.4.2**: Implement `internal/planning` Package
  - Create package for planning logic.
  - Define `Planner` interface and implementation.
  - Implement function to load/parse `.repomapintent` or `plan.json`.

## Static Analysis (DRY & Duplication)

- [x] **Task 1.4.3**: Integrate `repomap` with `gemini` (or `qodercli`)
  - [x] Create/Update `pkg/providers/gemini_cli` wrapper
  - [x] Remove internal `gemini` API client
    - [x] Implement session management (`internal/session`)
    - [x] Update `pkg/server/server.go` to use `gemini` CLI wrapper
    - [x] Add `GET /api/chat` for history persistence
    - [x] Implement **New Session** button in UI
    - [ ] Implement auto-summarization for long context (Deferred)

- [x] **Task 1.4.4**: Implement DRY Violation Detection
  - [x] Create utility to tokenize Go code (ignoring comments/whitespace).
  - [x] Use rolling hash or similar to detect duplicate blocks.
  - [x] Identify functions with >80% similarity.
  - [x] Scan the AST for identical structures (functions, methods).
  - [x] Flag violations in the `Issues` field of `FileNode`.
  - [x] Report file paths and line numbers.

- [x] **Task 1.4.5**: Implement Intent Validation
  - [x] Check imports against intent constraints (e.g., "modules in `domain` cannot import `infrastructure`").
  - [x] Flag circular dependencies.
  - [x] Flag violations in `Issues`.

## CLI Integration

- [x] **Task 1.4.6**: Add `--plan` Flag
  - [x] Allow loading a `plan.json` file.
  - [x] Merge the plan with the scanned repository graph.
  - [x] Ensure planned nodes are marked with `status: "planned"`.

- [x] **Task 1.4.7**: Add `--analyze` Flag
  - [x] Trigger the static analysis (Tasks 1.4.3 - 1.4.5).
  - [x] Populate `issues` field in the output.
  - [x] Output summary report to console if verbose.

## Visualizer Updates (Backend Support)

- [x] **Task 1.4.8**: Verify JSON Output Compatibility
  - [x] Ensure new fields are correctly marshaled to JSON.
  - [x] Verify Visualizer UI can load and display the enhanced JSON.
  - [x] Test with sample `plan.json` + `repomap.json`.

## Documentation & Testing

- [ ] **Task 1.4.9**: Document Planning Features
  - Update `README.md` with usage instructions for `--plan` and `--analyze`.
  - Provide examples of `plan.json` format.
  - Explain how AI agents can leverage this for architecture validation.

- [ ] **Task 1.4.10**: Unit & Integration Tests
  - Test duplication detection with fixtures.
  - Test plan merging logic.
  - Regression test existing functionality.

## Milestone 1.5: Interactive Flow & Chat Integration

- [ ] **Task 1.5.1**: Chat-Aware Node Highlighting
  - [ ] Update `App.tsx` to manage highlighted files state.
  - [ ] Update `Chat.tsx` to detect filenames in messages and trigger highlight.
  - [ ] Update `Graph.tsx` to visually isolate highlighted nodes (dim others).

- [ ] **Task 1.5.2**: Enhanced Graph Interactivity
  - [ ] Add "Focus" mode to only show connected subgraph of highlighted nodes.
  - [ ] Add animations for focus transitions.
