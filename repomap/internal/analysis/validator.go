package analysis

import (
	"fmt"
	"strings"

	"github.com/spanexx/agents-cli/repomap/internal/output"
)

// DependencyGraph represents the import relationships
type DependencyGraph struct {
	Nodes map[string]*output.FileNode
	Edges map[string][]string // Path -> Imported Paths
}

// BuildGraph constructs a dependency graph from file nodes.
// It resolves imports to file paths where possible.
// Note: Go imports are package-based, repomap nodes are file-based.
// This is an approximation. Ideally we map packages.
// For MVP, we'll try to map import paths to local directories.
func BuildGraph(files []*output.FileNode, rootModule string) *DependencyGraph {
	g := &DependencyGraph{
		Nodes: make(map[string]*output.FileNode),
		Edges: make(map[string][]string),
	}

	// Index files by package path (directory)?
	// Or just keep them as files.
	// Imports are typically "github.com/user/repo/pkg/foo".
	// We need to match that to local files.
	// If rootModule is provided (e.g. from go.mod), we can strip it.

	// 1. Index all files by their likely import path
	// We don't have the go.mod module name easily here unless passed.
	// Let's assume absolute paths or relative paths in FileNode.
	// FileNode.Path is relative to root.

	// Quick hack: Index by directory
	dirFiles := make(map[string][]*output.FileNode)
	for _, f := range files {
		g.Nodes[f.Path] = f
		dir := getDir(f.Path)
		dirFiles[dir] = append(dirFiles[dir], f)
	}

	// 2. Build Edges
	for _, f := range files {
		for _, imp := range f.Imports {
			// If import is internal to the project
			// We need to know if 'imp' points to a dir in our project.
			// This requires knowing the module name to match internal imports.
			// Or check if 'imp' ends with any known directory suffix?

			// For now, let's look for "known directories".
			// If we have a directory "pkg/server" and import is "github.com/StartRP/repomap/pkg/server",
			// we match "pkg/server".

			matchedDir := ""
			for dir := range dirFiles {
				if strings.HasSuffix(imp, dir) { // Heuristic
					matchedDir = dir
					break
				}
			}

			if matchedDir != "" {
				// Add edges to ALL files in that directory?
				// Dependencies are package-level.
				// A file imports a package, thus depends on all files in that package (roughly).
				for _, target := range dirFiles[matchedDir] {
					if target.Path == f.Path {
						continue
					}
					g.Edges[f.Path] = append(g.Edges[f.Path], target.Path)
				}
			}
		}
	}
	return g
}

func getDir(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return strings.Join(parts[:len(parts)-1], "/")
	}
	return "."
}

// DetectCycles finds circular dependencies using DFS.
func (g *DependencyGraph) DetectCycles() []output.Issue {
	var issues []output.Issue
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	var dfs func(current string, path []string)
	dfs = func(current string, path []string) {
		visited[current] = true
		recursionStack[current] = true
		path = append(path, current)

		for _, neighbor := range g.Edges[current] {
			if !visited[neighbor] {
				dfs(neighbor, path)
			} else if recursionStack[neighbor] {
				// Cycle detected
				// Reconstruct cycle path
				cycleStart := -1
				for i, p := range path {
					if p == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart != -1 {
					cycle := path[cycleStart:]
					desc := fmt.Sprintf("Circular dependency detected: %s -> %s", strings.Join(cycle, " -> "), neighbor)

					issues = append(issues, output.Issue{
						Type:        "circular_dependency",
						Severity:    "high",
						Description: desc,
					})
				}
			}
		}
		recursionStack[current] = false
	}

	for node := range g.Nodes {
		if !visited[node] {
			dfs(node, []string{})
		}
	}
	return issues
}

// IntentValidator enforces architectural constraints based on file intents.
type IntentValidator struct {
	// Map of Intent -> disallowed imports (by keyword or path substring)
	DisallowedImports map[string][]string
}

// NewIntentValidator creates a default validator for Clean Architecture.
func NewIntentValidator() *IntentValidator {
	return &IntentValidator{
		DisallowedImports: map[string][]string{
			"domain":      {"infrastructure", "adapter", "delivery", "cmd"},
			"core":        {"infrastructure", "adapter", "delivery", "cmd"},
			"application": {"infrastructure", "delivery", "cmd"},
			"usecase":     {"infrastructure", "delivery", "cmd"},
		},
	}
}

// Validate checks files for intent violations.
func (v *IntentValidator) Validate(files []*output.FileNode) []output.Issue {
	var issues []output.Issue

	for _, f := range files {
		// Normalize intent to lowercase for matching
		intent := strings.ToLower(f.Intent)
		if intent == "" {
			continue
		}

		// Check if this intent has restrictions
		for key, disallowed := range v.DisallowedImports {
			if strings.Contains(intent, key) {
				for _, imp := range f.Imports {
					for _, bad := range disallowed {
						if strings.Contains(imp, bad) {
							desc := fmt.Sprintf("Architecture Violation: '%s' layer (%s) imports '%s' (%s)", intent, f.Path, bad, imp)
							issues = append(issues, output.Issue{
								Type:        "architectural_violation",
								Severity:    "high",
								Description: desc,
							})
						}
					}
				}
			}
		}
	}
	return issues
}
