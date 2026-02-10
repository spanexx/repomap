package qwen_cli

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
	argsTemplate string
}

func New() *Provider {
	p := &Provider{timeout: 120 * time.Second}
	if d := strings.TrimSpace(os.Getenv("QWENCLI_CWD")); d != "" {
		p.workingDir = d
	}
	if t := strings.TrimSpace(os.Getenv("QWENCLI_TIMEOUT_SECONDS")); t != "" {
		if parsed, err := time.ParseDuration(t + "s"); err == nil {
			p.timeout = parsed
		}
	}
	p.argsTemplate = strings.TrimSpace(os.Getenv("QWENCLI_ARGS"))
	if bp := strings.TrimSpace(os.Getenv("QWENCLI_PATH")); bp != "" {
		p.binaryPath = bp
		return p
	}
	if lp, err := exec.LookPath("qwen"); err == nil {
		p.binaryPath = lp
	}
	return p
}

func (p *Provider) Name() string { return "qwen-cli" }

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
		return "", fmt.Errorf("qwen CLI binary not found (set QWENCLI_PATH or ensure qwen is in PATH)")
	}
	fullPrompt, err := buildPrompt(prompt, attachments, p.systemPrompt)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	start := time.Now()

	args := buildArgs(p.argsTemplate, p.model, fullPrompt)

	cmd := exec.CommandContext(ctx, p.binaryPath, args...)
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
			return "", fmt.Errorf("qwen-cli timeout")
		}
		errStr := strings.TrimSpace(stderr.String())
		outStr := strings.TrimSpace(stdout.String())
		if errStr != "" {
			return "", fmt.Errorf("qwen-cli failed: %v; stderr: %s", err, errStr)
		}
		if outStr != "" {
			return "", fmt.Errorf("qwen-cli failed: %v; stdout: %s", err, outStr)
		}
		return "", fmt.Errorf("qwen-cli failed: %w", err)
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

func buildArgs(argsTemplate string, model string, fullPrompt string) []string {
	// Qwen CLI variants differ. We keep this adapter flexible via QWENCLI_ARGS.
	// - If QWENCLI_ARGS is empty: we pass the prompt as a single positional argument.
	// - If QWENCLI_ARGS contains {prompt}: we replace it with the prompt.
	// - If QWENCLI_ARGS contains {model}: we replace it with the model.
	// - If {prompt} is missing, we append prompt at the end.
	tpl := strings.TrimSpace(argsTemplate)
	if tpl == "" {
		return []string{fullPrompt}
	}

	fields := strings.Fields(tpl)
	out := make([]string, 0, len(fields)+1)
	hasPrompt := false
	for _, f := range fields {
		if strings.Contains(f, "{prompt}") {
			hasPrompt = true
			out = append(out, strings.ReplaceAll(f, "{prompt}", fullPrompt))
			continue
		}
		if strings.Contains(f, "{model}") {
			out = append(out, strings.ReplaceAll(f, "{model}", strings.TrimSpace(model)))
			continue
		}
		out = append(out, f)
	}
	if !hasPrompt {
		out = append(out, fullPrompt)
	}
	return out
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
			return "", fmt.Errorf("qwen-cli provider does not support image attachments")
		case adapter.AttachmentTypeFolder:
			return "", fmt.Errorf("qwen-cli provider does not support folder attachments")
		default:
			return "", fmt.Errorf("unsupported attachment type: %s", att.Type)
		}
	}
	b.WriteString(prompt)
	return b.String(), nil
}
