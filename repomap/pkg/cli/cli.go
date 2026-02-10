package cli

// App represents the CLI application.
type App struct {
	Name    string
	Version string

	// Future fields (to be implemented in subsequent tasks)
	// Flags registry (Task 1.2.2)
	// Output writer (Task 1.2.4)
	// Error handler (Task 1.2.12)
}

// NewApp creates a new instance of the CLI application.
func NewApp(name, version string) *App {
	return &App{
		Name:    name,
		Version: version,
	}
}
