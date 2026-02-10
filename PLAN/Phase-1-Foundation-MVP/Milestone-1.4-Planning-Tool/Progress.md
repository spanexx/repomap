# Tasks: Milestone 1.4 - Planning & Intent Analysis

This milestone introduces features to transform `repomap` into a proactive planning tool for AI agents.

## Schema & Core Planning Model

- [ ] **Task 1.4.1**: Update `FileNode` Schema @[repomap/visualizer-ui/src/types.ts]
  - Add optional fields: `Status`, `Intent`, `Issues` (array of objects), `Comments` (array of objects)
  - Ensure Go struct (`internal/output/types.go`) matches TypeScript interface.
  - Test serialization/deserialization.

- [ ] **Task 1.4.2**: Implement `internal/planning` Package
  - Create package for planning logic.
  - Define `Planner` interface and implementation.
  - Implement function to load/parse `.repomapintent` or `plan.json`.

## Static Analysis (DRY & Duplication)

- [ ] **Task 1.4.3**: Implement Token-Based Hashing
  - Create utility to tokenize Go code (ignoring comments/whitespace).
  - Use rolling hash or similar to detect duplicate blocks.
  - Identify functions with >80% similarity.

- [ ] **Task 1.4.4**: Implement DRY Violation Detection
  - Scan the AST for identical structures (functions, methods).
  - Flag violations in the `Issues` field of `FileNode`.
  - Report file paths and line numbers.

- [ ] **Task 1.4.5**: Implement Intent Validation
  - Check imports against intent constraints (e.g., "modules in `domain` cannot import `infrastructure`").
  - Flag circular dependencies.
  - Flag violations in `Issues`.

## CLI Integration

- [ ] **Task 1.4.6**: Add `--plan` Flag
  - Allow loading a `plan.json` file.
  - Merge the plan with the scanned repository graph.
  - Ensure planned nodes are marked with `status: "planned"`.

- [ ] **Task 1.4.7**: Add `--analyze` Flag
  - Trigger the static analysis (Tasks 1.4.3 - 1.4.5).
  - Populate `issues` field in the output.
  - Output summary report to console if verbose.

## Visualizer Updates (Backend Support)

- [ ] **Task 1.4.8**: Verify JSON Output Compatibility
  - Ensure new fields are correctly marshaled to JSON.
  - Verify Visualizer UI can load and display the enhanced JSON.
  - Test with sample `plan.json` + `repomap.json`.

## Documentation & Testing

- [ ] **Task 1.4.9**: Document Planning Features
  - Update `README.md` with usage instructions for `--plan` and `--analyze`.
  - Provide examples of `plan.json` format.
  - Explain how AI agents can leverage this for architecture validation.

- [ ] **Task 1.4.10**: Unit & Integration Tests
  - Test duplication detection with fixtures.
  - Test plan merging logic.
  - Regression test existing functionality.
