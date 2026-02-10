package errors

import "fmt"

// Message templates
const (
	MsgConfigLoadFailed = "failed to load configuration"
	MsgParseFailed      = "failed to parse input"
	MsgValidationFailed = "validation failed"
)

// FormatMessage formats a message with arguments.
func FormatMessage(template string, args ...interface{}) string {
	return fmt.Sprintf(template, args...)
}
