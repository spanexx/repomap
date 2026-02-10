package output

import (
	"fmt"
)

// NewWriter creates a new Writer implementation based on the format string.
// Supported formats: xml, json, text.
// Currently returns an error as no formatters are implemented yet.
func NewWriter(format string) (Writer, error) {
	switch format {
	// Future tasks will implement these cases:
	// case "xml": return NewXMLWriter(...)
	// case "json": return NewJSONWriter(...)
	// case "text": return NewTextWriter(...)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
