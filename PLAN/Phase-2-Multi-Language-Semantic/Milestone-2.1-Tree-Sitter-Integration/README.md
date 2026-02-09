# Milestone 2.1: Tree-Sitter Integration

## Objective

Integrate Tree-sitter for multi-language parsing and upgrade repomap to support Go, Python, JavaScript, TypeScript, Rust, Java, and C++.

## Scope

### In-Scope

1. **Tree-Sitter Integration**
   - Add `go-tree-sitter` dependency
   - Implement language detection
   - Create parsers for 7 major languages
   - Handle language-specific constructs

2. **Repomap Upgrade**
   - Refactor parsing to use Tree-sitter
   - Update CLI with `--include-lang` and `--exclude-lang`
   - Language-specific ranking strategies
   - Extend definition extraction for each language

3. **Language Support**
   - **Go** – Functions, methods, structs, interfaces (already have)
   - **Python** – Classes, functions, decorators, async
   - **JavaScript/TypeScript** – Functions, classes, exports, interfaces
   - **Rust** – Functions, structs, traits, enums, macros
   - **Java** – Classes, methods, interfaces, annotations
   - **C++** – Classes, functions, templates, namespaces

### Out-of-Scope

- Language-specific optimizations (Phase 2.2+)
- Advanced type system support
- LSP integration (deferred to Phase 3+)

## Deliverables

1. **Tree-Sitter Wrapper Package** (`pkg/treesitter/`)
2. **Language-Specific Extractors** (one per language)
3. **Updated Repomap** (supports all 7 languages)
4. **CLI Extensions** (language selection flags)
5. **Documentation** (language support guide)
6. **Tests** (>80% coverage)

## Success Criteria

- ✅ Repomap correctly parses all 7 languages
- ✅ Definitions extracted accurately for each language
- ✅ Import graph works across language boundaries
- ✅ Ranking strategies language-aware
- ✅ No performance regression vs. Go-only
- ✅ >80% code coverage
- ✅ All tests passing

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** ~25–30

---

## Key Implementation Details

### Supported Languages

| Language | Constructs | Priority |
|----------|-----------|----------|
| Go | func, type, interface | High |
| Python | class, def, async def, decorator | High |
| JavaScript | function, class, export, const/let | High |
| TypeScript | interface, type, class, function | High |
| Rust | fn, struct, trait, enum, macro | Medium |
| Java | class, interface, method, annotation | Medium |
| C++ | class, function, namespace, template | Medium |

### Architecture

```
pkg/
├── treesitter/
│   ├── parser.go         # Tree-sitter wrapper
│   ├── language.go       # Language registry
│   └── query.go          # Query builder
├── languages/
│   ├── go/
│   ├── python/
│   ├── javascript/
│   ├── typescript/
│   ├── rust/
│   ├── java/
│   └── cpp/
└── extractor/            # Generalized extractor
    └── extractor.go
```

---

## Next Steps

1. Evaluate Tree-sitter bindings (stability, performance)
2. Create language-specific test fixtures
3. Implement extractors one language at a time
4. Update repomap incrementally
5. Benchmark performance

## Dependencies

- Phase 1 completion
- Tree-sitter and language definitions
- Language-specific test repositories
