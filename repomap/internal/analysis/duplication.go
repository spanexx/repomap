package analysis

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"go/scanner"
	"go/token"
	"strings"

	"github.com/spanexx/agents-cli/repomap/internal/output"
)

// DuplicationDetector scans files for duplicate code blocks.
type DuplicationDetector struct {
	MinTokens int // Minimum tokens to consider a block
}

// NewDuplicationDetector creates a new detector with default settings.
func NewDuplicationDetector() *DuplicationDetector {
	return &DuplicationDetector{
		MinTokens: 50, // Approx 10-15 lines of code
	}
}

// Tokenize processes a file strictly for structural hashing (ignoring comments/whitespace).
func (d *DuplicationDetector) Tokenize(content []byte) ([]string, error) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(content))

	s.Init(file, content, nil, 0) // No comments

	var tokens []string
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		// We only care about the token type for structure, or literal for exact match?
		// For stricter "Copy-Paste" detection, include literals.
		// For "Structural" similarity, maybe just types.
		// Let's go with exact token values for now (Type + Literal).
		// But normalize whitespace.
		if lit == "" {
			lit = tok.String()
		}
		tokens = append(tokens, lit)
	}
	return tokens, nil
}

// Analyze checks a set of files for duplicates.
// This is a simplified approach: exact sliding window hash match across all files.
func (d *DuplicationDetector) Analyze(files map[string][]byte) ([]output.Issue, error) {
	// 1. Tokenize all files
	fileTokens := make(map[string][]string)
	for path, content := range files {
		toks, err := d.Tokenize(content)
		if err != nil {
			continue // Skip unparseable
		}
		if len(toks) < d.MinTokens {
			continue
		}
		fileTokens[path] = toks
	}

	// 2. Build Hash Map
	// Key: Hash of N tokens
	// Value: List of (File, StartIndex)
	blockMap := make(map[string][]struct {
		File  string
		Start int
	})

	windowSize := d.MinTokens

	for path, tokens := range fileTokens {
		for i := 0; i <= len(tokens)-windowSize; i++ {
			// Compute hash of window
			window := strings.Join(tokens[i:i+windowSize], "|")
			hash := sha256.Sum256([]byte(window))
			hashStr := hex.EncodeToString(hash[:])

			blockMap[hashStr] = append(blockMap[hashStr], struct {
				File  string
				Start int
			}{path, i})
		}
	}

	// 3. Identify Duplicates
	var issues []output.Issue
	// To avoid overlapping reportspam, we need to merge contiguous blocks.
	// This is complex. For MVP, we'll just report "Found duplication involving X files".
	// Or we report distinct blocks.

	// Set of reported (File, Start) to avoid redundant noise
	reported := make(map[string]bool)

	for _, occurrences := range blockMap {
		if len(occurrences) > 1 {
			// Found duplicate!
			// Check if we already reported a block covering this
			// (Simplification: just take first occurrence and report it if not overlaps drastically)

			// Group by file to see cross-file vs intra-file
			filesInvolved := make(map[string]bool)
			for _, occ := range occurrences {
				filesInvolved[occ.File] = true
			}

			if len(filesInvolved) > 1 {
				// Cross-file duplication
				desc := fmt.Sprintf("Found duplicate code block (%d tokens) across %d files: ", windowSize, len(filesInvolved))
				var paths []string
				for f := range filesInvolved {
					paths = append(paths, f)
				}
				desc += strings.Join(paths, ", ")

				// De-duplicate the issue report string itself
				if !reported[desc] {
					issues = append(issues, output.Issue{
						Type:        "duplication",
						Severity:    "medium",
						Description: desc,
					})
					reported[desc] = true
				}
			}
		}
	}

	return issues, nil
}
