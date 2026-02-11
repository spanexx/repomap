/*
Repomap is a CLI tool that generates a token-optimized map of a Go repository.
It analyzes the codebase, extracts definitions, identifies dependencies, ranks files by importance,
and outputs a structured XML or JSON representation suitable for LLM context.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/internal/analysis"
	"github.com/spanexx/agents-cli/repomap/internal/discovery"
	"github.com/spanexx/agents-cli/repomap/internal/graph"
	"github.com/spanexx/agents-cli/repomap/internal/output" // Internal output structs
	"github.com/spanexx/agents-cli/repomap/internal/parsing"
	"github.com/spanexx/agents-cli/repomap/internal/planning"
	"github.com/spanexx/agents-cli/repomap/internal/ranking"
	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/cli"
	"github.com/spanexx/agents-cli/repomap/pkg/config"
	pkgOutput "github.com/spanexx/agents-cli/repomap/pkg/output" // Framework output writer
	"github.com/spanexx/agents-cli/repomap/pkg/providers/gemini_cli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qodercli"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/qwen_cli"
	"github.com/spanexx/agents-cli/repomap/pkg/server"
	"github.com/spanexx/agents-cli/repomap/pkg/util"
)

const version = "0.1.0"

func main() {
	app := cli.NewApp("repomap", version)
	app.SetDescription("Generate a token-optimized map of your Go repository.")
	app.AddExample("repomap --root . --output xml")
	app.AddExample("repomap --output json --max-tokens 4000")

	// CLI Flags
	app.AddFlag("root", "Repository root directory", ".")
	app.AddFlag("output", "Output format (xml|json|text)", "xml")
	app.AddFlag("max-tokens", "Maximum token budget (0 for unlimited)", 0)
	app.AddFlag("include-ext", "Comma-separated extensions to include (default: .go)", "")
	app.AddFlag("exclude-ext", "Comma-separated extensions to exclude", "")
	app.AddFlag("ignore-tests", "Ignore test files (*_test.go)", false)
	app.AddFlag("verbose", "Enable verbose logging", false)
	app.AddFlag("version", "Show version information", false)

	// Server Flags
	app.AddFlag("serve", "Start Visualizer server and Agent API", false)
	app.AddFlag("port", "Server port", "8080")
	app.AddFlag("plan", "Path to plan file for agent interaction", "PLAN/plan.json")
	app.AddFlag("analyze", "Run static analysis (duplication, intent, cycles)", false)

	// Parse Flags
	flags, err := app.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if flags.GetBool("version") {
		fmt.Printf("repomap version %s\n", version)
		os.Exit(0)
	}

	// Initialize Logger
	logger := util.NewLogger(os.Stdout, os.Stderr, flags.GetBool("verbose"))

	// Load Configuration (merging defaults -> config -> flags)
	// For now, we simple take flags as source of truth, but loading config is part of framework.
	cfgPaths := config.DefaultPaths("repomap")
	cfg, err := config.LoadConfig(cfgPaths)
	if err != nil && flags.GetBool("verbose") {
		logger.Warn("Failed to load config: %v", err)
	}

	// TODO: Implement proper hierarchy merging. For now, we prefer flags.
	// In a real implementation of Task 1.2.11, we would merge these.
	// Here we just use flags directly as the priority.

	visited := flags.GetVisitedValues()

	// 1. Validation
	rootDir := flags.GetString("root")
	if cfg.GetString("root") != "" {
		if _, ok := visited["root"]; !ok {
			// Only use config when the flag wasn't explicitly set.
			// This prevents config from overriding a user's explicit `--root .`.
			if rootDir == "." {
				rootDir = cfg.GetString("root")
			}
		}
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		logger.Error("Failed to resolve root directory: %v", err)
		os.Exit(1)
	}

	outputFmt := flags.GetString("output")
	if cfg.GetString("output") != "" {
		if _, ok := visited["output"]; !ok {
			// Use config if flag wasn't explicitly set
			outputFmt = cfg.GetString("output")
		}
	}

	// Validate output format using framework factory check (or just try to create one)
	writer, err := pkgOutput.NewWriter(outputFmt)
	if err != nil {
		logger.Error("Invalid output format: %v", err)
		os.Exit(1)
	}

	start := time.Now()
	// Use explicit logging instead of 'log' function
	if flags.GetBool("verbose") {
		logger.Info("Starting Repomap on %s", absRoot)
	}

	// 2. Discovery
	logger.Debug("Phase A: Discovering files...")

	files, err := discovery.Walk(absRoot)
	if err != nil {
		logger.Error("Discovery failed: %v", err)
		os.Exit(1)
	}

	// Apply CLI filters
	filteredFiles := filterFiles(files, flags)
	logger.Debug("Found %d files", len(filteredFiles))

	// 3. Parsing (Definitions & Imports)
	logger.Debug("Phase B: Parsing %d files...", len(filteredFiles))

	graphBuilder := graph.NewBuilder()
	fileNodes := make([]*output.FileNode, 0, len(filteredFiles))

	for _, path := range filteredFiles {
		// Extract Definitions via Registry
		extractor := parsing.DefaultRegistry.Get(path)
		var defs []string
		if extractor != nil {
			defs, err = extractor.ExtractDefinitions(path)
			if err != nil {
				logger.Warn("Failed to parse definitions for %s: %v", path, err)
			}
		}

		// Extract Imports via Registry
		var imports []string
		if extractor != nil {
			imports, err = extractor.ExtractImports(path)
			if err != nil {
				logger.Warn("Failed to parse imports for %s: %v", path, err)
			}
		}

		// Add to Graph Builder (relative to root)
		relPath, _ := filepath.Rel(absRoot, path)
		relPath = filepath.ToSlash(relPath)

		graphBuilder.AddFile(relPath, imports)

		ext := strings.ToLower(filepath.Ext(path))
		lang := strings.TrimPrefix(ext, ".")
		if lang == "" {
			lang = "unknown"
		}

		fileNodes = append(fileNodes, &output.FileNode{
			Path:        relPath,
			Language:    lang,
			Definitions: defs,
			Imports:     imports,
			TokenCount:  0,
		})
	}

	// 4. Graph Construction
	logger.Debug("Phase C: Building import graph...")
	moduleName := findModuleName(absRoot)
	importGraph := graphBuilder.Build(moduleName)

	// 5. Ranking
	logger.Debug("Phase D: Ranking files...")
	ranks := ranking.Rank(importGraph)
	importance := ranking.AssignImportance(ranks)

	// Enrich file nodes
	for _, node := range fileNodes {
		if r, ok := ranks[node.Path]; ok {
			node.Rank = r
		}
		if imp, ok := importance[node.Path]; ok {
			node.Importance = imp
		} else {
			node.Importance = "low"
		}
	}

	// Sort by Rank
	sort.Slice(fileNodes, func(i, j int) bool {
		return fileNodes[i].Rank > fileNodes[j].Rank
	})

	// 5.5 Intent Assignment (Heuristic + LLM)
	// We do this before output so it's included in the result.
	var provider adapter.Provider
	if flags.GetBool("analyze") {
		// Initialize provider for analysis if requested
		var err error
		provider, err = initProvider()
		if err != nil {
			logger.Warn("Failed to initialize LLM provider for intent analysis: %v (falling back to heuristics)", err)
		}
	}
	analysis.AssignIntent(fileNodes, provider)

	// 6. Output
	logger.Debug("Phase E: Rendering output...")

	// Wrap in RepoMap for proper root element in XML/JSON
	result := &output.RepoMap{
		Files: fileNodes,
	}

	// 5.5 Planning (Merge Plan)
	planPath := flags.GetString("plan")
	if planPath != "" {
		// Only try to load if file exists
		if _, err := os.Stat(planPath); err == nil {
			logger.Debug("Loading plan from %s", planPath)
			planner := planning.NewPlanner()
			if err := planner.LoadPlan(planPath); err != nil {
				logger.Warn("Failed to load plan: %v", err)
			} else {
				if err := planner.ApplyPlan(result); err != nil {
					logger.Warn("Failed to apply plan: %v", err)
				} else {
					logger.Info("Applied plan from %s", planPath)
				}
			}
		}
	}

	// 5.6 Static Analysis
	if flags.GetBool("analyze") {
		logger.Info("Running static analysis...")

		// Map generated result files to a map for easy lookup
		nodeMap := make(map[string]*output.FileNode)
		contentMap := make(map[string][]byte)

		for _, node := range result.Files {
			nodeMap[node.Path] = node
			// Read content for duplication detection
			// Only for existing files (Status != "planned")
			fullPath := filepath.Join(absRoot, node.Path)
			if _, err := os.Stat(fullPath); err == nil {
				data, err := os.ReadFile(fullPath)
				if err == nil {
					contentMap[node.Path] = data
				}
			}
		}

		// A. Duplication Detection
		dupeDetector := analysis.NewDuplicationDetector()
		dupes, err := dupeDetector.Analyze(contentMap)
		if err != nil {
			logger.Warn("Duplication analysis failed: %v", err)
		} else {
			logger.Info("Found %d duplication issues", len(dupes))
			// Distribute issues to files?
			// The issues describe "files involved".
			// We can attach the issue to all involved files.
			// Currently Issue struct is simple.
			// Let's attach to the first file mentioned in description?
			// Or just attach to a random one?
			// Better: Duplicate detector should probably return map[filename][]Issue?
			// Current implementation returns []Issue.
			// For MVP, simplistic attachment:
			for _, issue := range dupes {
				// Parse paths from description? (Brittle but working based on implementation)
				// "Found duplicate code block ... across X files: path/a.go, path/b.go"
				blockParts := strings.Split(issue.Description, ": ")
				if len(blockParts) > 1 {
					paths := strings.Split(blockParts[1], ", ")
					for _, p := range paths {
						if node, ok := nodeMap[p]; ok {
							node.Issues = append(node.Issues, issue)
						}
					}
				}
			}
		}

		// B. Intent Validation
		intentValidator := analysis.NewIntentValidator()
		intentIssues := intentValidator.Validate(result.Files)
		logger.Info("Found %d intent violations", len(intentIssues))
		for _, issue := range intentIssues {
			// Parse path from description?
			// "Architecture Violation: 'intent' layer (path/to/file) imports..."
			// We need a better way to link issues to files.
			// Validator.Validate loops over files.
			// Maybe we should update Validator to attach issues directly?
			// Or just parse.
			start := strings.Index(issue.Description, "(")
			end := strings.Index(issue.Description, ")")
			if start != -1 && end != -1 && end > start {
				path := issue.Description[start+1 : end]
				if node, ok := nodeMap[path]; ok {
					node.Issues = append(node.Issues, issue)
				}
			}
		}

		// C. Circular Dependencies
		depGraph := analysis.BuildGraph(result.Files, "")
		cycleIssues := depGraph.DetectCycles()
		logger.Info("Found %d circular dependencies", len(cycleIssues))
		for _, issue := range cycleIssues {
			// "Circular dependency detected: a -> b -> a"
			// Attach to the last node (neighbor)
			parts := strings.Split(issue.Description, " -> ")
			if len(parts) > 0 {
				lastPath := parts[len(parts)-1]
				if node, ok := nodeMap[lastPath]; ok {
					node.Issues = append(node.Issues, issue)
				}
			}
		}
	}

	// If serving, start server instead of writing to file/stdout (or do both?)
	// Typically we might want to see the output AND serve it.
	// But let's assume if --serve is on, we block on server.
	if flags.GetBool("serve") {
		port := flags.GetString("port")
		planPath := flags.GetString("plan")

		logger.Info("Starting server with generated map...")
		srv := server.New(port, planPath, result)

		if err := srv.Start(); err != nil {
			logger.Error("Server failed: %v", err)
			os.Exit(1)
		}
		// Server blocks, so we exit when it stops
		os.Exit(0)
	}

	if err := writer.Write(result); err != nil {
		logger.Error("Rendering failed: %v", err)
		os.Exit(1)
	}

	// Ensure we print a newline? Writer might do it.

	if flags.GetBool("verbose") {
		logger.Info("Completed in %v", time.Since(start))
	}
}

func filterFiles(files []string, flags *cli.Flags) []string {
	var filtered []string

	includeStr := flags.GetString("include-ext")
	excludeStr := flags.GetString("exclude-ext")
	ignoreTests := flags.GetBool("ignore-tests")

	incExts := splitExts(includeStr)
	excExts := splitExts(excludeStr)

	if len(incExts) == 0 {
		// Default to all supported extensions if none specified
		incExts = parsing.DefaultRegistry.SupportedExtensions()
	}

	for _, f := range files {
		ext := filepath.Ext(f)

		if contains(excExts, ext) {
			continue
		}

		if ignoreTests && strings.HasSuffix(f, "_test.go") {
			continue
		}

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

func initProvider() (adapter.Provider, error) {
	ctx := context.Background()
	providerName := os.Getenv("REPOMAP_PROVIDER")
	if providerName == "" {
		return nil, nil // No provider configured, use static only
	}

	providerNames := strings.Split(providerName, ",")
	var providers []adapter.Provider

	for _, name := range providerNames {
		name = strings.TrimSpace(name)
		var p adapter.Provider

		switch name {
		case "qwen-cli":
			p = qwen_cli.New()
		case "qodercli":
			p = qodercli.New()
		case "gemini-cli":
			gp, err := gemini_cli.New(ctx)
			if err != nil {
				fmt.Printf("Warning: failed to create gemini provider: %v\n", err)
				continue
			}
			if err := gp.Init(ctx); err != nil {
				fmt.Printf("Warning: Gemini provider init failed: %v\n", err)
			}
			p = gp
		default:
			fmt.Printf("Warning: unknown provider: %s\n", name)
			continue
		}
		if p != nil {
			providers = append(providers, p)
		}
	}

	if len(providers) == 0 {
		return nil, fmt.Errorf("no valid providers found in: %s", providerName)
	}

	if len(providers) == 1 {
		return providers[0], nil
	}

	return adapter.NewFallbackProvider(providers, true), nil
}
