package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_JSON(t *testing.T) {
	// Create a temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".testrc")

	content := `{"output": "xml", "verbose": true}`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig([]string{configPath})
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if val := cfg.GetString("output"); val != "xml" {
		t.Errorf("expected output=xml, got %s", val)
	}

	if val := cfg.GetBool("verbose"); !val {
		t.Error("expected verbose=true")
	}
}

func TestLoadConfig_Missing(t *testing.T) {
	// Should return empty config, no error
	cfg, err := LoadConfig([]string{"/non/existent/path"})
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}

	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if len(cfg.Settings) != 0 {
		t.Errorf("expected empty settings, got %v", cfg.Settings)
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".badrc")

	content := `{invalid-json`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfig([]string{configPath})
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestConfig_Getters_Defaults(t *testing.T) {
	cfg := &Config{Settings: make(map[string]interface{})}

	if val := cfg.GetString("missing"); val != "" {
		t.Errorf("expected empty string default, got %s", val)
	}

	if val := cfg.GetBool("missing"); val != false {
		t.Errorf("expected false default, got %v", val)
	}
}
