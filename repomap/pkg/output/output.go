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
	case "json":
		return NewJSONWriter(), nil
	case "text":
		return NewTextWriter(), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
