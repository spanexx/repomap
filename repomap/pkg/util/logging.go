package util

import (
	"fmt"
	"io"
	"os"
)

// Logger defines the interface for logging operations.
type Logger interface {
	Info(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	SetVerbose(verbose bool)
}

// StandardLogger implements Logger using provided writers (e.g., stdout/stderr).
type StandardLogger struct {
	out     io.Writer
	err     io.Writer
	verbose bool
}

// NewLogger creates a new StandardLogger.
// If out or err are nil, they default to os.Stdout and os.Stderr.
func NewLogger(out, err io.Writer, verbose bool) *StandardLogger {
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}
	return &StandardLogger{
		out:     out,
		err:     err,
		verbose: verbose,
	}
}

func (l *StandardLogger) Info(msg string, args ...interface{}) {
	fmt.Fprintf(l.out, msg+"\n", args...)
}

func (l *StandardLogger) Debug(msg string, args ...interface{}) {
	if l.verbose {
		fmt.Fprintf(l.out, "[DEBUG] "+msg+"\n", args...)
	}
}

func (l *StandardLogger) Warn(msg string, args ...interface{}) {
	fmt.Fprintf(l.err, "[WARN] "+msg+"\n", args...)
}

func (l *StandardLogger) Error(msg string, args ...interface{}) {
	fmt.Fprintf(l.err, "[ERROR] "+msg+"\n", args...)
}

func (l *StandardLogger) SetVerbose(verbose bool) {
	l.verbose = verbose
}
