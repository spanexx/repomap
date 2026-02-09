package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

// Walk traverses the directory tree rooted at root and returns a list of files
// that match the default filtering criteria (Go files, non-binary, non-hidden)
// and respect .gitignore rules.
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

	// Parse .gitignore if it exists
	gitignore, err := ParseGitignore(root)
	// We ignore error here as ParseGitignore returns usable object even on error (empty)
	// or we can just proceed. Actually ParseGitignore returns nil on error.
	if err != nil {
		// If we can't parse gitignore (e.g. permission error), we proceed without it?
		// Or strictly fail? Let's proceed with empty one.
		gitignore = &Gitignore{root: root}
	}

	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check gitignore first
		if gitignore != nil && gitignore.Matches(path) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
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
