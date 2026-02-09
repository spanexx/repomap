# Tasks: Milestone 1.3 – Testing Infrastructure & CI/CD

## Dependencies
- **Prerequisites:**
  - Milestone 1.1 (Repomap Core)
  - Milestone 1.2 (CLI Framework & Integration)
- **Provides for:**
  - Final Production Release
  - Phase 2 (Multi-language support) with high-confidence baseline.


## Unit Testing – Repomap Core

### Task 1.3.1: Discovery Package Tests
- [ ] Create `internal/discovery/walker_test.go`
- [ ] Test `Walk()` function with various directories
- [ ] Test `.gitignore` filtering with real patterns
- [ ] Test error handling (permission denied, etc.)
- [ ] Target: >90% coverage on discovery package

**Acceptance Criteria:**
- Unit tests pass: `go test ./internal/discovery`
- Coverage >90%
- All error paths tested

---

### Task 1.3.2: Parsing Package Tests
- [ ] Create `internal/parsing/extractor_test.go`
- [ ] Test `ExtractDefinitions()` on various Go files
- [ ] Test `ExtractImports()` on various Go files
- [ ] Test error handling (syntax errors, etc.)
- [ ] Create test fixtures in `internal/parsing/testdata/`
- [ ] Target: >90% coverage on parsing package

**Acceptance Criteria:**
- Unit tests pass: `go test ./internal/parsing`
- Coverage >90%
- Test fixtures for 10+ code samples

---

### Task 1.3.3: Graph Package Tests
- [ ] Create `internal/graph/builder_test.go`
- [ ] Test graph construction on various file sets
- [ ] Test in-degree calculation
- [ ] Test error handling
- [ ] Target: >85% coverage on graph package

**Acceptance Criteria:**
- Unit tests pass: `go test ./internal/graph`
- Coverage >85%
- Test with 10+ graph scenarios

---

### Task 1.3.4: Ranking Package Tests
- [ ] Create `internal/ranking/ranker_test.go`
- [ ] Test ranking algorithm on mock graphs
- [ ] Test normalization to [0.0, 1.0]
- [ ] Test importance level assignment
- [ ] Target: >90% coverage on ranking package

**Acceptance Criteria:**
- Unit tests pass: `go test ./internal/ranking`
- Coverage >90%
- Verify ranking correctness mathematically

---

### Task 1.3.5: Output Package Tests
- [ ] Create `internal/output/xml_test.go`
- [ ] Create `internal/output/json_test.go`
- [ ] Test XML/JSON rendering with various data
- [ ] Test token budget enforcement
- [ ] Verify output validity (XML well-formed, JSON valid)
- [ ] Target: >85% coverage on output package

**Acceptance Criteria:**
- Unit tests pass: `go test ./internal/output`
- Coverage >85%
- Output is valid XML/JSON

---

## Unit Testing – CLI Framework

### Task 1.3.6: CLI Package Tests
- [ ] Create `pkg/cli/cli_test.go`, `flags_test.go`, `help_test.go`
- [ ] Test flag parsing with various combinations
- [ ] Test help text generation
- [ ] Test error handling for invalid flags
- [ ] Target: >85% coverage on cli package

**Acceptance Criteria:**
- Unit tests pass: `go test ./pkg/cli`
- Coverage >85%
- All flag types tested

---

### Task 1.3.7: Output Package Tests (Framework)
- [ ] Create `pkg/output/writer_test.go`, formatter tests
- [ ] Test XML, JSON, text formatters
- [ ] Test factory pattern
- [ ] Test token counting integration
- [ ] Target: >85% coverage on output package

**Acceptance Criteria:**
- Unit tests pass: `go test ./pkg/output`
- Coverage >85%
- All formatters produce valid output

---

### Task 1.3.8: Config Package Tests
- [ ] Create `pkg/config/config_test.go`
- [ ] Test config file loading
- [ ] Test environment variable reading
- [ ] Test configuration hierarchy
- [ ] Create test fixtures for config files
- [ ] Target: >85% coverage on config package

