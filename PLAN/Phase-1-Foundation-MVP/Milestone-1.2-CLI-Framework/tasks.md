# Tasks: Milestone 1.2 – CLI Framework & Integration

## Dependencies
- **Prerequisites:** 
  - Milestone 1.1 Complete (specifically CLI logic 1.1.13 and Output logic 1.1.10)
- **Provides for:**
  - Milestone 1.3: Refactored library code for multi-tool integration testing.


## Framework Foundation

### Task 1.2.1: CLI Framework Base Package
- [ ] Create `pkg/cli/` package structure
- [ ] Define `App` struct with fields:
  - `name`, `version`
  - Flags registry
  - Output writer
  - Error handler
- [ ] Implement `NewApp(name, version) *App`
- [ ] Write unit tests

**Acceptance Criteria:**
- App struct is well-defined
- Can be instantiated and configured
- Unit tests pass

---

### Task 1.2.2: Flag Parsing Abstraction
- [ ] Create `flags.go` with `Flags` struct
- [ ] Implement `AddFlag(name, description, default) *App` method
- [ ] Implement `Parse(args []string) (*Flags, error)` method
- [ ] Support common flag types: string, int, bool, []string
- [ ] Generate help text from flags
- [ ] Write unit tests on various flag combinations

**Acceptance Criteria:**
- Flags parse correctly
- Help text is generated automatically
- Error handling for invalid flags
- Unit tests pass with >90% coverage

---

### Task 1.2.3: Help Text Generation
- [ ] Implement `help.go` with `GenerateHelp(app *App) string`
- [ ] Format help text with usage, options, examples
- [ ] Make help text customizable
- [ ] Write unit tests

**Acceptance Criteria:**
- Help text is formatted correctly
- Includes all flags and descriptions
- Easy to read and understand

---

## Output Formatting

### Task 1.2.4: Output Writer Interface
- [ ] Define `Writer` interface:
  ```go
  type Writer interface {
    Write(data interface{}) error
    WriteString(s string) error
  }
  ```
- [ ] Create `output.go` with factory `NewWriter(format string) Writer`
- [ ] Write unit tests

**Acceptance Criteria:**
- Interface is clean and extensible
- Factory pattern works
- Unit tests pass

---

### Task 1.2.5: Reusable XML Formatter
- [ ] Extract XML formatting logic from repomap
- [ ] Create `xml.go` with `XMLWriter` implementing `Writer`
- [ ] Generalize to handle any struct with tags
- [ ] Write unit tests on various data types

**Acceptance Criteria:**
- Works with repomap data without changes
- Can format other struct types
- Produces valid XML
- Unit tests pass

---

### Task 1.2.6: Reusable JSON Formatter
- [ ] Extract JSON formatting logic from repomap
- [ ] Create `json.go` with `JSONWriter` implementing `Writer`
- [ ] Generalize to handle any struct
- [ ] Write unit tests on various data types

**Acceptance Criteria:**
- Works with repomap data without changes
- Can format other struct types
- Produces valid JSON
- Unit tests pass

---

### Task 1.2.7: Text/Table Formatter
- [ ] Create `text.go` with `TextWriter` for simple text/table output
- [ ] Support table formatting with columns, alignment
- [ ] Support simple key-value output
- [ ] Write unit tests

**Acceptance Criteria:**
- Tables are well-formatted
- Easy to read
- Unit tests pass

---

### Task 1.2.8: Output Writer Factory
- [ ] Implement `formatter.go` with smart factory
- [ ] Route to correct formatter based on format string
- [ ] Handle unknown formats gracefully
- [ ] Write unit tests

**Acceptance Criteria:**
- Factory correctly instantiates formatters
- Error handling for unknown formats
- Unit tests pass

---

## Configuration Management

### Task 1.2.9: Configuration Loader
- [ ] Create `config/config.go` with `Config` struct
- [ ] Implement `LoadConfig(root string) (*Config, error)`
- [ ] Read `.repomaprc`, `.semgreprc`, etc. (tool-specific)
- [ ] Parse YAML/TOML/JSON formats
- [ ] Write unit tests with test fixtures

**Acceptance Criteria:**
- Reads config files correctly
- Handles various formats
- Graceful handling of missing files
- Unit tests pass

---

### Task 1.2.10: Environment Variable Support
- [ ] Implement `env.go` with `GetEnv(key, default) string`
- [ ] Support tool-specific prefixes (REPOMAP_, SEMGREP_, etc.)
- [ ] Document environment variable naming conventions
- [ ] Write unit tests

**Acceptance Criteria:**
- Environment variables are read correctly
- Prefixes work as expected
- Unit tests pass

---

### Task 1.2.11: Configuration Hierarchy
- [x] Implement `hierarchy.go` with merge logic
- [x] Apply hierarchy: defaults → env → config file → CLI flags
- [x] Last-in-wins override strategy
- [x] Write unit tests on various scenarios

