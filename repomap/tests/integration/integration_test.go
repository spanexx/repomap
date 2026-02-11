package main_test

import (
	"bytes"
	"encoding/xml"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration_SimpleFixture(t *testing.T) {
	// Build the binary first
	cmdBuild := exec.Command("go", "build", "-o", "../../repomap_test_bin", "./cmd/repomap")
	cmdBuild.Dir = "../.." // root/repomap
	if out, err := cmdBuild.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, out)
	}
	defer os.Remove("../../repomap_test_bin")

	// Verify fixture exists
	fixturePath, _ := filepath.Abs("../../tests/fixtures/simple")
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		t.Fatalf("Fixture not found at %s", fixturePath)
	}

	// Run repomap against the fixture
	cmd := exec.Command("../../repomap_test_bin", "--root", fixturePath, "--output", "xml")
	cmd.Dir = "../.."
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Repomap run failed: %v\nStderr: %s", err, stderr.String())
	}

	output := stdout.String()

	// Basic XML validation
	if !strings.Contains(output, "<repomap>") {
		t.Errorf("Output does not contain <repomap> tag")
	}
	if !strings.Contains(output, "main.go") {
		t.Errorf("Output does not contain main.go")
	}
	if !strings.Contains(output, "func main()") {
		t.Errorf("Output does not contain func main()")
	}

	// Verify XML structure validity
	var xmlData struct {
		XMLName xml.Name `xml:"repomap"`
		Files   []struct {
			Path string `xml:"path,attr"`
		} `xml:"file"`
	}

	// Strip XML header for Unmarshal if needed, or Unmarshal parses it fine.
	if err := xml.Unmarshal([]byte(output), &xmlData); err != nil {
		t.Errorf("Failed to unmarshal XML output: %v", err)
	}

	if len(xmlData.Files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(xmlData.Files))
	} else {
		if xmlData.Files[0].Path != "main.go" { // Rel path from root
			// Wait, filepath.Rel might behave differently depending on how we constructed it.
			// In main.go: relPath, _ := filepath.Rel(absRoot, path)
			// fixturePath is absolute. file is fixturePath/main.go.
			// Rel should be "main.go".
			t.Errorf("Expected path 'main.go', got '%s'", xmlData.Files[0].Path)
		}
	}
}

func TestFramework_Flags(t *testing.T) {
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	// Test --version
	cmd := exec.Command(binPath, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run --version: %v", err)
	}
	if !strings.Contains(string(out), "repomap version") {
		t.Errorf("Expected version output, got: %s", out)
	}

	// Test --help
	cmd = exec.Command(binPath, "--help")
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run --help: %v", err)
	}
	if !strings.Contains(string(out), "Usage:") {
		t.Errorf("Expected help usage, got: %s", out)
	}
}

func TestFramework_OutputFormats(t *testing.T) {
	binPath := buildBinary(t)
	defer os.Remove(binPath)

	fixturePath, _ := filepath.Abs("../../tests/fixtures/simple")

	// Test JSON output
	cmd := exec.Command(binPath, "--root", fixturePath, "--output", "json")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to run JSON output: %v", err)
	}

	if !strings.Contains(string(out), "\"files\":") {
		t.Errorf("Expected JSON output containing 'files', got: %s", out)
	}
}

func buildBinary(t *testing.T) string {
	binName := filepath.Join(os.TempDir(), "repomap_test_"+t.Name())
	if strings.Contains(t.Name(), "/") {
		binName = filepath.Join(os.TempDir(), "repomap_test_bin")
	}

	cmd := exec.Command("go", "build", "-o", binName, "../../cmd/repomap")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, out)
	}
	return binName
}
