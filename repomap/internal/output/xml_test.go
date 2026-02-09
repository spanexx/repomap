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

	xmlOutput, err := RenderXML(nodes)
	if err != nil {
		t.Fatalf("RenderXML failed: %v", err)
	}

	// Verify basic structure
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
			t.Logf("Got:\n%s", xmlOutput)
		}
	}
}
