package parsing

import (
	"path/filepath"
	"strings"
)

// Extractor interface for language-specific definition extraction.
type Extractor interface {
	ExtractDefinitions(filePath string) ([]string, error)
	ExtractImports(filePath string) ([]string, error)
}

// Registry manages language-specific extractors.
type Registry struct {
	extractors map[string]Extractor
}

// NewRegistry creates a new extractor registry.
func NewRegistry() *Registry {
	return &Registry{
		extractors: make(map[string]Extractor),
	}
}

// Register adds an extractor for a given extension.
func (r *Registry) Register(ext string, extractor Extractor) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	r.extractors[ext] = extractor
}

// Get returns the extractor for the given file path.
func (r *Registry) Get(path string) Extractor {
	ext := strings.ToLower(filepath.Ext(path))
	if extractor, ok := r.extractors[ext]; ok {
		return extractor
	}
	return nil
}

// SupportedExtensions returns a list of extensions with registered extractors.
func (r *Registry) SupportedExtensions() []string {
	exts := make([]string, 0, len(r.extractors))
	for ext := range r.extractors {
		exts = append(exts, ext)
	}
	return exts
}

// DefaultRegistry is a global registry instance for convenience.
var DefaultRegistry = NewRegistry()
