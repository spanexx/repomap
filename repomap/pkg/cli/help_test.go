package cli

import (
	"strings"
	"testing"
)

func TestGenerateHelp(t *testing.T) {
	app := NewApp("test-app", "1.0.0")

	// Setup app details
	app.SetDescription("A test application for CLI framework")
	app.SetUsage("test-app [options] <command>")

	// Add flags
	app.AddFlag("verbose", "Enable verbose output", false)
	app.AddFlag("output", "Output format (xml|json)", "xml")

	// Add examples
	app.AddExample("test-app --verbose run")
	app.AddExample("test-app --output json scan")

	// Generate help
	helpText := app.GenerateHelp()

	// Verify content
	expectedSections := []string{
		"test-app v1.0.0",
		"A test application for CLI framework",
		"Usage:",
		"  test-app [options] <command>",
		"Options:",
		"  --verbose", "Enable verbose output", "(default: false)",
		"  --output", "Output format (xml|json)", "(default: xml)",
		"Examples:",
		"  test-app --verbose run",
		"  test-app --output json scan",
	}

	for _, section := range expectedSections {
		if !strings.Contains(helpText, section) {
			t.Errorf("Help text missing expected section: %q", section)
		}
	}
}

func TestGenerateHelp_Minimal(t *testing.T) {
	app := NewApp("minimal-app", "0.0.1")

	helpText := app.GenerateHelp()

	if !strings.Contains(helpText, "minimal-app v0.0.1") {
		t.Error("Help text missing name/version")
	}

	if !strings.Contains(helpText, "Usage:") {
		t.Error("Help text missing Usage section")
	}

	// Should contain default usage pattern
	if !strings.Contains(helpText, "minimal-app [options]") {
		t.Error("Help text missing default usage pattern")
	}

	// Should NOT contain Options section as no flags added
	if strings.Contains(helpText, "Options:") {
		t.Error("Help text should not contain Options section when no flags are present")
	}

	// Should NOT contain Examples section
	if strings.Contains(helpText, "Examples:") {
		t.Error("Help text should not contain Examples section when no examples are present")
	}
}
