# Architecture Overview

## The Funnel Concept

The tool needs to act as a **funnel**: taking a massive file system and reducing it to a high-density "map." This approach allows us to extract essential structural information from large codebases without overwhelming context windows.

## The Complete Pipeline

The architecture follows a five-stage pipeline:

### 1. **Discovery**
Walk the directory tree quickly to identify all files in the codebase.

### 2. **Filtering**
Apply `.gitignore` rules and exclude binary/irrelevant files to reduce noise.

### 3. **Parsing**
Use **Tree-sitter** (via Go bindings) to extract definitions (functions, classes, structs) without reading the whole file body. For the MVP, Go's standard `go/ast` can be used.

### 4. **Ranking**
Score files based on "centrality"—how many other files import them. Files that are imported frequently are more likely to be important for understanding the codebase.

### 5. **Rendering**
Generate the output tree (XML/JSON), cutting off details once the `max-tokens` budget is reached.

## Design Principles

- **Token Efficiency:** Only extract function/struct signatures, not full implementations
- **Context Optimization:** Prioritize files based on their importance to the overall system
- **Language Agnostic:** Support multiple programming languages through Tree-sitter
- **Fast Execution:** Use optimized libraries for file walking and parsing
- **Budgeted Output:** Respect token limits to fit within LLM context windows
