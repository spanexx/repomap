package output

import (
	"testing"
)

func TestNewJSONWriter(t *testing.T) {
	writer := NewJSONWriter()
	if writer == nil {
		t.Fatal("NewJSONWriter returned nil")
	}

	_, ok := writer.(*JSONWriter)
	if !ok {
		t.Error("NewJSONWriter did not return *JSONWriter")
	}
}

func TestNewWriter_JSON(t *testing.T) {
	writer, err := NewWriter("json")
	if err != nil {
		t.Fatalf("NewWriter('json') returned error: %v", err)
	}
	if writer == nil {
		t.Fatal("NewWriter('json') returned nil writer")
	}
	_, ok := writer.(*JSONWriter)
	if !ok {
		t.Error("NewWriter('json') did not return *JSONWriter")
	}
}
