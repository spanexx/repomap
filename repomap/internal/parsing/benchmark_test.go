package parsing

import (
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkExtractDefinitions measures parsing performance including I/O.
// MVP Goal: < 1s for 1K files. So < 1ms per file.
func BenchmarkExtractDefinitions(b *testing.B) {
	// Create a temp file once
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "bench_file.go")
	src := `package main

import (
	"fmt"
	"strings"
)

// BenchStruct represents a data structure for benchmarking.
type BenchStruct struct {
	ID   int
	Name string
	Tags []string
}

// Process handles the processing logic.
func (bs *BenchStruct) Process(data map[string]interface{}) (int, error) {
	return len(bs.Tags), nil
}

// HelperFunc assists with string manipulation.
func HelperFunc(s string) string {
	return strings.ToUpper(s)
}

// Interface defines the contract.
type Interface interface {
	DoWork()
}
`
	if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractGoDefinitions(filePath)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkExtractImports checks import extraction speed
func BenchmarkExtractImports(b *testing.B) {
	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "import_bench.go")
	src := `package main
import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
	"github.com/example/pkg/v2"
)
`
	if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractImports(filePath)
		if err != nil {
			b.Fatal(err)
		}
	}
}
