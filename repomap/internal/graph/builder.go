package graph

import (
	"strings"
	"sync"
)

// Graph represents a directed graph of file dependencies.
type Graph struct {
	// Nodes maps file paths to their node information.
	Nodes map[string]*Node
	// Edges maps a source file path to a list of destination file paths (imports).
	Edges map[string][]string
}

// Node represents a file in the import graph.
type Node struct {
	Path     string
	Imports  []string
	InDegree int
}

// Builder constructs a dependency graph.
type Builder struct {
	mu    sync.Mutex
	files map[string][]string // file path -> list of raw imports
}

// NewBuilder creates a new graph builder.
func NewBuilder() *Builder {
	return &Builder{
		files: make(map[string][]string),
	}
}

// AddFile adds a file and its imports to the builder.
// path should be relative to the repository root.
// imports is a list of imported package paths.
func (b *Builder) AddFile(path string, imports []string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.files[path] = imports
}

// Build constructs the final graph.
func (b *Builder) Build(moduleName string) *Graph {
	b.mu.Lock()
	defer b.mu.Unlock()

	g := &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}

	// Index files by directory (package path)
	// Map: directory path -> list of files
	pkgFiles := make(map[string][]string)

	for path, imports := range b.files {
		// Use slash for consistency in graph keys
		normalizedPath := strings.ReplaceAll(path, "\\", "/")

		dir := "."
		lastSlash := strings.LastIndex(normalizedPath, "/")
		if lastSlash != -1 {
			dir = normalizedPath[:lastSlash]
		}
		pkgFiles[dir] = append(pkgFiles[dir], normalizedPath)

		// Initialize node with normalized path
		g.Nodes[normalizedPath] = &Node{
			Path:    normalizedPath,
			Imports: imports,
		}
	}

	// Build edges
	for srcFileRaw, imports := range b.files {
		srcFile := strings.ReplaceAll(srcFileRaw, "\\", "/")

		for _, imp := range imports {
			// Normalize import: strip module name if present
			targetPkg := imp
			if moduleName != "" && strings.HasPrefix(imp, moduleName) {
				targetPkg = strings.TrimPrefix(imp, moduleName)
				targetPkg = strings.TrimPrefix(targetPkg, "/")
			}

			// If targetPkg is a directory in our index, add edges
			if destFiles, ok := pkgFiles[targetPkg]; ok {
				for _, destFile := range destFiles {
					// Avoid self-loops
					if srcFile == destFile {
						continue
					}

					g.Edges[srcFile] = append(g.Edges[srcFile], destFile)
					// Ensure destination node exists before updating (should exist from init loop)
					if node, ok := g.Nodes[destFile]; ok {
						node.InDegree++
					}
				}
			}
		}
	}

	return g
}