**Acceptance Criteria:**
- Unit tests pass: `go test ./pkg/config`
- Coverage >85%
- Test all levels of hierarchy

---

### Task 1.3.9: Error & Utility Package Tests
- [ ] Create `pkg/errors/errors_test.go`
- [ ] Create `pkg/util/*_test.go` for all utilities
- [ ] Test error type construction
- [ ] Test path operations, token counting, filtering
- [ ] Target: >80% coverage on all packages

**Acceptance Criteria:**
- Unit tests pass: `go test ./pkg/errors ./pkg/util`
- Coverage >80%

---

## Integration Testing

### Task 1.3.10: Test Fixtures Setup
- [ ] Create `test/fixtures/` directory
- [ ] Create small Go project fixture (10-20 files)
  - `test/fixtures/simple-go-app/`
  - Includes various import patterns and structures
- [ ] Create medium Go project fixture (100-200 files)
  - `test/fixtures/medium-go-project/`
- [ ] Create large Go project fixture (1K+ files)
  - `test/fixtures/large-go-project/` (or reference external repo)
- [ ] Document fixtures in `test/fixtures/README.md`

**Acceptance Criteria:**
- Fixtures exist and are usable
- Each fixture has clear purpose
- Documented for test authors

---

### Task 1.3.11: Full Pipeline Integration Tests
- [ ] Create `test/integration_test.go`
- [ ] Test complete repomap flow on small fixture
- [ ] Test on medium fixture
- [ ] Test on large fixture
- [ ] Verify output correctness for each fixture
- [ ] Test various `.gitignore` patterns

**Acceptance Criteria:**
- Integration tests pass: `go test ./test`
- All fixtures tested
- Output validated for correctness

---

### Task 1.3.12: CLI Integration Tests
- [ ] Test repomap CLI with various flags
- [ ] Test error handling on invalid inputs
- [ ] Test output format switching (--output xml/json)
- [ ] Test token budget enforcement
- [ ] Test configuration file loading

**Acceptance Criteria:**
- CLI integration tests pass
- All major CLI paths tested
- Error handling verified

---

### Task 1.3.13: Cross-Tool Integration Tests
- [ ] Test CLI framework with mock tools
- [ ] Test output writers with multiple data types
- [ ] Test configuration hierarchy across tools
- [ ] Verify no coupling between tools

**Acceptance Criteria:**
- Framework works with multiple tool types
- No breakage when using different tools

---

## Performance Benchmarking

### Task 1.3.14: Discovery Benchmarks
- [ ] Create `internal/discovery/walker_bench_test.go`
- [ ] Benchmark on 1K, 10K, 100K file sets
- [ ] Measure time and memory usage
- [ ] Establish baselines

**Acceptance Criteria:**
- Benchmarks run: `go test -bench=. ./internal/discovery`
- Results recorded
- Baselines established

---

### Task 1.3.15: Parsing Benchmarks
- [ ] Create `internal/parsing/extractor_bench_test.go`
- [ ] Benchmark definition extraction
- [ ] Benchmark import extraction
- [ ] Test on files of various sizes

**Acceptance Criteria:**
- Benchmarks run and complete
- Baselines recorded

---

### Task 1.3.16: End-to-End Benchmarks
- [ ] Create `test/benchmarks_test.go`
- [ ] Benchmark full repomap pipeline
- [ ] Test on 1K, 10K, 100K file sets
- [ ] Measure discovery, parsing, ranking, output phases
- [ ] Record memory usage

**Acceptance Criteria:**
- Full pipeline benchmarks recorded
- Performance acceptable (<30s for 100K files)
- Memory usage reasonable

---

### Task 1.3.17: Benchmark Tracking Script
- [ ] Create `scripts/compare_benchmarks.sh`
- [ ] Script to run benchmarks and compare to baseline
- [ ] Generate benchmark report
- [ ] Alert on performance regressions

