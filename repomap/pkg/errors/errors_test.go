package errors

import (
	"fmt"
	"testing"
)

func TestErrorCreation(t *testing.T) {
	err := New(CodeConfigError, "config error")
	if err.Code != CodeConfigError {
		t.Errorf("expected code %d, got %d", CodeConfigError, err.Code)
	}
	if err.Message != "config error" {
		t.Errorf("expected message 'config error', got '%s'", err.Message)
	}
}

func TestErrorWrapping(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := Wrap(CodeExecutionError, "exec error", cause)

	if err.Cause != cause {
		t.Error("expected cause to be preserved")
	}

	expectedStr := "[102] exec error: root cause"
	if err.Error() != expectedStr {
		t.Errorf("expected string '%s', got '%s'", expectedStr, err.Error())
	}
}

func TestSpecializedErrors(t *testing.T) {
	err := NewValidationError("invalid input")
	if err.Code != CodeValidationError {
		t.Errorf("expected validation code %d, got %d", CodeValidationError, err.Code)
	}
}
