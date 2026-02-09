package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitignore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitignore-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	gitignoreContent := `
# Comments are ignored
*.log
!important.log
/build
node_modules/
docs/*.md
`
	err = os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	gi, err := ParseGitignore(tmpDir)
	if err != nil {
		t.Fatalf("ParseGitignore failed: %v", err)
	}

	tests := []struct {
		path     string
		expected bool
	}{
		{"app.log", true},
		{"src/app.log", true},
		{"important.log", false},
		{"src/important.log", false}, // !important.log un-ignores it anywhere

		{"build", true},
		{"src/build", false}, // /build is anchored
		{"node_modules/foo.js", true}, // node_modules/ matches anywhere
		{"src/node_modules/foo.js", true},
		{"docs/readme.md", true},
		{"docs/api/index.md", false}, // docs/*.md only matches files in docs/
	}

	for _, tt := range tests {
		fullPath := filepath.Join(tmpDir, tt.path)
		if got := gi.Matches(fullPath); got != tt.expected {
			t.Errorf("Matches(%q) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}
