package planning

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spanexx/agents-cli/repomap/internal/output"
)

// Planner handles the ingestion and application of architectural plans.
type Planner struct {
	// Plan is the loaded plan, if any.
	Plan *Plan
}

// Plan represents the structure of a plan.json file.
type Plan struct {
	// Version of the plan format.
	Version string `json:"version"`
	// Intent describes the high-level goal of this plan.
	Intent string `json:"intent"`
	// Changes is a list of planned changes to files.
	Changes []FileChange `json:"changes"`
}

// FileChange represents a planned change to a specific file.
type FileChange struct {
	Path        string           `json:"path"`
	Status      string           `json:"status"` // "planned", "modified", "deleted"
	Intent      string           `json:"intent,omitempty"`
	Description string           `json:"description,omitempty"`
	Issues      []output.Issue   `json:"issues,omitempty"`
	Comments    []output.Comment `json:"comments,omitempty"`
}

// NewPlanner creates a new Planner instance.
func NewPlanner() *Planner {
	return &Planner{}
}

// LoadPlan loads a plan from a JSON file.
func (p *Planner) LoadPlan(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan Plan
	if err := json.Unmarshal(data, &plan); err != nil {
		return fmt.Errorf("failed to parse plan file: %w", err)
	}

	p.Plan = &plan
	return nil
}

// ApplyPlan merges the loaded plan into the existing repository map.
// It updates existing nodes or appends new ones based on the plan.
func (p *Planner) ApplyPlan(repoMap *output.RepoMap) error {
	if p.Plan == nil {
		return fmt.Errorf("no plan loaded")
	}

	// Create a lookup map for existing files
	fileMap := make(map[string]*output.FileNode)
	for _, file := range repoMap.Files {
		fileMap[file.Path] = file
	}

	for _, change := range p.Plan.Changes {
		node, exists := fileMap[change.Path]

		if !exists {
			// Create new node for planned file
			node = &output.FileNode{
				Path:       change.Path,
				Status:     change.Status,
				Intent:     change.Intent,
				Issues:     change.Issues,
				Comments:   change.Comments,
				Importance: "low", // Default for new files until ranked
			}
			repoMap.Files = append(repoMap.Files, node)
		} else {
			// Update existing node
			node.Status = change.Status
			if change.Intent != "" {
				node.Intent = change.Intent
			}
			if len(change.Issues) > 0 {
				node.Issues = append(node.Issues, change.Issues...)
			}
			if len(change.Comments) > 0 {
				node.Comments = append(node.Comments, change.Comments...)
			}
		}
	}

	return nil
}
