package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONWriter implements the Writer interface for JSON formatting.
type JSONWriter struct {
	// Configuration options could be added here
}

// NewJSONWriter creates a new JSONWriter.
func NewJSONWriter() Writer {
	return &JSONWriter{}
}

// Write writes the data as indented JSON to stdout.
func (w *JSONWriter) Write(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return nil
}

// WriteString writes a raw string to stdout.
func (w *JSONWriter) WriteString(s string) error {
	_, err := fmt.Fprintln(os.Stdout, s)
	return err
}
