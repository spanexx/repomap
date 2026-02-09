# Milestone 2.1: Progress Tracking

## Phase Overview
Multi-language parsing support through Tree-sitter integration. Extends Phase 1 repomap to support 7 programming languages with language-aware ranking strategies.

---

## Task Status Matrix

| Task | Title | Status | Owner | % Complete | Target Date |
|------|-------|--------|-------|------------|------------|
| T1 | Evaluate Tree-Sitter Bindings | — | — | 0% | — |
| T2 | Set Up Tree-Sitter Dependency | — | — | 0% | — |
| T3 | Create Tree-Sitter Wrapper | — | — | 0% | — |
| T4 | Implement Language Registry | — | — | 0% | — |
| T5 | Create Go Extractor | — | — | 0% | — |
| T6 | Create Python Extractor | — | — | 0% | — |
| T7 | Create JavaScript Extractor | — | — | 0% | — |
| T8 | Create TypeScript Extractor | — | — | 0% | — |
| T9 | Create Rust Extractor | — | — | 0% | — |
| T10 | Create Java Extractor | — | — | 0% | — |
| T11 | Create C++ Extractor | — | — | 0% | — |
| T12 | Create Generalized Interface | — | — | 0% | — |
| T13 | Update Repomap Integration | — | — | 0% | — |
| T14 | Test Go Support | — | — | 0% | — |
| T15 | Test Python Support | — | — | 0% | — |
| T16 | Test JavaScript Support | — | — | 0% | — |
| T17 | Test TypeScript Support | — | — | 0% | — |
| T18 | Test Rust Support | — | — | 0% | — |
| T19 | Test Java Support | — | — | 0% | — |
| T20 | Test C++ Support | — | — | 0% | — |
| T21 | Add CLI Language Flags | — | — | 0% | — |
| T22 | Language-Aware Ranking | — | — | 0% | — |
| T23 | Performance Benchmarking | — | — | 0% | — |
| T24 | Documentation | — | — | 0% | — |
| T25 | Final Integration Tests | — | — | 0% | — |

---

## Burn-Down Chart

```
Week 1: Setup & Planning
  - Complete: T1, T2, T3, T4
  - In Progress: T5
  - Remaining: T6-T25

Week 2: Language Extractors
  - Complete: T5, T6, T7, T8, T9, T10, T11, T12, T13
  - In Progress: T14, T15
  - Remaining: T16-T25

Week 3: Testing & Polish
  - Complete: T14-T23
  - In Progress: T24, T25
  - Remaining: None
```

---

## Risk Register

| Risk | Severity | Probability | Mitigation |
|------|----------|-------------|-----------|
| Tree-sitter binding unstable | High | Medium | Start with evaluation, test early |
| Language definition accuracy | High | Low | Use official language grammars |
| Performance degradation | Medium | Medium | Benchmark continuously |
| Import resolution complexity | Medium | High | Start simple, iterate |

---

## Success Criteria Tracking

- [ ] All 7 languages supported with extractors
- [ ] 100% definition extraction accuracy on test repos
- [ ] >80% code coverage
- [ ] All tests passing
- [ ] No performance regression
- [ ] Language flags working in CLI
- [ ] Documentation complete

---

## Notes & Decisions

- **Language Priorities:** Go (Phase 1), Python/JS/TS (high demand), Rust/Java/C++ (common)
- **Fallback Strategy:** If Tree-sitter unstable, may implement simplified parsers
- **Performance:** Plan benchmarking early to catch regressions
