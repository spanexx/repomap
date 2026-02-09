# Milestone 1.2: CLI Framework & Integration

## Objective

Extract and generalize reusable CLI patterns from repomap into a shared framework that all future agent tools can leverage. This milestone establishes standard conventions for flag parsing, output formatting, configuration, and error handling.

## Scope

### In-Scope

1. **CLI Framework Package** (`pkg/cli/`)
   - Flag parsing abstraction (wrapper around `flag` package)
   - Standard options structure (root, output, verbose, etc.)
   - Help text generation and customization
   - Version information management

2. **Output Formatting Utilities** (`pkg/output/`)
   - Abstract output writer interface
   - XML formatter (reusable from repomap)
   - JSON formatter (reusable from repomap)
   - Table/text formatter for CLI tools
   - Output validation utilities

3. **Configuration Management** (`pkg/config/`)
   - Config file loading (`.repomaprc`, `.semgreprc`, etc.)
   - Environment variable integration
   - Flag override hierarchy (defaults → env vars → config file → CLI flags)

4. **Error Handling** (`pkg/errors/`)
   - Standardized error types and codes
   - Error message formatting
   - User-friendly error output

5. **Shared Utilities** (`pkg/util/`)
   - Path normalization
   - Token counting utilities
   - File filtering helpers
   - Logging abstractions

6. **Documentation**
   - CLI framework guide for tool developers
   - Examples for building new tools
   - Best practices document

### Out-of-Scope

- Server/API framework (deferred to Phase 4+)
- Configuration file format standards (simple YAML/TOML for now)
- Distributed configuration management
- Authentication/authorization (deferred to Phase 3+)

## Deliverables

1. **CLI Framework Package** (`pkg/cli/`)
   - `cli.go` – Main framework
   - `flags.go` – Flag definitions
   - `help.go` – Help text generation

2. **Output Formatters** (`pkg/output/`)
   - `writer.go` – Output interface
   - `xml.go` – XML formatter
   - `json.go` – JSON formatter
   - `text.go` – Table/text formatter
   - `formatter.go` – Abstraction layer

3. **Configuration** (`pkg/config/`)
   - `config.go` – Config loader
   - `env.go` – Environment variable support
   - `hierarchy.go` – Config override logic

4. **Error Handling** (`pkg/errors/`)
   - `errors.go` – Error types
   - `codes.go` – Exit codes
   - `messages.go` – Error message templates

5. **Utilities** (`pkg/util/`)
   - `paths.go` – Path utilities
   - `tokens.go` – Token counting
   - `filter.go` – File filtering
   - `logging.go` – Logging helpers

6. **Integration Guide**
   - `FRAMEWORK_GUIDE.md` – How to use the CLI framework
   - `EXAMPLES.md` – Example implementations
   - `BEST_PRACTICES.md` – Standards and conventions

## Success Criteria

- ✅ CLI framework is reusable by at least 2 future tools (semgrep-cli, web-ray)
- ✅ All output formatters work with any tool's data structure
- ✅ Configuration hierarchy works correctly (env vars override defaults, etc.)
- ✅ Error handling is consistent across all tools
- ✅ Framework documentation is clear and includes examples
- ✅ Refactored repomap uses the new framework without losing functionality
- ✅ No breaking changes to repomap CLI interface
- ✅ >80% code coverage on framework packages
- ✅ Framework adds <10% overhead to CLI startup time

## Key Implementation Details

### Architecture

```
pkg/
├── cli/
│   ├── cli.go          # Main framework
│   ├── flags.go        # Flag definitions
│   └── help.go         # Help text generation
├── output/
│   ├── writer.go       # Output interface
│   ├── xml.go          # XML formatter
│   ├── json.go         # JSON formatter
│   ├── text.go         # Text formatter
│   └── formatter.go    # Formatter factory
├── config/
│   ├── config.go       # Config loader
│   ├── env.go          # Environment support
│   └── hierarchy.go    # Override logic
├── errors/
│   ├── errors.go       # Error types
│   ├── codes.go        # Exit codes
│   └── messages.go     # Error messages
└── util/
    ├── paths.go        # Path utilities
    ├── tokens.go       # Token counting
    ├── filter.go       # File filtering
    └── logging.go      # Logging helpers
```

### CLI Framework Example

```go
package main

import "github.com/spanexx/agents-cli/pkg/cli"

func main() {
    app := cli.NewApp("mytool", "1.0.0")
    app.AddFlag("--root", "Repository root", ".")
    app.AddFlag("--output", "Output format", "json")
    
    flags := app.Parse(os.Args[1:])
    
    // Use framework's output writer
    writer := app.NewOutputWriter(flags.Output)
    writer.Write([]byte(`{"result": "success"}`))
}
```

### Configuration Hierarchy

```
1. Default values (hardcoded)
   ↓
2. Config file (~/.semgreprc, ./.semgreprc)
   ↓
3. Environment variables (SEMGREP_ROOT, SEMGREP_OUTPUT)
   ↓
4. CLI flags (--root, --output)
   ↓
Final configuration
```

## Testing Strategy

### Unit Tests
- Flag parsing and validation
- Configuration loading and hierarchy
- Error type construction
- Output formatter validity
- Utility function correctness

### Integration Tests
- CLI framework with repomap
- Multiple tools sharing framework
- Configuration file loading
- Output formatting with various data types

### Compatibility Tests
- Repomap CLI interface unchanged after refactoring
- All repomap tests still pass
- No performance regression

## Definition of Done

- [ ] All CLI framework packages implemented
- [ ] Output formatters work with multiple data types
- [ ] Configuration management tested
- [ ] Error handling standardized
- [ ] Repomap refactored to use framework
- [ ] All tests passing (>80% coverage)
- [ ] Documentation complete with examples
- [ ] No breaking changes to repomap CLI
- [ ] Code reviewed and approved
- [ ] Framework tested with 2+ future tools (stubs)

## Dependencies

- Milestone 1.1 (Repomap Core) – framework patterns based on repomap implementation
- Go 1.19+ (same as Phase 1)

## Timeline Estimate

- **Duration:** 1 week
- **Parallel with:** Milestone 1.3 (testing can start once core framework is done)

## Risks

| Risk | Mitigation |
|------|-----------|
| Over-engineering the framework | Keep it simple; only generalize what's actually needed |
| Breaking repomap during refactor | Use feature flags and comprehensive testing |
| Future tools don't fit the framework | Design for extensibility; document patterns clearly |

## Success Metrics

| Metric | Target |
|--------|--------|
| Framework reusability | Used by 2+ tools in Phase 2 |
| Code coverage | >80% |
| Performance overhead | <10% startup time increase |
| API stability | No breaking changes in Phase 2+ |

## Next Phase

Milestone 1.3 will test both Repomap Core and CLI Framework comprehensively, ensuring the foundation is solid for Phase 2.
