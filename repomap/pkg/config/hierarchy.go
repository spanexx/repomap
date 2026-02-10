package config

import (
	"fmt"
	"github.com/spanexx/agents-cli/repomap/pkg/cli"
)

// MergeConfigs merges configurations from multiple sources.
// Priority (lowest to highest):
// 1. Defaults (provided as map)
// 2. File Configuration (*Config)
// 3. Environment Variables (map[string]string)
// 4. CLI Flags (*cli.Flags)
func MergeConfigs(defaults map[string]interface{}, fileCfg *Config, envCfg map[string]string, flags *cli.Flags) *Config {
	final := make(map[string]interface{})

	// 1. Defaults
	for k, v := range defaults {
		final[k] = v
	}

	// 2. File Config
	if fileCfg != nil {
		for k, v := range fileCfg.Settings {
			final[k] = v
		}
	}

	// 3. Environment Variables
	for k, v := range envCfg {
		final[k] = v
	}

	// 4. CLI Flags
	if flags != nil {
		// Use getters to retrieve typed values, but ONLY for visited flags.
		// If we use GetValues() (all flags), defaults will overwrite Env/File.
		flagVals := flags.GetVisitedValues()
		for k := range flagVals {
			v := flagVals[k]

			switch val := v.(type) {
			case *string:
				final[k] = *val
			case *int:
				final[k] = *val
			case *bool:
				final[k] = *val
			case fmt.Stringer:
				final[k] = val.String()
			default:
				final[k] = v
			}
		}
	}

	return &Config{Settings: final}
}
