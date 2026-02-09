package discovery

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Gitignore represents a set of ignore patterns parsed from a .gitignore file.
type Gitignore struct {
	patterns []string
	root     string
}

// ParseGitignore reads the .gitignore file from the given root directory
// and returns a Gitignore object. If no .gitignore exists, it returns
// a Gitignore with no patterns and nil error.
func ParseGitignore(root string) (*Gitignore, error) {
	gitignorePath := filepath.Join(root, ".gitignore")

	// If file doesn't exist, return empty gitignore
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return &Gitignore{root: root, patterns: []string{}}, nil
	}

	file, err := os.Open(gitignorePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Gitignore{
		patterns: patterns,
		root:     root,
	}, nil
}

// Matches returns true if the given file path should be ignored.
// path should be absolute or relative to the execution context, but
// logic will normalize it relative to g.root.
func (g *Gitignore) Matches(path string) bool {
	if len(g.patterns) == 0 {
		return false
	}

	// Calculate relative path from root
	relPath, err := filepath.Rel(g.root, path)
	if err != nil {
		return false
	}

	// If file is outside root, don't ignore
	if strings.HasPrefix(relPath, "..") {
		return false
	}

	// Clean path for consistent matching
	relPath = filepath.ToSlash(relPath)

	ignored := false
	for _, pattern := range g.patterns {
		negate := false
		if strings.HasPrefix(pattern, "!") {
			negate = true
			pattern = strings.TrimPrefix(pattern, "!")
		}

		// Handle directory-only patterns (ending with /)
		dirOnly := strings.HasSuffix(pattern, "/")
		cleanPattern := pattern
		if dirOnly {
			cleanPattern = strings.TrimSuffix(pattern, "/")
		}

		// Determine if pattern is anchored
		// A pattern is anchored if it contains a slash (not counting trailing slash)
		// exception: strict leading slash
		isAnchored := strings.Contains(cleanPattern, "/")

		if strings.HasPrefix(cleanPattern, "/") {
			cleanPattern = strings.TrimPrefix(cleanPattern, "/")
			isAnchored = true // Leading slash forces anchor
		}

		match := false

		if isAnchored {
			// Anchored matching: pattern must match the start of relPath

			// 1. Check exact match
			if cleanPattern == relPath {
				match = true
			}

			// 2. Check directory prefix match
			// e.g. /build matches build/foo.log
			if strings.HasPrefix(relPath, cleanPattern + "/") {
				match = true
			}

			// 3. Glob match on full path
			// e.g. docs/*.md matches docs/readme.md
			if m, _ := filepath.Match(cleanPattern, relPath); m {
				match = true
			}
		} else {
			// Unanchored matching: pattern matches any component of the path
			// e.g. *.log matches app.log, src/app.log
			// e.g. node_modules matches node_modules, src/node_modules

			// Check if any path component matches the pattern
			parts := strings.Split(relPath, "/")
			for _, part := range parts {
				if m, _ := filepath.Match(cleanPattern, part); m {
					match = true
					break
				}
			}
		}

		if match {
			if negate {
				ignored = false
			} else {
				ignored = true
			}
		}
	}

	return ignored
}
