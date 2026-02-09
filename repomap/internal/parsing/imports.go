package parsing

import (
	"go/parser"
	"go/token"
	"strconv"
)

// ExtractImports parses a Go file and returns a list of imported package paths.
func ExtractImports(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	// Parse only imports for performance
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	var imports []string
	if node.Imports != nil {
		for _, imp := range node.Imports {
			if imp.Path != nil {
				// Import path is a string literal, e.g. "fmt" or "github.com/pkg/errors"
				// We need to unquote it.
				path, err := strconv.Unquote(imp.Path.Value)
				if err != nil {
					// Fallback to raw value if unquote fails (should strictly not happen for valid Go files)
					path = imp.Path.Value
				}
				imports = append(imports, path)
			}
		}
	}

	return imports, nil
}
