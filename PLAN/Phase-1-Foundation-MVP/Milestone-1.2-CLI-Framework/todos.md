# TODO: CLI Framework Refactoring & Improvements

## Immediate Actions
- [ ] **Refactor Output Writers to use `io.Writer`**
  - **Why:** Current implementation writes directly to `os.Stdout`, making unit testing unreliable (hijacking stdout) and preventing output redirection.
  - **Plan:**
    - Update `NewWriter` factory to accept an `io.Writer`.
    - Update `XMLWriter`, `JSONWriter`, and `TextWriter` to write to the provided writer instead of `fmt.Println`.
    - Update unit tests to use a `bytes.Buffer` and verify the output string.

## Future Improvements
- [ ] **Error Handling in `AddFlag`**
  - **Why:** Currently `AddFlag` panics on unsupported types.
  - **Plan:** Return an error instead of panicking, allowing the caller to handle configuration errors gracefully.

- [x] **Wire up qwen_cli & qoder_cli**
  - Updated `pkg/providers/qwen_cli` and `qodercli` to use `repomap/pkg/adapter`.
  - Updated `pkg/server/server.go` to select provider via `REPOMAP_PROVIDER`.

- [x] **Implement Provider Fallback Strategy**
  - Copied `FallbackProvider` logic from `llm-adapter`.
  - Updated `pkg/server/server.go` to support comma-separated list of providers (e.g., `qwen-cli,gemini-cli`).
  - Verified fallback behavior by forcing `qwen-cli` to fail and observing `qodercli` take over.
