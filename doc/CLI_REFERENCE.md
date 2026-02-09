# CLI Reference

Complete documentation of all repomap command-line options and usage patterns.

---

## Basic Syntax

```bash
repomap [FLAGS] [OPTIONS]
```

---

## Global Flags

### `--root <PATH>`

**Description:** Root directory of the repository to analyze.

**Required:** No (defaults to current directory `.`)

**Example:**
```bash
repomap --root ./src
repomap --root /home/user/projects/myapp
```

---

## Output Options

### `--output <FORMAT>`

**Description:** Output format for the generated repomap.

**Accepted Values:**
- `xml` - XML format (default)
- `json` - JSON format

**Default:** `xml`

**Examples:**
```bash
repomap --output xml
repomap --output json
```

### `--max-tokens <NUMBER>`

**Description:** Maximum token budget for the output. The tool will include files in rank order until this budget is exhausted.

**Type:** Positive integer

**Default:** 5000

**Examples:**
```bash
repomap --max-tokens 500      # Tight budget
repomap --max-tokens 2000     # Medium budget
repomap --max-tokens 10000    # Generous budget
```

**Notes:**
- Respects token counting based on the `cl100k_base` encoding (OpenAI)
- Output will indicate remaining token budget in a comment
- Files are included in importance order until budget exhausted

---

## Filtering Options

### `--include-ext <EXTENSIONS>`

**Description:** Only include files with specific extensions.

**Type:** Comma-separated list of extensions (with dots)

**Default:** All extensions

**Examples:**
```bash
repomap --include-ext .go
repomap --include-ext .go,.js
repomap --include-ext .py,.pyx,.pyi
```

### `--exclude-ext <EXTENSIONS>`

**Description:** Exclude files with specific extensions.

**Type:** Comma-separated list of extensions (with dots)

**Default:** None (all extensions included, except filtered by `.gitignore`)

**Examples:**
```bash
repomap --exclude-ext .test.js
repomap --exclude-ext .pyc,.pyo,.so
```

### `--ignore-tests`

**Description:** Exclude test files from the repomap.

**Type:** Boolean flag

**Default:** false

**Patterns Excluded:**
- `*_test.go`
- `*_test.py`
- `test_*.go`
- `test_*.py`
- Any file in `test/` or `tests/` directories
- Files matching `*Test.js`, `*.spec.js`, `*.test.js` patterns

**Examples:**
```bash
repomap --ignore-tests
repomap --root ./src --ignore-tests --output json
```

### `--ignore-vendor`

**Description:** Explicitly exclude vendor directories (usually handled by `.gitignore`).

**Type:** Boolean flag

**Default:** true (already handled by `.gitignore` parsing)

**Examples:**
```bash
repomap --ignore-vendor
```

---

## Behavior Flags

### `--verbose`

**Description:** Enable verbose output to stderr for debugging.

**Type:** Boolean flag

**Default:** false

**Output Includes:**
- File discovery statistics
- Parsing progress
- Graph building details
- Ranking information
- Token counting details

**Examples:**
```bash
repomap --verbose
repomap --root . --max-tokens 2000 --verbose
```

### `--follow-symlinks`

**Description:** Follow symbolic links during directory traversal.

**Type:** Boolean flag

**Default:** false (symlinks are ignored for safety)

**Examples:**
```bash
repomap --follow-symlinks
```

### `--gitignore-file <PATH>`

**Description:** Specify a custom `.gitignore` file instead of using the repository's.

**Type:** File path

**Default:** `.gitignore` in the root directory

**Examples:**
```bash
repomap --gitignore-file ./custom-ignore
repomap --root /external/repo --gitignore-file /external/repo/.gitignore
```

---

## Language-Specific Options

### `--include-lang <LANGUAGES>`

**Description:** Only include files of specific programming languages (future feature for Tree-sitter support).

**Type:** Comma-separated list of language identifiers

**Default:** All supported languages

