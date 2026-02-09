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

	// Case 1: No limit
	jsonOutput, err := RenderJSON(nodes, 0)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

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
		}
	}

	// Case 2: Budget limit
	// Wrapper ~3 tokens. Node 1 cost 100. Total 103.
	// We set limit 110. Node 2 cost 50. 103 + 50 = 153 > 110.

	jsonOutputLimited, err := RenderJSON(nodes, 110)
	if err != nil {
		t.Fatalf("RenderJSON limited failed: %v", err)
	}

	if !strings.Contains(jsonOutputLimited, "main.go") {
		t.Error("Limited JSON should contain main.go")
	}
	if strings.Contains(jsonOutputLimited, "pkg/utils.go") {
		t.Error("Limited JSON should NOT contain pkg/utils.go")
	}
	if !strings.Contains(jsonOutputLimited, "Output truncated") {
		t.Error("Limited JSON should contain truncation notice")
	}
}
