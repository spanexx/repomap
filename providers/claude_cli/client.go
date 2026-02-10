package claude_cli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"llm-adapter/pkg/adapter"
	"llm-adapter/pkg/providers/auditlog"
)

type Provider struct {
	binaryPath   string
	model        string
	systemPrompt string
	workingDir   string
	timeout      time.Duration
}

func New() *Provider {
	p := &Provider{timeout: 120 * time.Second}
	if d := strings.TrimSpace(os.Getenv("CLAUDE_CWD")); d != "" {
		p.workingDir = d
	}
	if t := strings.TrimSpace(os.Getenv("CLAUDE_TIMEOUT_SECONDS")); t != "" {
		if parsed, err := time.ParseDuration(t + "s"); err == nil {
			p.timeout = parsed
		}
	}
	if bp := strings.TrimSpace(os.Getenv("CLAUDE_PATH")); bp != "" {
		p.binaryPath = bp
		return p
	}
	if lp, err := exec.LookPath("claude"); err == nil {
		p.binaryPath = lp
	}
	return p
}

func (p *Provider) Name() string { return "claude" }

func (p *Provider) SetModel(model string) {
	if strings.TrimSpace(model) != "" {
		p.model = strings.TrimSpace(model)
	}
}

func (p *Provider) SetSystemPrompt(prompt string) { p.systemPrompt = prompt }

func (p *Provider) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		p.timeout = timeout
	}
}

func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	if strings.TrimSpace(p.binaryPath) == "" {
		return "", fmt.Errorf("claude binary not found (set CLAUDE_PATH or ensure claude is in PATH)")
	}
	fullPrompt, err := buildPrompt(prompt, attachments, p.systemPrompt)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	start := time.Now()

	args := []string{"-p", "--output-format", "text"}
	if strings.TrimSpace(p.model) != "" {
		args = append(args, "--model", strings.TrimSpace(p.model))
	}

	cmd := exec.CommandContext(ctx, p.binaryPath, args...)
	cmd.Stdin = strings.NewReader(fullPrompt)
	if strings.TrimSpace(p.workingDir) != "" {
		cmd.Dir = p.workingDir
	}
	cmd.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	dur := time.Since(start)
	exit := 0
	if err != nil {
		exit = -1
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		}
	}
	if wd := strings.TrimSpace(cmd.Dir); wd == "" {
		wd = strings.TrimSpace(p.workingDir)
	}
	auditlog.LogExec(auditlog.ExecEvent{
		Provider:    p.Name(),
		Binary:      p.binaryPath,
		Dir:         cmd.Dir,
		Timeout:     p.timeout,
		Duration:    dur,
		ExitCode:    exit,
		TimedOut:    ctx.Err() != nil,
		StdoutBytes: stdout.Len(),
		StderrBytes: stderr.Len(),
		ErrorSummary: func() string {
			if err != nil {
				return err.Error()
			}
			return ""
		}(),
	})

	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("claude timeout")
		}
		errStr := strings.TrimSpace(stderr.String())
		outStr := strings.TrimSpace(stdout.String())
		if errStr != "" {
			return "", fmt.Errorf("claude failed: %v; stderr: %s", err, errStr)
		}
		if outStr != "" {
			return "", fmt.Errorf("claude failed: %v; stdout: %s", err, outStr)
		}
		return "", fmt.Errorf("claude failed: %w", err)
	}

	return strings.TrimSpace(stdout.String()), nil
}

func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	resp, err := p.Generate(prompt, attachments)
	if err != nil {
		return err
	}
	tokens <- resp
	return nil
}

func buildPrompt(prompt string, attachments []adapter.Attachment, systemPrompt string) (string, error) {
	var b strings.Builder
	sp := strings.TrimSpace(systemPrompt)
	if sp != "" {
		b.WriteString(sp)
		b.WriteString("\n\n")
	}
	for _, att := range attachments {
		switch att.Type {
		case adapter.AttachmentTypeText:
			name := att.Name
			if strings.TrimSpace(name) == "" {
				name = att.Path
			}
			b.WriteString("--- File: ")
			b.WriteString(name)
			b.WriteString(" ---\n")
			b.WriteString(att.Data)
			b.WriteString("\n--- End File ---\n\n")
		case adapter.AttachmentTypeImage:
			return "", fmt.Errorf("claude provider does not support image attachments")
		case adapter.AttachmentTypeFolder:
			return "", fmt.Errorf("claude provider does not support folder attachments")
		default:
			return "", fmt.Errorf("unsupported attachment type: %s", att.Type)
		}
	}
	b.WriteString(prompt)
	return b.String(), nil
}
