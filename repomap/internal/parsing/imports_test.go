package parsing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractImports(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsing-imports-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := `
package main

import (
	"fmt"
	"os"
)

import "net/http"
import log "log" // aliased
import . "path/filepath" // dot import
import _ "image/png" // side-effect import

func main() {
	fmt.Println("Hello")
}
`
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	imports, err := ExtractImports(filePath)
	if err != nil {
		t.Fatalf("ExtractImports failed: %v", err)
	}

	expected := []string{
		"fmt",
		"os",
		"net/http",
		"log",
		"path/filepath",
		"image/png",
	}

	if len(imports) != len(expected) {
		t.Errorf("Expected %d imports, got %d", len(expected), len(imports))
		for i, imp := range imports {
			t.Logf("Got[%d]: %s", i, imp)
		}
	}

	expectedMap := make(map[string]bool)
	for _, e := range expected {
		expectedMap[e] = true
	}

	for _, imp := range imports {
		if !expectedMap[imp] {
			t.Errorf("Unexpected import found: %s", imp)
		}
	}
}
