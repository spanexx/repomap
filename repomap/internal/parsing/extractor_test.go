package parsing

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractDefinitions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parsing-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	src := `
package main

import "fmt"

type Config struct {
	Port int
}

type Service interface {
	Start() error
}

func main() {
	fmt.Println("Hello")
}

func NewServer(port int) *Server {
	return &Server{Port: port}
}

func (s *Server) Start() error {
	return nil
}

func (s Server) Stop() (error, bool) {
	return nil, true
}

type Handler func(w http.ResponseWriter, r *http.Request)

type MyInt int
`
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(src), 0644); err != nil {
		t.Fatal(err)
	}

	defs, err := ExtractDefinitions(filePath)
	if err != nil {
		t.Fatalf("ExtractDefinitions failed: %v", err)
	}

	expected := []string{
		"type Config struct",
		"type Service interface",
		"func main()",
		"func NewServer(int) *Server",
		"func (*Server) Start() error",
		"func (Server) Stop() (error, bool)",
		"type Handler func",
		"type MyInt int",
	}

	if len(defs) != len(expected) {
		t.Errorf("Expected %d definitions, got %d", len(expected), len(defs))
		for i, d := range defs {
			t.Logf("Got[%d]: %s", i, d)
		}
	}

	for i, d := range defs {
		if i < len(expected) && d != expected[i] {
			t.Errorf("Definition mismatch at %d: got %q, want %q", i, d, expected[i])
		}
	}
}
