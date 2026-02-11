package util

import (
	"bytes"
	"strings"
	"testing"
)

func TestStandardLogger_Info(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	logger := NewLogger(&out, &err, false)

	logger.Info("Hello %s", "World")

	expected := "Hello World\n"
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
	if err.String() != "" {
		t.Errorf("Expected empty error output, got %q", err.String())
	}
}

func TestStandardLogger_Debug(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	logger := NewLogger(&out, &err, false)

	// Non-verbose: Should not log
	logger.Debug("Secret message")
	if out.String() != "" {
		t.Errorf("Expected no output in non-verbose mode, got %q", out.String())
	}

	// Verbose: Should log
	logger.SetVerbose(true)
	logger.Debug("Secret message")
	expected := "[DEBUG] Secret message\n"
	if out.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, out.String())
	}
}

func TestStandardLogger_Warn(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	logger := NewLogger(&out, &err, false)

	logger.Warn("Watch out!")

	expected := "[WARN] Watch out!\n"
	if err.String() != expected {
		t.Errorf("Expected error output %q, got %q", expected, err.String())
	}
	if out.String() != "" {
		t.Errorf("Expected empty stdout, got %q", out.String())
	}
}

func TestStandardLogger_Error(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	logger := NewLogger(&out, &err, false)

	logger.Error("Something failed")

	expected := "[ERROR] Something failed\n"
	if err.String() != expected {
		t.Errorf("Expected error output %q, got %q", expected, err.String())
	}
}

func TestNewLoggerHelpers(t *testing.T) {
	// Identify if os.Stdout/Stderr defaults are used when nil is passed
	// We can't easily capture os.Stdout here without more complex mocking,
	// but we can ensure it facilitates nil arguments without panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("NewLogger panicked with nil arguments")
		}
	}()
	NewLogger(nil, nil, true)
}

func TestLoggerFormatting(t *testing.T) {
	var out bytes.Buffer
	logger := NewLogger(&out, nil, true)

	logger.Debug("Value: %d", 42)
	if !strings.Contains(out.String(), "Value: 42") {
		t.Errorf("Formatting failed, got %q", out.String())
	}
}