**Acceptance Criteria:**
- Script runs and generates reports
- Regressions are detectable

---

## CI/CD Pipeline

### Task 1.3.18: GitHub Actions – Test Workflow
- [ ] Create `.github/workflows/test.yml`
- [ ] Test on Linux, macOS, Windows
- [ ] Test on Go 1.19, 1.20, 1.21+
- [ ] Run: `go test ./...`, coverage checks
- [ ] Block merges on test failure

**Acceptance Criteria:**
- Workflow runs on push/PR
- Tests pass on all platforms
- Coverage reports generated

---

### Task 1.3.19: GitHub Actions – Build Workflow
- [ ] Create `.github/workflows/build.yml`
- [ ] Build binaries for Linux, macOS, Windows
- [ ] Create checksums (SHA256)
- [ ] Publish as artifacts

**Acceptance Criteria:**
- Binaries built successfully
- Checksums created

---

### Task 1.3.20: GitHub Actions – Release Workflow
- [ ] Create `.github/workflows/release.yml`
- [ ] Triggered on git tag (v1.0.0, etc.)
- [ ] Create GitHub Release
- [ ] Upload binaries
- [ ] Generate changelog

**Acceptance Criteria:**
- Release workflow runs on tag
- Release created on GitHub
- Binaries published

---

### Task 1.3.21: Code Coverage Integration
- [ ] Integrate with Codecov or Coveralls
- [ ] Generate coverage reports in CI
- [ ] Publish coverage badges
- [ ] Set minimum coverage threshold (85%)

**Acceptance Criteria:**
- Coverage reports published
- Badge displays on README

---

## Documentation

### Task 1.3.22: Testing Documentation
- [ ] Create `TESTING.md`
- [ ] Document testing philosophy
- [ ] How to run tests locally
- [ ] How to write new tests
- [ ] Test naming conventions
- [ ] Coverage targets

**Acceptance Criteria:**
- Documentation is comprehensive
- New contributors can write tests

---

### Task 1.3.23: Benchmarking Documentation
- [ ] Create `BENCHMARKING.md`
- [ ] How to run benchmarks
- [ ] How to interpret results
- [ ] How to track performance over time
- [ ] Current performance baselines

**Acceptance Criteria:**
- Clear instructions for running benchmarks
- Baselines documented

---

### Task 1.3.24: CI/CD Documentation
- [ ] Create `CI_CD.md`
- [ ] Document workflow structure
- [ ] How to add new workflows
- [ ] Troubleshooting guide
- [ ] Release procedure

**Acceptance Criteria:**
- CI/CD pipeline is fully documented
- Release procedure is clear

---

### Task 1.3.25: Contributing Guide
- [ ] Update `CONTRIBUTING.md`
- [ ] Testing requirements
- [ ] Code review process
- [ ] Release process
- [ ] Code of conduct

**Acceptance Criteria:**
- Contributors know what to do
- Standards are clear

---

## Validation & Finalization

### Task 1.3.26: Run Full Test Suite
- [ ] Run all tests locally: `go test ./...`
- [ ] Verify coverage: `go test -cover ./...`
- [ ] Run benchmarks: `go test -bench=. ./...`
- [ ] Verify CI/CD passes

**Acceptance Criteria:**
- All tests pass
- Coverage >85%
- No errors

---

### Task 1.3.27: Documentation Review
- [ ] Review all test documentation
- [ ] Ensure clarity and completeness
- [ ] Add examples where needed

**Acceptance Criteria:**
- Documentation is polished
- Ready for publication

---

## Summary

**Total Tasks:** 27
**Estimated Effort:** 1 week
**Dependencies:** Milestones 1.1 & 1.2 (mostly complete)
**Success:** All tests passing + CI/CD working + benchmarks recorded + documentation complete
