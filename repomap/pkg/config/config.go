package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the configuration data.
// For now, it's a simple map, but can be typed later.
type Config struct {
	Settings map[string]interface{}
}

// LoadConfig attempts to load configuration from the given paths.
// It stops at the first successful load.
// If no file is found, it returns an empty config and no error.
// Supported format: JSON (e.g., .repomaprc as JSON object).
func LoadConfig(paths []string) (*Config, error) {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
		}

		var settings map[string]interface{}
		// We only support JSON for now to avoid external dependencies
		if err := json.Unmarshal(data, &settings); err != nil {
			// If JSON fails, we could try other parsers if we had them.
			// For MVP, we treat it as an error if the file exists but isn't valid JSON.
			return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
		}

		return &Config{Settings: settings}, nil
	}

	// No config file found, return empty config
	return &Config{Settings: make(map[string]interface{})}, nil
}

// GetString returns a string value from the config.
func (c *Config) GetString(key string) string {
	if v, ok := c.Settings[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetBool returns a bool value from the config.
func (c *Config) GetBool(key string) bool {
	if v, ok := c.Settings[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// DefaultPaths returns standard configuration paths for a tool.
// e.g. ./.toolrc, ~/.toolrc
func DefaultPaths(toolName string) []string {
	paths := []string{
		fmt.Sprintf("./.%src", toolName),
	}

	home, err := os.UserHomeDir()
	if err == nil {
		paths = append(paths, filepath.Join(home, fmt.Sprintf(".%src", toolName)))
	}

	return paths
}
