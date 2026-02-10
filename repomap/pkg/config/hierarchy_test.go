package config

import (
	"testing"

	"github.com/spanexx/agents-cli/repomap/pkg/cli"
)

func TestMergeConfigs(t *testing.T) {
	// 1. Defaults
	defaults := map[string]interface{}{
		"timeout": 30,
		"verbose": false,
		"output":  "text",
	}

	// 2. File Config (overrides timeout)
	fileCfg := &Config{
		Settings: map[string]interface{}{
			"timeout": 60,
			"theme":   "dark",
		},
	}

	// 3. Env Vars (overrides verbose)
	envCfg := map[string]string{
		"verbose": "true",
		"apikey":  "secret",
	}

	// 4. CLI Flags (overrides output)
	app := cli.NewApp("test", "1.0")
	app.AddFlag("output", "fmt", "json")
	// Simulate parsing to populate values
	flags, _ := app.Parse([]string{"--output", "xml"})

	// Merge
	merged := MergeConfigs(defaults, fileCfg, envCfg, flags)

	// Assertions
	check := func(key string, expected interface{}) {
		val, ok := merged.Settings[key]
		if !ok {
			t.Errorf("missing key %s", key)
			return
		}
		if val != expected {
			t.Errorf("key %s: expected %v, got %v", key, expected, val)
		}
	}

	// timeout should be 60 (File overrides Default)
	check("timeout", 60)

	// verbose should be "true" (Env overrides Default)
	// Note: Env vars are strings, so it overrides boolean false with string "true"
	check("verbose", "true")

	// output should be "xml" (CLI overrides Default/File/Env)
	check("output", "xml")

	// theme should be "dark" (from File)
	check("theme", "dark")

	// apikey should be "secret" (from Env)
	check("apikey", "secret")
}