**Supported Languages (future):**
- `go`, `python`, `javascript`, `typescript`, `rust`, `java`, `c`, `cpp`

**Examples:**
```bash
repomap --include-lang go
repomap --include-lang python,javascript
```

### `--exclude-lang <LANGUAGES>`

**Description:** Exclude files of specific programming languages.

**Type:** Comma-separated list of language identifiers

**Default:** None

**Examples:**
```bash
repomap --exclude-lang test
repomap --exclude-lang vendor
```

---

## Performance Options

### `--max-files <NUMBER>`

**Description:** Maximum number of files to process (useful for testing).

**Type:** Positive integer

**Default:** Unlimited

**Examples:**
```bash
repomap --max-files 100     # Process only first 100 files
repomap --max-files 1000
```

### `--parallel <NUMBER>`

**Description:** Number of parallel workers for parsing files (future optimization).

**Type:** Positive integer (1-32)

**Default:** Auto-detect (based on CPU cores)

**Examples:**
```bash
repomap --parallel 4
repomap --parallel 1        # Single-threaded
```

---

## Help and Version

### `--help` or `-h`

**Description:** Display help message and exit.

**Examples:**
```bash
repomap --help
repomap -h
```

### `--version` or `-v`

**Description:** Display version information and exit.

**Examples:**
```bash
repomap --version
repomap -v
```

---

## Command Examples

### Example 1: Basic Usage

```bash
repomap --root . --output xml
```

Analyzes current directory, outputs XML to stdout.

### Example 2: JSON Output with Token Budget

```bash
repomap --root ./src --output json --max-tokens 2000
```

Analyzes `./src`, outputs JSON with 2000-token budget.

### Example 3: Go Files Only

```bash
repomap --root . --include-ext .go --output xml
```

Analyzes only Go files in current directory.

### Example 4: Exclude Test Files and Vendor

```bash
repomap --root . --ignore-tests --ignore-vendor --output json
```

Analyzes repository excluding tests and vendor directories.

### Example 5: Verbose Output with Filtering

```bash
repomap --root . --include-lang go --verbose --output xml
```

Analyzes Go files with verbose debug output.

### Example 6: Save Output to File

```bash
repomap --root . --max-tokens 5000 --output xml > repomap.xml
```

Saves XML repomap to `repomap.xml`.

### Example 7: Process with jq

```bash
repomap --root . --output json | jq '.repomap.files[0:5]'
```

Show first 5 files from JSON output.

### Example 8: Create Comparison Reports

```bash
# Tight budget
repomap --root . --max-tokens 500 --output json > tight.json

# Medium budget
repomap --root . --max-tokens 2000 --output json > medium.json

# Generous budget
repomap --root . --max-tokens 5000 --output json > generous.json
```

Generate multiple reports for comparison.

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid arguments, file not found, etc.) |
| 2 | Parse error (invalid syntax in source files) |
| 3 | Output error (cannot write to file) |
| 4 | Timeout error (processing took too long) |

---

## Error Messages

### "no files found"

**Cause:** No files matching criteria in the specified root.

**Solution:**
- Verify `--root` path exists
- Check that files aren't all filtered by `.gitignore`
- Try `--verbose` to see filtering details

### "invalid token budget"

**Cause:** `--max-tokens` value is not a valid positive integer.

**Solution:**
```bash
repomap --max-tokens 2000  # Valid
repomap --max-tokens abc   # Invalid
```

### "unknown output format"

**Cause:** `--output` value is not `xml` or `json`.

**Solution:**
```bash
repomap --output xml    # Valid
repomap --output yaml   # Invalid
```

---

## Configuration Files (Future Feature)

A `.repomaprc` or `.repomap.yml` file could specify defaults:

```yaml
# .repomap.yml
root: ./src
output: json
max-tokens: 2000
ignore-tests: true
include-ext:
  - .go
  - .mod
```

Usage would then be:
```bash
repomap  # Uses settings from .repomap.yml
```

