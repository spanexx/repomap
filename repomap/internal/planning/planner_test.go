package planning

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spanexx/agents-cli/repomap/internal/output"
)

func TestPlanner_LoadPlan(t *testing.T) {
	// Create a temporary plan file
	planContent := `{
		"version": "1.0",
		"intent": "Refactor auth",
		"changes": [
			{
				"path": "pkg/auth/new_auth.go",
				"status": "planned",
				"intent": "New authentication logic"
			}
		]
	}`

	tmpDir := t.TempDir()
	planPath := filepath.Join(tmpDir, "plan.json")
	if err := os.WriteFile(planPath, []byte(planContent), 0644); err != nil {
		t.Fatalf("failed to write temp plan: %v", err)
	}

	planner := NewPlanner()
	if err := planner.LoadPlan(planPath); err != nil {
		t.Fatalf("LoadPlan failed: %v", err)
	}

	if planner.Plan.Intent != "Refactor auth" {
		t.Errorf("Expected intent 'Refactor auth', got '%s'", planner.Plan.Intent)
	}

	if len(planner.Plan.Changes) != 1 {
		t.Fatalf("Expected 1 change, got %d", len(planner.Plan.Changes))
	}

	change := planner.Plan.Changes[0]
	if change.Path != "pkg/auth/new_auth.go" {
		t.Errorf("Expected path 'pkg/auth/new_auth.go', got '%s'", change.Path)
	}
}

func TestPlanner_ApplyPlan(t *testing.T) {
	// Setup initial repo map
	repoMap := &output.RepoMap{
		Files: []*output.FileNode{
			{Path: "main.go", Language: "go"},
			{Path: "pkg/old.go", Language: "go"},
		},
	}

	// Setup planner with a plan
	planner := NewPlanner()
	planner.Plan = &Plan{
		Changes: []FileChange{
			{
				Path:   "pkg/old.go",
				Status: "modified",
				Intent: "Deprecate this",
			},
			{
				Path:   "pkg/new.go",
				Status: "planned",
				Intent: "Replacement",
			},
		},
	}

	// Apply plan
	if err := planner.ApplyPlan(repoMap); err != nil {
		t.Fatalf("ApplyPlan failed: %v", err)
	}

	// Verify results
	if len(repoMap.Files) != 3 {
		t.Fatalf("Expected 3 files, got %d", len(repoMap.Files))
	}

	// Check mapping
	fileMap := make(map[string]*output.FileNode)
	for _, f := range repoMap.Files {
		fileMap[f.Path] = f
	}

	// Verify modified file
	oldFile, ok := fileMap["pkg/old.go"]
	if !ok {
		t.Fatal("pkg/old.go missing")
	}
	if oldFile.Status != "modified" {
		t.Errorf("pkg/old.go status: expected 'modified', got '%s'", oldFile.Status)
	}
	if oldFile.Intent != "Deprecate this" {
		t.Errorf("pkg/old.go intent: expected 'Deprecate this', got '%s'", oldFile.Intent)
	}

	// Verify new file
	newFile, ok := fileMap["pkg/new.go"]
	if !ok {
		t.Fatal("pkg/new.go missing")
	}
	if newFile.Status != "planned" {
		t.Errorf("pkg/new.go status: expected 'planned', got '%s'", newFile.Status)
	}
}
