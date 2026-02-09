package output

import (
	"strings"
	"testing"
)

func TestRenderXML(t *testing.T) {
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
	xmlOutput, err := RenderXML(nodes, 0)
	if err != nil {
		t.Fatalf("RenderXML failed: %v", err)
	}

	expectedSubstrings := []string{
		`<?xml version="1.0" encoding="UTF-8"?>`,
		`<repomap>`,
		`<file path="main.go" language="go" importance="high" rank="0.95" token_count="100">`,
		`<definition>func main()</definition>`,
		`<file path="pkg/utils.go" language="go" importance="medium" rank="0.5" token_count="50">`,
		`<definition>func Helper()</definition>`,
	}

	for _, s := range expectedSubstrings {
		if !strings.Contains(xmlOutput, s) {
			t.Errorf("XML output missing expected string: %s", s)
		}
	}

	// Case 2: Budget limit (allow first node, truncates second)
	// Header ~12 tokens. Node 1 cost 100. Total 112.
	// We set maxTokens to 120. Node 2 cost 50. 112 + 50 = 162 > 120.
	// So Node 2 should be truncated.

	xmlOutputLimited, err := RenderXML(nodes, 120)
	if err != nil {
		t.Fatalf("RenderXML limited failed: %v", err)
	}

	if !strings.Contains(xmlOutputLimited, "main.go") {
		t.Error("Limited XML should contain main.go")
	}
	if strings.Contains(xmlOutputLimited, "pkg/utils.go") {
		t.Error("Limited XML should NOT contain pkg/utils.go")
	}
	if !strings.Contains(xmlOutputLimited, "Output truncated") {
		t.Error("Limited XML should contain truncation notice")
	}
}
