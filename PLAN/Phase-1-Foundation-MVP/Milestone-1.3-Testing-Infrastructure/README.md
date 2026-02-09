# Milestone 1.3: Testing Infrastructure & CI/CD

## Objective

Establish comprehensive testing infrastructure, performance benchmarking, and CI/CD pipelines to ensure code quality, reliability, and continuous integration for both Repomap Core and the CLI Framework.

## Scope

### In-Scope

1. **Unit Testing**
   - Achieve >85% code coverage across all packages
   - Test frameworks and utilities
   - Test error paths and edge cases

2. **Integration Testing**
   - Full pipeline tests on real Go repositories
   - Multi-tool integration tests
   - End-to-end CLI testing

3. **Performance Benchmarking**
   - Benchmark discovery, parsing, ranking phases
   - Memory profiling
   - Startup time measurement
   - Track metrics over time

4. **CI/CD Pipeline**
   - GitHub Actions or equivalent
   - Automated testing on commits
   - Code coverage reporting
   - Binary releases
   - Automated deployment

5. **Testing Documentation**
   - Test strategy and procedures
   - How to add new tests
   - Performance benchmarking guide
   - CI/CD configuration documentation

### Out-of-Scope

- Distributed testing infrastructure (deferred to Phase 4+)
- End-to-end testing with external services (deferred to Phase 3+)
- Performance optimization (optimize based on benchmarks in Phase 2)
- Load testing (deferred to Phase 4+)

## Deliverables

1. **Test Suites**
   - Unit tests for repomap (discovery, parsing, graph, ranking, output)
   - Unit tests for CLI framework (cli, output, config, errors, util)
   - Integration tests for full pipeline
   - Test fixtures and test data

2. **Performance Benchmarks**
   - Go benchmark files (`*_bench_test.go`)
   - Benchmark results and baselines
   - Performance tracking script

3. **CI/CD Pipeline**
   - `.github/workflows/` YAML files
   - Build job (go build on multiple platforms)
   - Test job (unit + integration tests)
   - Coverage reporting job
   - Release job (create binaries)

4. **Documentation**
   - `TESTING.md` – Testing strategy and procedures
   - `BENCHMARKING.md` – How to run and interpret benchmarks
   - `CI_CD.md` – CI/CD pipeline documentation
   - `CONTRIBUTING.md` – Contributing guidelines

## Success Criteria

- ✅ >85% code coverage on repomap and framework packages
- ✅ All unit tests pass locally and in CI
- ✅ Integration tests pass on 5+ real Go repositories
- ✅ Performance benchmarks recorded and tracked
- ✅ CI/CD pipeline runs automatically on push/PR
- ✅ Binary releases created for Linux, macOS, Windows
- ✅ Code coverage reports published
- ✅ Test failures block merges to main
- ✅ Documentation is comprehensive and clear
- ✅ New contributors can run tests locally in <5 minutes

## Key Implementation Details

### Testing Structure

```
├── repomap/
│   ├── internal/
│   │   ├── discovery/
│   │   │   ├── walker_test.go
│   │   │   ├── walker_bench_test.go
│   │   │   └── testdata/
│   │   ├── parsing/
│   │   │   ├── extractor_test.go
│   │   │   ├── extractor_bench_test.go
│   │   │   └── testdata/
│   │   └── ...
│   └── cmd/
│       └── repomap/
│           └── main_test.go
├── pkg/
│   ├── cli/
│   │   ├── cli_test.go
│   │   ├── flags_test.go
│   │   └── testdata/
│   ├── output/
│   │   ├── xml_test.go
│   │   ├── json_test.go
│   │   └── testdata/
│   └── ...
└── test/
    ├── fixtures/          # Real Go projects for integration tests
    ├── integration_test.go
    └── benchmarks_test.go
```

### Test Fixtures

```
test/fixtures/
├── simple-go-app/         # Small Go project (10-20 files)
│   ├── main.go
│   ├── go.mod
│   └── util/
├── medium-go-project/     # Medium project (100-200 files)
└── large-go-project/      # Large project (1K+ files)
```

### CI/CD Pipeline Structure

```
.github/workflows/
├── test.yml               # Run tests on PR/push
├── coverage.yml           # Generate coverage reports
├── release.yml            # Create releases on tag
└── benchmark.yml          # Track performance over time
```

## Testing Strategy

### Unit Testing (TDD)
- Write tests first, then implementation
- Test both happy path and error cases
- Mock external dependencies (files, networks)
- Aim for >90% coverage on critical paths

### Integration Testing
- Use real test fixtures (Go repositories)
- Test full pipeline from discovery to output
- Verify correctness of output (XML/JSON validity)
- Test with various `.gitignore` patterns

### Performance Testing
- Benchmark on fixed-size test sets (1K, 10K, 100K files)
- Track metrics: memory, time, throughput
- Establish baselines in Phase 1
- Improve in Phase 2

### Regression Testing
- Prevent regressions from new features
- Run full test suite before releases
- Track test failure rate over time

## CI/CD Workflow

```
Developer Push → GitHub Actions
    ↓
Test Job (Linux, macOS, Windows)
    ├─ go test ./...
    ├─ Coverage report
    └─ Lint checks
    ↓
[If Tests Pass]
    ├─ Build Job
    │   ├─ Compile for Linux, macOS, Windows
    │   └─ Create checksums
    └─ Coverage Report
        └─ Upload to Codecov
    ↓
[If Tag Push]
    └─ Release Job
        ├─ Create GitHub Release
        ├─ Upload binaries
        └─ Create changelog
```

## Definition of Done

- [ ] All unit tests implemented and passing (>85% coverage)
- [ ] Integration tests implemented and passing
- [ ] Benchmark suite created and baselines recorded
- [ ] CI/CD pipeline configured and working
- [ ] Code coverage reports generated
- [ ] Documentation complete (TESTING.md, etc.)
- [ ] All tests pass on main branch
- [ ] Binary releases created
- [ ] Contributing guide updated

## Dependencies

- Milestone 1.1 (Repomap Core) – core code to test
- Milestone 1.2 (CLI Framework) – framework code to test
- GitHub Actions (or equivalent CI/CD platform)
- Go 1.19+

## Timeline Estimate

- **Duration:** 1 week
- **Parallel with:** Milestones 1.1 & 1.2 (testing starts partway through)

## Risks

| Risk | Mitigation |
|------|-----------|
| Difficult to achieve 85%+ coverage | Focus on critical paths first; use code coverage tools to identify gaps |
| CI/CD setup complexity | Start simple; iterate based on needs |
| Test fixtures are outdated | Maintain fixtures; update with new patterns |
| Performance regression not caught | Set baselines early; alert on regressions |

## Success Metrics

| Metric | Target |
|--------|--------|
| Code coverage | >85% |
| Test pass rate | 100% |
| CI/CD uptime | >99% |
| Release frequency | Weekly (or on-demand) |
| Bug detection | Major bugs caught before release |

## Next Phase

Milestone 1.3 completes Phase 1. Phase 2 will:
- Add multi-language support to repomap (Tree-sitter)
- Build `semgrep-cli` using the CLI framework
- Extend testing for new tools

---

## Key Deliverables Summary

1. **Test Suites** – Comprehensive tests for repomap and framework
2. **Performance Benchmarks** – Track metrics over time
3. **CI/CD Pipeline** – Automated testing and releases
4. **Documentation** – Testing procedures and guidelines
