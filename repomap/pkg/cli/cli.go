package cli

import (
	"flag"
	"fmt"
	"os"
)

// App represents the CLI application.
type App struct {
	Name        string
	Version     string
	Description string
	UsageText   string
	Examples    []string

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

	app := &App{
		Name:       name,
		Version:    version,
		flagSet:    fs,
		flagValues: make(map[string]interface{}),
	}

	// Wire up custom help generation
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, app.GenerateHelp())
	}

	return app
}

// SetDescription sets the application description.
func (a *App) SetDescription(desc string) *App {
	a.Description = desc
	return a
}

// SetUsage sets the usage text.
func (a *App) SetUsage(usage string) *App {
	a.UsageText = usage
	return a
}

// AddExample adds an example to the help text.
func (a *App) AddExample(example string) *App {
	a.Examples = append(a.Examples, example)
	return a
}
