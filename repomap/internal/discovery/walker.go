package discovery

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spanexx/agents-cli/repomap/internal/parsing"
)

// Walk traverses the directory tree rooted at root and returns a list of files
// that match the default filtering criteria (Go files, non-binary, non-hidden)
// and respect .gitignore rules.
func Walk(root string) ([]string, error) {
	var files []string

	// Get supported extensions from the parsing registry
	supportedExts := make(map[string]bool)
	for _, ext := range parsing.DefaultRegistry.SupportedExtensions() {
		supportedExts[ext] = true
	}

	// Binary extensions to exclude
	excludeExts := map[string]bool{
		".exe":   true,
		".o":     true,
		".a":     true,
		".so":    true,
		".dylib": true,
		".dll":   true,
		".bin":   true,
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

		// Skip hidden directories (starting with .) or node_modules
		if d.IsDir() {
			name := d.Name()
			if (strings.HasPrefix(name, ".") && name != ".") || name == "node_modules" || name == "vendor" || name == "dist" || name == "build" {
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

		// Include only allowed extensions from the registry
		if supportedExts[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
