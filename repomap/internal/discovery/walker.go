package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

// Walk traverses the directory tree rooted at root and returns a list of files
// that match the default filtering criteria (Go files, non-binary, non-hidden).
func Walk(root string) ([]string, error) {
	var files []string

	// Default extensions to include
	includeExts := map[string]bool{
		".go": true,
	}

	// Binary extensions to exclude
	excludeExts := map[string]bool{
		".exe": true,
		".o":   true,
		".a":   true,
		".so":  true,
		".dylib": true,
		".dll": true,
		".bin": true,
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories (starting with .)
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") && d.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip non-regular files
		if !d.Type().IsRegular() {
			return nil
		}

		ext := filepath.Ext(path)

		// Skip binary files
		if excludeExts[ext] {
			return nil
		}

		// Include only allowed extensions (default to .go)
		if includeExts[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
