package discovery

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalk(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "discovery-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files
	files := []string{
		"main.go",
		"pkg/utils.go",
		"vendor/dep.go", // Should be included unless excluded explicitly (not yet)
		"bin/app.exe",    // Should be excluded (binary)
		".git/config",    // Should be excluded (hidden dir)
		"README.md",      // Should be excluded (not .go)
		"test_data/data.txt", // Should be excluded (not .go)
		"ignored.go", // Will be ignored by .gitignore
		"sub/ignored/file.go", // Will be ignored by directory match
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create .gitignore
	gitignoreContent := `
ignored.go
sub/ignored/
vendor/
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a hidden directory with a .go file inside
	hiddenDirPath := filepath.Join(tmpDir, ".hidden")
	if err := os.MkdirAll(hiddenDirPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDirPath, "secret.go"), []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a binary file
	binPath := filepath.Join(tmpDir, "binary.bin")
	if err := os.WriteFile(binPath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}


	// Run Walk
	foundFiles, err := Walk(tmpDir)
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Verify results
	expected := map[string]bool{
		filepath.Join(tmpDir, "main.go"):      true,
		filepath.Join(tmpDir, "pkg/utils.go"): true,
		// vendor/dep.go is ignored by .gitignore
		// ignored.go is ignored by .gitignore
		// sub/ignored/file.go is ignored by .gitignore
	}

	if len(foundFiles) != len(expected) {
		t.Errorf("Expected %d files, got %d", len(expected), len(foundFiles))
		for _, f := range foundFiles {
			t.Logf("Found: %s", f)
		}
	}

	for _, file := range foundFiles {
		if !expected[file] {
			t.Errorf("Unexpected file found: %s", file)
		}
	}
}
