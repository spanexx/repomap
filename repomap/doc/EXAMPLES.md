# CLI Framework Examples

## Minimal Tool

This example demonstrates a minimal tool that prints "Hello, World!" using the framework.

```go
package main

import (
    "fmt"
    "os"
    "github.com/spanexx/agents-cli/repomap/pkg/cli"
)

func main() {
    app := cli.NewApp("hello", "0.1.0")
    app.AddFlag("name", "Name to greet", "World")

    flags, err := app.Parse(os.Args[1:])
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    name := flags.GetString("name")
    fmt.Printf("Hello, %s!\n", name)
}
```

## Tool with Output Formatting

This example shows how to use the `pkg/output` package to support JSON and XML output.

```go
package main

import (
    "os"
    "github.com/spanexx/agents-cli/repomap/pkg/cli"
    "github.com/spanexx/agents-cli/repomap/pkg/output"
)

type Result struct {
    Message string `json:"message" xml:"message"`
    Count   int    `json:"count" xml:"count"`
}

func main() {
    app := cli.NewApp("counter", "1.0.0")
    app.AddFlag("format", "Output format (json|xml)", "json")
    app.AddFlag("count", "Count value", 1)

    flags, _ := app.Parse(os.Args[1:])
    
    writer, _ := output.NewWriter(flags.GetString("format"))
    
    res := Result{
        Message: "Success",
        Count:   flags.GetInt("count"),
    }
    
    writer.Write(res)
}
```
