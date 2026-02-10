package util

import (
	"path/filepath"
	"strings"
)

// FilterByExtension returns a new slice containing only paths that match one of the given extensions.
// Extensions should include the dot (e.g., ".go").
func FilterByExtension(paths []string, exts []string) []string {
	if len(exts) == 0 {
		return paths
	}

	var filtered []string
	for _, path := range paths {
		ext := filepath.Ext(path)
		for _, e := range exts {
			if strings.EqualFold(ext, e) {
				filtered = append(filtered, path)
				break
			}
		}
	}
	return filtered
}

// ExcludeByPattern returns a new slice containing only paths that DO NOT match any of the given glob patterns.
func ExcludeByPattern(paths []string, patterns []string) []string {
	if len(patterns) == 0 {
		return paths
	}

	var filtered []string
	for _, path := range paths {
		matched := false
		for _, pattern := range patterns {
			// filepath.Match handles glob patterns
			if m, err := filepath.Match(pattern, filepath.Base(path)); err == nil && m {
				matched = true
				break
			}
			// Also check full path if it contains slashes?
			// Standard .gitignore logic is complex.
			// Task says "ExcludeByPattern".
			// Let's assume simple glob matching on the name or relative path.
			// filepath.Match only matches against the name unless pattern contains separators.

			// Let's try matching against the full path if strictly needed,
			// but for simple file filtering, usually we match against basename.
			// Or we can try matching the full path relative to root?
			// Let's stick to standard filepath.Match behavior on the basename for safety first,
			// or if pattern has slash, match whole path.

			// Improving matching logic:
			if strings.Contains(pattern, "/") || strings.Contains(pattern, "\\") {
				if m, err := filepath.Match(pattern, path); err == nil && m {
					matched = true
					break
				}
			}
		}

		if !matched {
			filtered = append(filtered, path)
		}
	}
	return filtered
}
