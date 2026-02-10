package config

import (
	"os"
	"strings"
)

// GetEnv retrieves an environment variable with a default value.
// The key should be the full environment variable name.
func GetEnv(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// LoadEnv loads all environment variables matching the prefix into a map.
// The keys in the map will be lowercased and the prefix removed.
// e.g. REPOMAP_VERBOSE=true -> {"verbose": "true"}
func LoadEnv(prefix string) map[string]string {
	envVars := make(map[string]string)

	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key, value := pair[0], pair[1]
		if strings.HasPrefix(key, prefix) {
			// Strip prefix and lowercase
			newKey := strings.ToLower(strings.TrimPrefix(key, prefix))
			envVars[newKey] = value
		}
	}

	return envVars
}
