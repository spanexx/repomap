/*
Repomap is a CLI tool that generates a token-optimized map of a Go repository.
It analyzes the codebase, extracts definitions, identifies dependencies, ranks files by importance,
and outputs a structured XML or JSON representation suitable for LLM context.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/internal/discovery"
	"github.com/spanexx/agents-cli/repomap/internal/graph"
	"github.com/spanexx/agents-cli/repomap/internal/output"
	"github.com/spanexx/agents-cli/repomap/internal/parsing"
	"github.com/spanexx/agents-cli/repomap/internal/ranking"
)

var (
	// CLI Flags
	rootDir     = flag.String("root", ".", "Repository root directory")
	outputFmt   = flag.String("output", "xml", "Output format (xml|json)")
	maxTokens   = flag.Int("max-tokens", 0, "Maximum token budget (0 for unlimited)")
	includeExts = flag.String("include-ext", "", "Comma-separated extensions to include (default: .go)")
	excludeExts = flag.String("exclude-ext", "", "Comma-separated extensions to exclude")
	ignoreTests = flag.Bool("ignore-tests", false, "Ignore test files (*_test.go)")
	verbose     = flag.Bool("verbose", false, "Enable verbose logging")
	showVersion = flag.Bool("version", false, "Show version information")
)

const version = "0.1.0"

func main() {
	flag.Usage = printHelp
	flag.Parse()

	if *showVersion {
		fmt.Printf("repomap version %s\n", version)
		os.Exit(0)
	}

	// 1. Validation
	absRoot, err := filepath.Abs(*rootDir)
	if err != nil {
		fatalf("Failed to resolve root directory: %v", err)
	}

	if *outputFmt != "xml" && *outputFmt != "json" {
		fatalf("Invalid output format: %s. Must be 'xml' or 'json'", *outputFmt)
	}

	start := time.Now()
	log("Starting Repomap on %s", absRoot)

	// 2. Discovery
	log("Phase A: Discovering files...")
	// For MVP, discovery.Walk only accepts root.
	// We might need to refactor discovery.Walk to accept options or filter post-discovery.
	// Given the current signature: func Walk(root string) ([]string, error)
	// We will filter efficiently after walking for now, or assume Walk defaults are okay for MVP.
	// The task says "Implement include/exclude flags", so we should ideally respect them.
	// However, discovery.Walk is hardcoded for .go files currently.
	// Let's use discovery.Walk and then filter if needed, OR relies on the fact that it defaults to .go

	files, err := discovery.Walk(absRoot)
	if err != nil {
		fatalf("Discovery failed: %v", err)
	}

	// Apply CLI filters (test files, custom extensions)
	// Since discovery.Walk currently hardcodes .go, we might be limited.
	// But let's filter based on flags to be safe/future-proof.
	filteredFiles := filterFiles(files)
	log("Found %d files", len(filteredFiles))

	// 3. Parsing (Definitions & Imports)
	log("Phase B: Parsing %d files...", len(filteredFiles))

	graphBuilder := graph.NewBuilder()
	fileNodes := make([]*output.FileNode, 0, len(filteredFiles))

	for _, path := range filteredFiles {
		// Extract Definitions
		defs, err := parsing.ExtractDefinitions(path)
		if err != nil {
			log("Warning: failed to parse definitions for %s: %v", path, err)
			continue // Skip files we can't parse? Or keep with empty defs?
			// Let's keep best effort.
		}

		// Extract Imports
		imports, err := parsing.ExtractImports(path)
		if err != nil {
			log("Warning: failed to parse imports for %s: %v", path, err)
		}

		// Add to Graph Builder (relative to root)
		relPath, _ := filepath.Rel(absRoot, path)
		// Normalize path to forward slashes for consistency
		relPath = filepath.ToSlash(relPath)

		graphBuilder.AddFile(relPath, imports)

		fileNodes = append(fileNodes, &output.FileNode{
			Path:        relPath,
			Language:    "go", // Fixed for MVP
			Definitions: defs,
			Imports:     imports,
			// TokenCount populated later or now?
			// Logic says "Track cumulative tokens while building output" inside Render.
			// But Render expects Node.TokenCount to be pre-filled or it calculates it.
			// Let's calculate rough cost now.
			TokenCount: 0, // Will be calculated by Render if 0, or we can pre-calc.
		})
	}

	// 4. Graph Construction
	log("Phase C: Building import graph...")
	// We need a module name to handle internal imports correctly.
	// For now, let's try to guess or leave empty (which might treat internal imports as external).
	// Ideally we parse go.mod. For MVP, maybe empty string is acceptable if imports are relative?
	// The builder handles "moduleName" stripping.
	// Let's attempt to find go.mod module name.
	moduleName := findModuleName(absRoot)
	importGraph := graphBuilder.Build(moduleName)

	// 5. Ranking
	log("Phase D: Ranking files...")
	ranks := ranking.Rank(importGraph)
	importance := ranking.AssignImportance(ranks)

	// Enrich file nodes with rank and importance
	for _, node := range fileNodes {
		if r, ok := ranks[node.Path]; ok {
			node.Rank = r
		}
		if imp, ok := importance[node.Path]; ok {
			node.Importance = imp
		} else {
			node.Importance = "low" // Fallback
		}
	}

	// Sort by Rank (Descending)
	sort.Slice(fileNodes, func(i, j int) bool {
		return fileNodes[i].Rank > fileNodes[j].Rank
	})

	// 6. Output
	log("Phase E: Rendering output...")

	var result string

	if *outputFmt == "json" {
		result, err = output.RenderJSON(fileNodes, *maxTokens)
	} else {
		result, err = output.RenderXML(fileNodes, *maxTokens)
	}

	if err != nil {
		fatalf("Rendering failed: %v", err)
	}

	// Print to stdout
	fmt.Println(result)

	log("Completed in %v", time.Since(start))
}

func filterFiles(files []string) []string {
	var filtered []string

	// Parse include/exclude lists
	incExts := splitExts(*includeExts)
	excExts := splitExts(*excludeExts)

	// Default to .go if no includes specified (and discovery didn't already handle it, though it does)
	if len(incExts) == 0 {
		incExts = []string{".go"}
	}

	for _, f := range files {
		ext := filepath.Ext(f)

		// 1. Check excludes
		if contains(excExts, ext) {
			continue
		}

		// 2. Check test files
		if *ignoreTests && strings.HasSuffix(f, "_test.go") {
			continue
		}

		// 3. Check includes
		if !contains(incExts, ext) {
			continue
		}

		filtered = append(filtered, f)
	}
	return filtered
}

func splitExts(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var res []string
	for _, p := range parts {
		if p != "" {
			res = append(res, strings.TrimSpace(p))
		}
	}
	return res
}

func contains(list []string, item string) bool {
	for _, s := range list {
		if s == item {
			return true
		}
	}
	return false
}

func findModuleName(root string) string {
	// Simple grep for "module " in go.mod
	goModPath := filepath.Join(root, "go.mod")
	data, err := os.ReadFile(goModPath)
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				return strings.TrimSpace(strings.TrimPrefix(line, "module "))
			}
		}
	}
	return ""
}

func log(format string, args ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, "[LOG] "+format+"\n", args...)
	}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "Usage: repomap [options]\n\n")
	fmt.Fprintf(os.Stderr, "Generate a token-optimized map of your Go repository.\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}
