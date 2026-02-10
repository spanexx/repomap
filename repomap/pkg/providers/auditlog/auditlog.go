package auditlog

import (
	"fmt"
	"log"
	"time"
)

type ExecEvent struct {
	Provider     string
	Binary       string
	Dir          string
	Timeout      time.Duration
	Duration     time.Duration
	ExitCode     int
	TimedOut     bool
	StdoutBytes  int
	StderrBytes  int
	ErrorSummary string
}

// LogExec records non-sensitive execution metadata for external CLI providers.
// IMPORTANT: Do not include user prompts, attachment contents, or environment variables.
func LogExec(ev ExecEvent) {
	log.Printf(
		"AUDIT external_cli_exec provider=%s binary=%q dir=%q timeout=%s duration=%s exit=%d timed_out=%t stdout_bytes=%d stderr_bytes=%d err=%q",
		ev.Provider,
		ev.Binary,
		ev.Dir,
		ev.Timeout,
		ev.Duration,
		ev.ExitCode,
		ev.TimedOut,
		ev.StdoutBytes,
		ev.StderrBytes,
		fmt.Sprintf("%s", ev.ErrorSummary),
	)
}
