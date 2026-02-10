package output

import (
	"fmt"
)

// NewWriter creates a new Writer implementation based on the format string.
// Supported formats: xml, json, text.
func NewWriter(format string) (Writer, error) {
	switch format {
	case "xml":
		return NewXMLWriter(), nil
	// Future tasks will implement these cases:
	// case "json": return NewJSONWriter(...)
	// case "text": return NewTextWriter(...)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
