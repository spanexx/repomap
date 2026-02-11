# Repomap CLI Framework Guide

This guide provides an overview of the Repomap CLI framework, designed to build consistent and robust CLI tools for the agent ecosystem.

## Components

The framework is divided into several reusable packages:

### 1. `pkg/cli` - Application Structure
Handles application initialization, flag parsing, and help text generation.

```go
import "github.com/spanexx/agents-cli/repomap/pkg/cli"

func main() {
    app := cli.NewApp("my-tool", "1.0.0")
    app.AddFlag("root", "Root directory", ".")
    
    flags, err := app.Parse(os.Args[1:])
    // ...
}
```

### 2. `pkg/config` - Configuration Loading
Loads configuration from files (JSON) and environment variables.

```go
import "github.com/spanexx/agents-cli/repomap/pkg/config"

paths := config.DefaultPaths("my-tool")
cfg, err := config.LoadConfig(paths)
```

### 3. `pkg/output` - Output Formatting
Standardizes output formats (XML, JSON, Text) for better agent consumption.

```go
import "github.com/spanexx/agents-cli/repomap/pkg/output"

writer, err := output.NewWriter("json")
writer.Write(myStruct)
```

### 4. `pkg/errors` - Error Handling
Provides standardized error types and codes (implementation details in `pkg/errors`).

### 5. `pkg/util` - Utilities
Shared helpers for logging, paths, and more.

```go
import "github.com/spanexx/agents-cli/repomap/pkg/util"

logger := util.NewLogger(os.Stdout, os.Stderr, verbose)
logger.Info("Starting...")
```

## Workflow

1.  **Initialize App**: Create a `cli.App` instance.
2.  **Define Flags**: Register flags using `app.AddFlag`.
3.  **Parse Arguments**: Call `app.Parse` to process CLI args.
4.  **Load Config**: value from `pkg/config` (optional merge with flags).
5.  **Setup Logger**: Initialize `util.Logger`.
6.  **Setup Output**: Initialize `output.Writer`.
7.  **Run Logic**: Execute tool logic using inputs.

## Extensibility

-   **New Flags**: Add types to `pkg/cli/flags.go`
-   **New Formats**: Implement `output.Writer` interface in `pkg/output`.
