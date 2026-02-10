package parsing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractDefinitionsComplexTypes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsing-complex-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := `
package main

type MyMap map[string]int
type MyChan chan int
type MyRecvChan <-chan int
type MySendChan chan<- int
type MyArray [5]int
type MySlice []string
type MyInterface interface{}
type MyEllipsisFunc func(...string)

func ComplexFunc(
	m map[string]interface{}, 
	ch chan bool, 
	arr [10]byte, 
	slice []string,
	ptr *int,
	variadic ...string,
) (func(int), interface{}) {
	return nil, nil
}
`
	filePath := filepath.Join(tmpDir, "complex.go")
	if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	defs, err := ExtractDefinitions(filePath)
	if err != nil {
		t.Fatalf("ExtractDefinitions failed: %v", err)
	}

	expected := []string{
		"type MyMap map[string]int",
		"type MyChan chan int",
		"type MyRecvChan chan int", // Direction ignored in MVP
		"type MySendChan chan int", // Direction ignored in MVP
		"type MyArray [5]int",
		"type MySlice []string",
		"type MyInterface interface",
		"type MyEllipsisFunc func",
		"func ComplexFunc(map[string]interface{}, chan bool, [10]byte, []string, *int, ...string) (func, interface{})",
	}

	if len(defs) != len(expected) {
		t.Errorf("Expected %d definitions, got %d", len(expected), len(defs))
		for i, d := range defs {
			t.Logf("Got[%d]: %s", i, d)
		}
		return
	}

	for i, d := range defs {
		if d != expected[i] {
			t.Errorf("Definition mismatch at %d: got %q, want %q", i, d, expected[i])
		}
	}
}
