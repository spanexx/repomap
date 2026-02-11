package parsing

import (
	"bufio"
	"os"
	"strings"
)

// GenericExtractor provides a basic line-based definition extraction for unsupported languages.
type GenericExtractor struct {
	DefKeywords    []string
	ImportKeywords []string
}

func (e *GenericExtractor) ExtractImports(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var imports []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || isComment(line) {
			continue
		}

		for _, kw := range e.ImportKeywords {
			if strings.HasPrefix(line, kw+" ") {
				// Try to find something in quotes
				start := strings.Index(line, "\"")
				if start == -1 {
					start = strings.Index(line, "'")
				}
				if start != -1 {
					remaining := line[start+1:]
					end := strings.Index(remaining, "\"")
					if end == -1 {
						end = strings.Index(remaining, "'")
					}
					if end != -1 {
						imports = append(imports, remaining[:end])
					}
				}
				break
			}
		}
	}
	return imports, scanner.Err()
}

func (e *GenericExtractor) ExtractDefinitions(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var definitions []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || isComment(line) {
			continue
		}

		for _, kw := range e.DefKeywords {
			if strings.HasPrefix(line, kw+" ") {
				definitions = append(definitions, line)
				break
			}
		}
	}
	return definitions, scanner.Err()
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "#")
}

func init() {
	generic := &GenericExtractor{
		DefKeywords:    []string{"class", "function", "func", "type", "interface", "export", "const", "var", "let", "def", "struct"},
		ImportKeywords: []string{"import", "from", "require"},
	}

	// Register for common web/scripting extensions as fallback
	exts := []string{".ts", ".tsx", ".js", ".jsx", ".py", ".rs", ".java", ".cpp", ".c", ".h", ".cs"}
	for _, ext := range exts {
		DefaultRegistry.Register(ext, generic)
	}
}
