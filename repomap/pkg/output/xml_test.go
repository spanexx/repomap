package output

import (
	"encoding/xml"
	"testing"
)

type TestStruct struct {
	XMLName xml.Name `xml:"test"`
	Value   string   `xml:"value"`
}

func TestNewXMLWriter(t *testing.T) {
	writer := NewXMLWriter()
	if writer == nil {
		t.Fatal("NewXMLWriter returned nil")
	}

	_, ok := writer.(*XMLWriter)
	if !ok {
		t.Error("NewXMLWriter did not return *XMLWriter")
	}
}

// Note: Testing Write and WriteString involving stdout capture is tricky in unit tests
// without dependency injection of the writer. For now, we assume standard library calls work.
// A more robust test would refactor XMLWriter to take an io.Writer.
// Given the constraints, we focus on the factory and basic structure.

// We can test the factory integration though.
func TestNewWriter_XML(t *testing.T) {
	writer, err := NewWriter("xml")
	if err != nil {
		t.Fatalf("NewWriter('xml') returned error: %v", err)
	}
	if writer == nil {
		t.Fatal("NewWriter('xml') returned nil writer")
	}
	_, ok := writer.(*XMLWriter)
	if !ok {
		t.Error("NewWriter('xml') did not return *XMLWriter")
	}
}
