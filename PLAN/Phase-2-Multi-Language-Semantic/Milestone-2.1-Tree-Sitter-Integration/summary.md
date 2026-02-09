# Milestone 2.1: Quick Reference

## Summary

Add multi-language parsing with Tree-sitter, extending repomap from Go-only to 7 languages (Go, Python, JavaScript, TypeScript, Rust, Java, C++).

## Timeline

- **Duration:** 3–4 weeks
- **Tasks:** 25 sequential
- **Target Completion:** Week 7 of Phase 2

## Key Components

| Component | Path | Purpose |
|-----------|------|---------|
| Parser Wrapper | `pkg/treesitter/parser.go` | Tree-sitter abstraction layer |
| Language Registry | `pkg/treesitter/language.go` | Language detection & routing |
| Extractors | `pkg/languages/{lang}/extractor.go` | Language-specific definition extraction |
| Interface | `pkg/extractor/extractor.go` | Unified extractor interface |
| CLI Update | `cmd/repomap/main.go` | Add language flags |

## Execution Strategy

1. **Setup Phase (Week 1)**
   - Evaluate Tree-sitter bindings
   - Create wrapper and registry
   - Set up infrastructure

2. **Extraction Phase (Week 2)**
   - Implement extractors for all 7 languages
   - Create generalized interface
   - Update repomap parser

3. **Validation Phase (Week 3)**
   - Test each language support
   - Benchmark performance
   - Document findings

## Success Criteria

| Criterion | Target | Status |
|-----------|--------|--------|
| Languages Supported | 7 | — |
| Extraction Accuracy | 100% | — |
| Code Coverage | >80% | — |
| Test Pass Rate | 100% | — |
| Performance | No regression | — |

## Acceptance Criteria

- ✅ Repomap parses all 7 languages correctly
- ✅ Definitions extracted with 100% accuracy
- ✅ >80% code coverage achieved
- ✅ CLI flags work: `repomap --include-lang go,python`
- ✅ Language-aware ranking implemented
- ✅ Complete documentation with examples

## Metrics

| Metric | Baseline | Target |
|--------|----------|--------|
| Parse Time (10k LOC) | — | <5s |
| Memory per 1k LOC | — | <10MB |
| Definitions Found | 100% | 100% |
