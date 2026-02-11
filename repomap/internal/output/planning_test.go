package output

import (
	"encoding/json"
	"testing"
)

func TestFileNodePlanningSchema(t *testing.T) {
	node := FileNode{
		Path:       "pkg/auth/auth.go",
		Language:   "go",
		Importance: "high",
		Rank:       0.95,
		Status:     "planned",
		Intent:     "Handle authentication",
		Issues: []Issue{
			{Type: "circular_dependency", Description: "Imports pkg/user", Severity: "high"},
		},
		Comments: []Comment{
			{User: "alice", Text: "Needs OAuth support"},
		},
	}

	data, err := json.Marshal(node)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	jsonStr := string(data)

	// Basic string checks to ensure fields are present
	expectedFields := []string{
		`"status":"planned"`,
		`"intent":"Handle authentication"`,
		`"issues":[`,
		`"type":"circular_dependency"`,
		`"comments":[`,
		`"user":"alice"`,
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON missing field %s: %s", field, jsonStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(s)-len(substr) >= 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))
}
