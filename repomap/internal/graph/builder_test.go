package graph

import (
	"testing"
)

func TestGraphBuilder(t *testing.T) {
	// Create a builder
	builder := NewBuilder()

	// Add files and their imports
	// Simulating a project with module "github.com/example/repo"
	// Files:
	// - main.go -> github.com/example/repo/pkg/utils
	// - pkg/utils/util.go -> (no internal imports)
	// - pkg/utils/helper.go -> (no internal imports)
	// - pkg/config/config.go -> github.com/example/repo/pkg/utils

	moduleName := "github.com/example/repo"

	files := map[string][]string{
		"main.go": {
			"fmt",
			"github.com/example/repo/pkg/utils",
			"github.com/example/repo/pkg/config",
		},
		"pkg/utils/util.go": {
			"fmt",
		},
		"pkg/utils/helper.go": {
			"fmt",
		},
		"pkg/config/config.go": {
			"github.com/example/repo/pkg/utils",
		},
	}

	for path, imports := range files {
		builder.AddFile(path, imports)
	}

	// Build graph
	g := builder.Build(moduleName)

	// Verify nodes
	if len(g.Nodes) != len(files) {
		t.Errorf("Expected %d nodes, got %d", len(files), len(g.Nodes))
	}

	// Verify edges and in-degrees
	// main.go -> pkg/utils/util.go
	// main.go -> pkg/utils/helper.go
	// main.go -> pkg/config/config.go

	// pkg/config/config.go -> pkg/utils/util.go
	// pkg/config/config.go -> pkg/utils/helper.go

	// In-degrees:
	// main.go: 0
	// pkg/utils/util.go: 2 (from main and config)
	// pkg/utils/helper.go: 2 (from main and config)
	// pkg/config/config.go: 1 (from main)

	expectedInDegree := map[string]int{
		"main.go":              0,
		"pkg/utils/util.go":    2,
		"pkg/utils/helper.go":  2,
		"pkg/config/config.go": 1,
	}

	for path, degree := range expectedInDegree {
		if node, ok := g.Nodes[path]; !ok {
			t.Errorf("Node not found: %s", path)
		} else if node.InDegree != degree {
			t.Errorf("Node %s: expected in-degree %d, got %d", path, degree, node.InDegree)
		}
	}

	// Verify edge correctness (simplified check)
	edges := g.Edges["main.go"]

	// We expect 3 edges: config.go, util.go, helper.go
	expectedEdges := map[string]bool{
		"pkg/utils/util.go":    true,
		"pkg/utils/helper.go":  true,
		"pkg/config/config.go": true,
	}

	if len(edges) != 3 {
		t.Errorf("main.go: expected 3 edges, got %d", len(edges))
	}

	for _, edge := range edges {
		if !expectedEdges[edge] {
			t.Errorf("main.go: unexpected edge to %s", edge)
		}
	}
}
