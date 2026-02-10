package output

// Writer defines the interface for output formatting.
type Writer interface {
	// Write formats and writes the given data structure.
	// The implementation should handle type assertion and formatting.
	Write(data interface{}) error

	// WriteString writes a raw string to the output.
	WriteString(s string) error
}
