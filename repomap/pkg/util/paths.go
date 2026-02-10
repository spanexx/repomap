package util

import (
	"path/filepath"
	"strings"
)

// NormalizePath cleans the path and ensures forward slashes for cross-platform consistency.
func NormalizePath(path string) string {
	cleaned := filepath.Clean(path)
	return filepath.ToSlash(cleaned)
}

// IsAbsolutePath checks if the path is absolute.
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// MakeRelative returns the path relative to the base path.
// It normalizes the result to use forward slashes.
func MakeRelative(basePath, targetPath string) string {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		// If relative path cannot be determined, return the original target path
		// normalized.
		return NormalizePath(targetPath)
	}
	return filepath.ToSlash(rel)
}

// HasExtension checks if the file has one of the given extensions.
// Extensions should include the dot (e.g., ".go").
func HasExtension(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, e := range extensions {
		if strings.EqualFold(ext, e) {
			return true
		}
	}
	return false
}
