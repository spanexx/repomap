package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Flags holds the parsed flag values.
type Flags struct {
	values map[string]interface{}
}

// AddFlag adds a flag to the application.
// The type of the flag is inferred from the defaultValue.
// Supported types: string, int, bool, []string.
func (a *App) AddFlag(name string, usage string, defaultValue interface{}) *App {
	switch v := defaultValue.(type) {
	case string:
		val := a.flagSet.String(name, v, usage)
		a.flagValues[name] = val
	case int:
		val := a.flagSet.Int(name, v, usage)
		a.flagValues[name] = val
	case bool:
		val := a.flagSet.Bool(name, v, usage)
		a.flagValues[name] = val
	case []string:
		// Make a copy of defaults
		defaults := make([]string, len(v))
		copy(defaults, v)

		val := &stringSlice{values: defaults}
		a.flagSet.Var(val, name, usage)
		a.flagValues[name] = val
	default:
		// Panicking here is acceptable as this is a developer error during setup
		panic(fmt.Sprintf("unsupported flag type for flag '%s': %T", name, defaultValue))
	}
	return a
}

// Parse parses the arguments and returns the flag values.
func (a *App) Parse(args []string) (*Flags, error) {
	if err := a.flagSet.Parse(args); err != nil {
		return nil, err
	}

	return &Flags{
		values: a.flagValues,
	}, nil
}

// GetString returns the string value of a flag.
// Returns empty string if the flag does not exist or is not a string.
func (f *Flags) GetString(name string) string {
	if v, ok := f.values[name]; ok {
		if s, ok := v.(*string); ok {
			return *s
		}
	}
	return ""
}

// GetInt returns the int value of a flag.
// Returns 0 if the flag does not exist or is not an int.
func (f *Flags) GetInt(name string) int {
	if v, ok := f.values[name]; ok {
		if i, ok := v.(*int); ok {
			return *i
		}
	}
	return 0
}

// GetBool returns the bool value of a flag.
// Returns false if the flag does not exist or is not a bool.
func (f *Flags) GetBool(name string) bool {
	if v, ok := f.values[name]; ok {
		if b, ok := v.(*bool); ok {
			return *b
		}
	}
	return false
}

// GetStringSlice returns the []string value of a flag.
// Returns nil if the flag does not exist or is not a []string.
func (f *Flags) GetStringSlice(name string) []string {
	if v, ok := f.values[name]; ok {
		if s, ok := v.(*stringSlice); ok {
			return s.values
		}
	}
	return nil
}

// stringSlice implements flag.Value interface for []string flags.
type stringSlice struct {
	values  []string
	changed bool
}

// verify stringSlice implements flag.Value
var _ flag.Value = (*stringSlice)(nil)

func (s *stringSlice) String() string {
	return strings.Join(s.values, ",")
}

func (s *stringSlice) Set(value string) error {
	// Standard behavior for slice flags:
	// If the user provides the flag, it REPLACES the default value on the first occurrence.
	// Subsequent occurrences APPEND to the list.
	if !s.changed {
		s.values = []string{}
		s.changed = true
	}

	// If the flag is provided multiple times, we append.
	// However, if the user provides comma-separated values, we split them.
	// Standard convention for repeatable flags is usually multiple flags or comma-separated.
	// Let's support both for maximum flexibility.
	parts := strings.Split(value, ",")
	for _, part := range parts {
		s.values = append(s.values, strings.TrimSpace(part))
	}
	return nil
}
