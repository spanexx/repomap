# Milestone 2.1: Tree-Sitter Integration – Tasks

## Task 1: Evaluate Tree-Sitter Bindings
- **Objective:** Select best Go Tree-sitter binding for stability and performance
- **Research:** `go-tree-sitter` vs `tree-sitter-go` vs manual bindings
- **Deliverable:** Evaluation report with performance benchmarks
- **Acceptance Criteria:** Binding chosen, performance metrics documented

## Task 2: Set Up Tree-Sitter Dependency
- **Objective:** Add Tree-sitter library to go.mod with language definitions
- **Steps:** Download language binaries, create vendor directory, update go.mod
- **Deliverable:** go.mod updated, build succeeds
- **Acceptance Criteria:** All language binaries load without error

## Task 3: Create Tree-Sitter Wrapper Package
- **Objective:** Build `pkg/treesitter/parser.go` with Tree-sitter abstractions
- **Includes:** Parser initialization, AST navigation, query builder
- **Deliverable:** parser.go with >80% test coverage
- **Acceptance Criteria:** Parse Go, Python, JavaScript successfully

## Task 4: Implement Language Registry
- **Objective:** Create `pkg/treesitter/language.go` for language detection and routing
- **Steps:** Language mime-type detection, extension mapping, parser routing
- **Deliverable:** language.go with language detection logic
- **Acceptance Criteria:** Correctly identifies 7 languages from file path

## Task 5: Create Go Extractor
- **Objective:** Implement `pkg/languages/go/extractor.go` for Go definitions
- **Steps:** Extract functions, methods, structs, interfaces using Tree-sitter queries
- **Deliverable:** Extractor handles all Go constructs
- **Acceptance Criteria:** Passes test on real Go codebase (stdlib)

## Task 6: Create Python Extractor
- **Objective:** Implement `pkg/languages/python/extractor.go`
- **Constructs:** Classes, functions, decorators, async functions
- **Deliverable:** Python extractor with query templates
- **Acceptance Criteria:** Works on Django, Flask, Pandas codebases

## Task 7: Create JavaScript Extractor
- **Objective:** Implement `pkg/languages/javascript/extractor.go`
- **Constructs:** Functions, classes, exports, async, arrow functions
- **Deliverable:** JavaScript extractor
- **Acceptance Criteria:** Works on Node.js and browser code

## Task 8: Create TypeScript Extractor
- **Objective:** Implement `pkg/languages/typescript/extractor.go`
- **Constructs:** Interfaces, types, classes, functions, generics
- **Deliverable:** TypeScript extractor
- **Acceptance Criteria:** Works on real TypeScript projects

## Task 9: Create Rust Extractor
- **Objective:** Implement `pkg/languages/rust/extractor.go`
- **Constructs:** Functions, structs, traits, enums, macros
- **Deliverable:** Rust extractor
- **Acceptance Criteria:** Works on real Rust projects

## Task 10: Create Java Extractor
- **Objective:** Implement `pkg/languages/java/extractor.go`
- **Constructs:** Classes, interfaces, methods, annotations
- **Deliverable:** Java extractor
- **Acceptance Criteria:** Works on Spring Boot projects

## Task 11: Create C++ Extractor
- **Objective:** Implement `pkg/languages/cpp/extractor.go`
- **Constructs:** Classes, functions, namespaces, templates
- **Deliverable:** C++ extractor
- **Acceptance Criteria:** Works on real C++ codebases

## Task 12: Create Generalized Extractor Interface
- **Objective:** Build `pkg/extractor/extractor.go` unifying all language extractors
- **Interface:** GetDefinitions(), GetImports(), GetClasses(), etc.
- **Deliverable:** Extractor interface implemented by all language extractors
- **Acceptance Criteria:** Single interface works for all 7 languages

## Task 13: Update Repomap Parser Integration
- **Objective:** Refactor repomap parsing to use Tree-sitter instead of ast/goparser
- **Steps:** Replace language-specific parsing logic with Tree-sitter calls
- **Deliverable:** repomap works with Tree-sitter backend
- **Acceptance Criteria:** All Phase 1 tests pass with new backend

## Task 14: Test Go Language Support
- **Objective:** Unit test Go extractor on Go stdlib
- **Test Files:** 5+ Go packages with all construct types
- **Deliverable:** Unit tests with >95% pass rate
- **Acceptance Criteria:** All expected definitions found

## Task 15: Test Python Language Support
- **Objective:** Unit test Python extractor on real Python projects
- **Test Files:** Django, Flask, Pandas samples
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Definitions match expected output

## Task 16: Test JavaScript Language Support
- **Objective:** Unit test JavaScript extractor
- **Test Files:** Node.js, React, vanilla JS samples
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Finds all exports and classes

## Task 17: Test TypeScript Language Support
- **Objective:** Unit test TypeScript extractor
- **Test Files:** Real TypeScript projects, interface definitions
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Handles interfaces and generics

## Task 18: Test Rust Language Support
- **Objective:** Unit test Rust extractor
- **Test Files:** Tokio, Serde samples
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Finds traits and macros

## Task 19: Test Java Language Support
- **Objective:** Unit test Java extractor
- **Test Files:** Spring Boot sample, JDK stdlib
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Finds annotations and inheritance

## Task 20: Test C++ Language Support
- **Objective:** Unit test C++ extractor
- **Test Files:** Real C++ projects with templates
- **Deliverable:** Integration tests
- **Acceptance Criteria:** Handles templates and namespaces

## Task 21: Add CLI Language Flags
- **Objective:** Add `--include-lang` and `--exclude-lang` to repomap CLI
- **Flags:** `repomap --include-lang go,python,rust`
- **Deliverable:** CLI updated with language filtering
- **Acceptance Criteria:** Flags work correctly, filters applied

## Task 22: Implement Language-Aware Ranking
- **Objective:** Update ranking strategy in Phase 1 to rank by language relevance
- **Strategy:** Prefer definitions in specified languages
- **Deliverable:** Ranking algorithm updated
- **Acceptance Criteria:** Results ranked by language precedence

## Task 23: Performance Benchmarking
- **Objective:** Benchmark Tree-sitter vs old ast-based parsing
- **Metrics:** Parse time, memory usage for multi-language repos
- **Deliverable:** Benchmark report
- **Acceptance Criteria:** No regression, <5% overhead

## Task 24: Documentation – Language Support Guide
- **Objective:** Write guide on supported languages and extractors
- **Content:** Language-specific instructions, example outputs, known limitations
- **Deliverable:** docs/LANGUAGE_SUPPORT.md
- **Acceptance Criteria:** Complete guide with examples

## Task 25: Final Integration Tests
- **Objective:** End-to-end tests with multi-language repository
- **Test Repo:** Repository with all 7 languages
- **Deliverable:** Integration tests passing
- **Acceptance Criteria:** All 7 languages detected and parsed correctly
