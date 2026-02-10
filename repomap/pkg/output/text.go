package output

import (
	"fmt"
	"os"
)

// TextWriter implements the Writer interface for plain text/table formatting.
type TextWriter struct {
	// Configuration options could be added here (e.g., table columns)
}

// NewTextWriter creates a new TextWriter.
func NewTextWriter() Writer {
	return &TextWriter{}
}

// Write writes the data as text to stdout.
// For complex structures, it uses basic Go formatting (%+v).
// Future improvements can add reflection-based table formatting.
func (w *TextWriter) Write(data interface{}) error {
	// Simple implementation for now: just print the data
	// If data is a string, print it directly
	if s, ok := data.(string); ok {
		return w.WriteString(s)
	}

	// Otherwise use formatted print
	_, err := fmt.Printf("%+v\n", data)
	return err
}

// WriteString writes a raw string to stdout.
func (w *TextWriter) WriteString(s string) error {
	_, err := fmt.Fprintln(os.Stdout, s)
	return err
}
