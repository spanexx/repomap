package errors

import (
	"testing"
)

func TestFormatMessage(t *testing.T) {
	msg := FormatMessage(MsgFileNotFound, "config.json")
	expected := "file not found: config.json"

	if msg != expected {
		t.Errorf("expected '%s', got '%s'", expected, msg)
	}

	msg2 := FormatMessage(MsgConfigLoadFailed)
	if msg2 != MsgConfigLoadFailed {
		t.Errorf("expected '%s', got '%s'", MsgConfigLoadFailed, msg2)
	}
}
