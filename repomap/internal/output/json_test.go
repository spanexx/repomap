package output

import (
	"strings"
	"testing"
)

func TestRenderJSON(t *testing.T) {
	nodes := []*FileNode{
		{
			Path:        "main.go",
			Language:    "go",
			Importance:  "high",
			Rank:        0.95,
			Definitions: []string{"func main()"},
			TokenCount:  100,
		},
		{
			Path:        "pkg/utils.go",
			Language:    "go",
			Importance:  "medium",
			Rank:        0.5,
			Definitions: []string{"func Helper()"},
			TokenCount:  50,
		},
	}

	jsonOutput, err := RenderJSON(nodes)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	// Verify basic structure
	expectedSubstrings := []string{
		`"files": [`,
		`"path": "main.go"`,
		`"language": "go"`,
		`"importance": "high"`,
		`"rank": 0.95`,
		`"definitions": [`,
		`"func main()"`,
		`"path": "pkg/utils.go"`,
		`"rank": 0.5`,
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(jsonOutput, s) {
			t.Errorf("JSON output missing expected string: %s", s)
			t.Logf("Got:\n%s", jsonOutput)
		}
	}
}
