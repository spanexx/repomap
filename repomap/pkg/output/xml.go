package output

import (
	"encoding/xml"
	"fmt"
	"os"
)

// XMLWriter implements the Writer interface for XML formatting.
type XMLWriter struct {
	// We might add configuration here later (e.g., indentation)
}

// NewXMLWriter creates a new XMLWriter.
func NewXMLWriter() Writer {
	return &XMLWriter{}
}

// Write writes the data as XML to stdout.
func (w *XMLWriter) Write(data interface{}) error {
	output, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal XML: %w", err)
	}

	// Add XML header
	fmt.Println(xml.Header + string(output))
	return nil
}

// WriteString writes a raw string to stdout.
func (w *XMLWriter) WriteString(s string) error {
	_, err := fmt.Fprintln(os.Stdout, s)
	return err
}
