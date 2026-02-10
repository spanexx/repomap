package errors

import "fmt"

// Code represents a standardized error code.
type Code int

// Error represents a standardized error with a code and optional cause.
type Error struct {
	Code    Code
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// New creates a new Error.
func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

// Wrap creates a new Error wrapping an existing error.
func Wrap(code Code, message string, cause error) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// Specialized Error Constructors

func NewConfigError(message string, cause error) *Error {
	return Wrap(CodeConfigError, message, cause)
}

func NewParseError(message string, cause error) *Error {
	return Wrap(CodeParseError, message, cause)
}

func NewExecutionError(message string, cause error) *Error {
	return Wrap(CodeExecutionError, message, cause)
}

func NewValidationError(message string) *Error {
	return New(CodeValidationError, message)
}
