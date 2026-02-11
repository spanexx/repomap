package analysis

import (
	"testing"

	"github.com/spanexx/agents-cli/repomap/internal/output"
)

func TestIntentValidator(t *testing.T) {
	validator := NewIntentValidator()

	files := []*output.FileNode{
		{
			Path:    "internal/domain/user.go",
			Intent:  "Domain Logic",
			Imports: []string{"fmt", "github.com/project/internal/infrastructure/db"},
		},
		{
			Path:    "internal/domain/order.go",
			Intent:  "Domain Logic",
			Imports: []string{"fmt"},
		},
		{
			Path:    "internal/usecase/handler.go",
			Intent:  "Application Use Case",
			Imports: []string{"github.com/project/internal/domain"},
		},
	}

	issues := validator.Validate(files)

	if len(issues) != 1 {
		t.Errorf("Expected 1 violation, got %d", len(issues))
	}

	if len(issues) > 0 {
		if issues[0].Type != "architectural_violation" {
			t.Errorf("Expected architectural_violation, got %s", issues[0].Type)
		}
		t.Logf("Found expected violation: %s", issues[0].Description)
	}
}

func TestCycleDetection(t *testing.T) {
	// A -> B -> A
	files := []*output.FileNode{
		{Path: "pkg/a/a.go", Imports: []string{"pkg/b"}},
		{Path: "pkg/b/b.go", Imports: []string{"pkg/a"}},
	}

	// We need to Mock BuildGraph logic?
	// BuildGraph relies on mapping import paths to files.
	// Our BuildGraph matches "pkg/b" to "pkg/b/b.go" if suffix matches.

	g := BuildGraph(files, "")
	issues := g.DetectCycles()

	// Recursion might find multiple paths for the same cycle A->B->A and B->A->B
	// But our DFS tracks visited per traversal.
	// We should find at least one.

	if len(issues) == 0 {
		t.Errorf("Expected circular dependency, got 0 issues")
	}

	for _, i := range issues {
		t.Logf("Found cycle: %s", i.Description)
	}
}
