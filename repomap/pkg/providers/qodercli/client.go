package qodercli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
)

// Provider implements adapter.Provider using the external qodercli binary.
// qodercli is treated as a one-shot CLI: each prompt runs a new process.
//
// Configuration (optional) via env vars:
// - QODERCLI_PATH: override binary path
// - QODERCLI_TIMEOUT_SECONDS: command timeout (default: 60)
// - QODERCLI_CWD: working directory for qodercli execution (default: current process cwd)
//
// Notes:
// - Attachments: text attachments are appended to the prompt; image attachments are not supported.
// - Streaming: implemented as "send all at once".
type Provider struct {
	binaryPath   string
	model        string
	systemPrompt string
	workingDir   string
	timeout      time.Duration
}

func New() *Provider {
	p := &Provider{
		model:   "",
		timeout: 60 * time.Second,
	}

	if d := strings.TrimSpace(os.Getenv("QODERCLI_CWD")); d != "" {
		p.workingDir = d
	}

	if t := strings.TrimSpace(os.Getenv("QODERCLI_TIMEOUT_SECONDS")); t != "" {
		if parsed, err := time.ParseDuration(t + "s"); err == nil {
			p.timeout = parsed
		}
	}

	if bp := strings.TrimSpace(os.Getenv("QODERCLI_PATH")); bp != "" {
		p.binaryPath = bp
		return p
	}

	if lp, err := exec.LookPath("qodercli"); err == nil {
		p.binaryPath = lp
	}

	return p
}

func (p *Provider) Name() string {
	return "qodercli"
}

func (p *Provider) SetModel(model string) {
	if strings.TrimSpace(model) != "" {
		p.model = strings.TrimSpace(model)
	}
}

func (p *Provider) SetSystemPrompt(prompt string) {
	p.systemPrompt = prompt
}

func (p *Provider) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		p.timeout = timeout
	}
}

func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	if strings.TrimSpace(p.binaryPath) == "" {
		return "", fmt.Errorf("qodercli binary not found (set QODERCLI_PATH or ensure qodercli is in PATH)")
	}

	fullPrompt, err := buildPrompt(prompt, attachments, p.systemPrompt)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	start := time.Now()

	cmd := exec.CommandContext(ctx, p.binaryPath, "-p", fullPrompt)
	if strings.TrimSpace(p.workingDir) != "" {
		cmd.Dir = p.workingDir
	}
	cmd.Env = os.Environ()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	_ = time.Since(start) // Duration unused without auditlog

	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("qodercli timeout")
		}
		errStr := strings.TrimSpace(stderr.String())
		outStr := strings.TrimSpace(stdout.String())
		if errStr != "" && outStr != "" {
			return "", fmt.Errorf("qodercli failed: %v; stderr: %s; stdout: %s", err, errStr, outStr)
		}
		if errStr != "" {
			return "", fmt.Errorf("qodercli failed: %v; stderr: %s", err, errStr)
		}
		if outStr != "" {
			return "", fmt.Errorf("qodercli failed: %v; stdout: %s", err, outStr)
		}
		return "", fmt.Errorf("qodercli failed: %w", err)
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
			return "", fmt.Errorf("qodercli provider does not support image attachments")
		case adapter.AttachmentTypeFolder:
			// Folder attachments are expanded earlier into text attachments in most flows.
			// If they reach here, treat as unsupported to avoid silent data loss.
			return "", fmt.Errorf("qodercli provider does not support folder attachments")
		default:
			return "", fmt.Errorf("unsupported attachment type: %s", att.Type)
		}
	}

	b.WriteString(prompt)
	return b.String(), nil
}
