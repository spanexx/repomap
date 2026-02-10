package output

import (
	"testing"
)

func TestNewTextWriter(t *testing.T) {
	writer := NewTextWriter()
	if writer == nil {
		t.Fatal("NewTextWriter returned nil")
	}

	_, ok := writer.(*TextWriter)
	if !ok {
		t.Error("NewTextWriter did not return *TextWriter")
	}
}

func TestNewWriter_Text(t *testing.T) {
	writer, err := NewWriter("text")
	if err != nil {
		t.Fatalf("NewWriter('text') returned error: %v", err)
	}
	if writer == nil {
		t.Fatal("NewWriter('text') returned nil writer")
	}
	_, ok := writer.(*TextWriter)
	if !ok {
		t.Error("NewWriter('text') did not return *TextWriter")
	}
}
