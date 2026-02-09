# Milestone 2.2: Semgrep-CLI – Tasks

## Task 1: Evaluate Embedding Models
- **Objective:** Test and select lightweight embedding model for code
- **Candidates:** all-MiniLM-L6-v2, CodeBERT, other code-specific models
- **Deliverable:** Evaluation report with accuracy, latency, memory metrics
- **Acceptance Criteria:** Model selected, justified by performance

## Task 2: Set Up Embedding Infrastructure
- **Objective:** Create `pkg/embeddings/` package with model loader
- **Steps:** Model download, caching, initialization logic
- **Deliverable:** embeddings/model.go with model management
- **Acceptance Criteria:** Model loads and initializes without error

## Task 3: Implement Code Snippet Vectorization
- **Objective:** Build `pkg/embeddings/vectorizer.go` to embed code snippets
- **Steps:** Tokenization, embedding computation, vector normalization
- **Deliverable:** Vectorizer with batch processing support
- **Acceptance Criteria:** Produces 384-dim vectors for any code snippet

## Task 4: Create Vector Index Structure
- **Objective:** Design `pkg/search/index.go` for vector storage and retrieval
- **Index Format:** In-memory with optional disk persistence
- **Deliverable:** Index with add(), search(), persist() methods
- **Acceptance Criteria:** Fast similarity search, <100ms for 10k vectors

## Task 5: Implement Vector Similarity Search
- **Objective:** Build cosine similarity search in index
- **Algorithm:** Cosine distance, k-nearest neighbors
- **Deliverable:** search/cosine.go with similarity metrics
- **Acceptance Criteria:** Returns top-k results ranked by similarity

## Task 6: Build Query Embedding Pipeline
- **Objective:** Create pipeline to embed natural language queries
- **Steps:** Normalize query, apply same tokenization as code, compute embedding
- **Deliverable:** query.go with query embedding logic
- **Acceptance Criteria:** Queries embedded consistently with code

## Task 7: Create Code Index Builder
- **Objective:** Implement indexing pipeline from repomap output
- **Steps:** Extract definitions from repomap, tokenize, embed, index
- **Deliverable:** builder.go with full indexing pipeline
- **Acceptance Criteria:** Index built from repomap in <5 seconds for 10k LOC

## Task 8: Implement Result Ranking Strategy
- **Objective:** Create ranking beyond vector similarity
- **Factors:** Similarity score, definition type (class vs function), language, popularity
- **Deliverable:** ranker.go with multi-factor ranking
- **Acceptance Criteria:** Results ranked intuitively for queries

## Task 9: Create semgrep-cli Command
- **Objective:** Build `cmd/semgrep-cli/main.go` with CLI interface
- **Flags:** --root, --query, --output-format, --top-k
- **Deliverable:** CLI tool executable
- **Acceptance Criteria:** Tool accepts queries and produces output

## Task 10: Implement JSON Output Format
- **Objective:** Format search results as structured JSON
- **Fields:** File, line, definition, snippet, similarity score
- **Deliverable:** output.go with JSON serialization
- **Acceptance Criteria:** Output parseable by downstream tools

## Task 11: Implement Text Output Format
- **Objective:** Create human-readable text output
- **Format:** Ranked list with scores and snippets
- **Deliverable:** Text formatter
- **Acceptance Criteria:** Clear, readable output for terminal

## Task 12: Add Query Expansion
- **Objective:** Enhance query with synonyms and related terms
- **Strategy:** Extract keywords, apply thesaurus, re-rank results
- **Deliverable:** query_expansion.go
- **Acceptance Criteria:** Expanded queries improve relevance

## Task 13: Implement Result Filtering
- **Objective:** Add language, file-type, and threshold filters
- **Flags:** --lang python, --min-score 0.7, --exclude-test
- **Deliverable:** Filter logic integrated into search pipeline
- **Acceptance Criteria:** Filters applied correctly

## Task 14: Build Index Persistence
- **Objective:** Support saving/loading indexes from disk
- **Format:** Binary format for fast loading
- **Deliverable:** persist.go with save() and load() methods
- **Acceptance Criteria:** Index loads in <500ms from disk

## Task 15: Create Integration Tests
- **Objective:** Test on real codebase with semantic queries
- **Test Cases:** 10+ semantic queries with expected results
- **Deliverable:** Integration tests
- **Acceptance Criteria:** >90% query relevance

## Task 16: Benchmark Index Performance
- **Objective:** Measure index size, build time, query latency
- **Metrics:** Time vs repository size, memory usage
- **Deliverable:** Benchmark report
- **Acceptance Criteria:** Query response <500ms for 50k LOC

## Task 17: Add Language-Specific Ranking
- **Objective:** Rank results higher for specified languages
- **Integration:** Use Phase 2.1 language information
- **Deliverable:** Language aware ranker
- **Acceptance Criteria:** Results prioritize specified languages

## Task 18: Implement Batch Queries
- **Objective:** Support multiple queries in single invocation
- **Format:** JSON file with query list
- **Deliverable:** Batch query processor
- **Acceptance Criteria:** Processes 10+ queries in <2 seconds

## Task 19: Add Interactive Mode
- **Objective:** Build REPL for interactive queries
- **Commands:** query, load, save, clear, help
- **Deliverable:** Interactive CLI mode
- **Acceptance Criteria:** REPL responsive, commands work

## Task 20: Create Documentation and Examples
- **Objective:** Write user guide with query examples
- **Content:** Basic queries, advanced techniques, performance tips
- **Deliverable:** docs/SEMGREP_CLI.md with 10+ examples
- **Acceptance Criteria:** Clear guide with working examples

## Task 21: Implement Multi-Repository Search
- **Objective:** Search across multiple indexed repositories
- **Flag:** --repos repo1,repo2,repo3
- **Deliverable:** Multi-repo support
- **Acceptance Criteria:** Searches multiple indexes

## Task 22: Add Result Caching
- **Objective:** Cache recent query results for repeated queries
- **Strategy:** LRU cache of recent results
- **Deliverable:** Cache implementation
- **Acceptance Criteria:** Cached queries return instantly

## Task 23: Performance Tuning
- **Objective:** Optimize embedding computation and index search
- **Profiling:** Identify and fix bottlenecks
- **Deliverable:** Performance improvements documented
- **Acceptance Criteria:** 20%+ faster than baseline

## Task 24: Error Handling and Recovery
- **Objective:** Handle corrupted indexes, missing models, invalid queries
- **Deliverable:** Robust error handling throughout
- **Acceptance Criteria:** Clear error messages, graceful degradation

## Task 25: Code Review and Final Testing
- **Objective:** Final QA and code review
- **Deliverable:** All tests passing, code review approved
- **Acceptance Criteria:** >85% code coverage, no critical issues
