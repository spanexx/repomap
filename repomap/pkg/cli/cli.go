package cli

import (
	"flag"
	"os"
)

// App represents the CLI application.
type App struct {
	Name    string
	Version string

	// Flags registry (Task 1.2.2)
	flagSet    *flag.FlagSet
	flagValues map[string]interface{} // Stores pointers to flag values

	// Future fields (to be implemented in subsequent tasks)
	// Output writer (Task 1.2.4)
	// Error handler (Task 1.2.12)
}

// NewApp creates a new instance of the CLI application.
func NewApp(name, version string) *App {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	// Redirect output to stderr by default, but can be customized later
	fs.SetOutput(os.Stderr)

	return &App{
		Name:       name,
		Version:    version,
		flagSet:    fs,
		flagValues: make(map[string]interface{}),
	}
}