**Acceptance Criteria:**
- Hierarchy is applied correctly
- CLI flags override config files
- Environment variables work as expected
- Unit tests pass

---

## Error Handling

### Task 1.2.12: Error Types & Codes
- [x] Create `errors/errors.go` with custom error types:
  - `ConfigError`
  - `ParseError`
  - `ExecutionError`
  - `ValidationError`
- [x] Create `codes.go` with standardized exit codes
- [x] Write unit tests

**Acceptance Criteria:**
- Error types are clear and useful
- Exit codes are consistent
- Unit tests pass

---

### Task 1.2.13: Error Message Templates
- [x] Create `messages.go` with standardized error messages
- [x] Support message templating with context
- [x] Write user-friendly error messages
- [x] Write unit tests

**Acceptance Criteria:**
- Messages are clear and actionable
- Support variable substitution
- Unit tests pass

---

## Shared Utilities

### Task 1.2.14: Path Utilities
- [x] Create `util/paths.go` with functions:
  - `NormalizePath(path string) string`
  - `IsAbsolutePath(path string) bool`
  - `MakeRelative(basePath, targetPath string) string`
- [x] Write unit tests on various paths

**Acceptance Criteria:**
- Path operations work correctly
- Handle Windows and Unix paths
- Unit tests pass

---

### Task 1.2.15: Token Counting Utilities
- [ ] Create `util/tokens.go` with reusable token counting
- [ ] Implement `CountTokens(text string) int`
- [ ] Document approximation method
- [ ] Write unit tests

**Acceptance Criteria:**
- Token counting is reusable
- Consistent with repomap
- Unit tests pass

---

### Task 1.2.16: File Filtering Helpers
- [ ] Create `util/filter.go` with:
  - `FilterByExtension(paths []string, exts []string) []string`
  - `ExcludeByPattern(paths []string, patterns []string) []string`
- [ ] Write unit tests

**Acceptance Criteria:**
- Filtering works correctly
- Handles edge cases
- Unit tests pass

---

### Task 1.2.17: Logging Abstraction
- [ ] Create `util/logging.go` with logger interface
- [ ] Implement simple stdout/stderr logging
- [ ] Support verbose mode
- [ ] Write unit tests

**Acceptance Criteria:**
- Logging is configurable
- Easy to use
- Unit tests pass

---

## Integration & Refactoring

### Task 1.2.18: Refactor Repomap to Use Framework
- [ ] Update `cmd/repomap/main.go` to use `pkg/cli`
- [ ] Update output formatting to use `pkg/output`
- [ ] Update configuration to use `pkg/config`
- [ ] Update error handling to use `pkg/errors`
- [ ] Verify all repomap tests still pass
- [ ] Run integration tests

**Acceptance Criteria:**
- Repomap works identically to before
- All repomap tests pass
- No CLI interface changes
- Code is cleaner and more maintainable

---

### Task 1.2.19: Framework Documentation
- [ ] Write `FRAMEWORK_GUIDE.md`:
  - Overview of framework components
  - How to add a new tool
  - API reference
- [ ] Write `EXAMPLES.md`:
  - Example tool using the framework
  - Step-by-step walkthrough
- [ ] Write `BEST_PRACTICES.md`:
  - CLI conventions
  - Error handling patterns
  - Output formatting guidelines
- [ ] Write godoc comments on all exports

**Acceptance Criteria:**
- Documentation is comprehensive
- Includes working examples
- Clear instructions for new tools

---

## Testing & Validation

### Task 1.2.20: Unit Tests (All Packages)
- [ ] Achieve >80% code coverage on `pkg/`
- [ ] Run: `go test ./pkg/... -cover`
- [ ] Write tests for all error paths

**Acceptance Criteria:**
- Coverage >80%
- All tests pass
- Error handling tested

---

### Task 1.2.21: Integration Tests
- [ ] Test framework with repomap
- [ ] Test framework with mock tools
- [ ] Verify configuration hierarchy works
- [ ] Test error handling end-to-end

**Acceptance Criteria:**
- Integration tests pass
- Framework works with multiple tools
- No regressions in repomap

---

### Task 1.2.22: Framework Compatibility Tests
- [ ] Run all repomap tests after refactoring
- [ ] Verify repomap CLI interface unchanged
- [ ] Performance benchmarks (should be minimal overhead)

**Acceptance Criteria:**
- Repomap tests pass 100%
- CLI interface unchanged
- Performance overhead <10%

---

## Summary

**Total Tasks:** 22
**Estimated Effort:** 1 week
**Dependencies:** Milestone 1.1 (Repomap Core)
**Success:** All tasks completed + tests passing + documentation complete + repomap refactored
