package output

import (
	"testing"
)

func TestNewWriter(t *testing.T) {
	tests := []struct {
		format    string
		expectErr bool
	}{
		{"xml", false},    // Implemented
		{"json", false},   // Implemented
		{"text", true},    // Not implemented yet
		{"unknown", true}, // Unknown format
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			writer, err := NewWriter(tt.format)
			if tt.expectErr {
				if err == nil {
					t.Errorf("expected error for format %s, got nil", tt.format)
				}
				if writer != nil {
					t.Errorf("expected nil writer for format %s, got %v", tt.format, writer)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for format %s, got %v", tt.format, err)
				}
				if writer == nil {
					t.Errorf("expected writer for format %s, got nil", tt.format)
				}
			}
		})
	}
}
