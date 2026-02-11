package gemini_cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
)

// Provider wrappers the external 'gemini' CLI tool.
type Provider struct {
	binaryPath string
	model      string
	timeout    time.Duration
}

// New creates a new Gemini CLI provider.
func New(ctx context.Context) (*Provider, error) {
	p := &Provider{
		timeout: 120 * time.Second,
	}

	// Allow override via env
	if bp := os.Getenv("GEMINI_CLI_PATH"); bp != "" {
		p.binaryPath = bp
	} else {
		// Look in PATH
		path, err := exec.LookPath("gemini")
		if err != nil {
			// Fallback or error? For now, default to "gemini" and let it fail at runtime if missing
			p.binaryPath = "gemini"
		} else {
			p.binaryPath = path
		}
	}

	return p, nil
}

func (p *Provider) Init(ctx context.Context) error {
	// Check if binary exists/runnable
	cmd := exec.CommandContext(ctx, p.binaryPath, "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gemini binary not found or not executable at %s", p.binaryPath)
	}
	return nil
}

func (p *Provider) Name() string {
	return "gemini-cli"
}

func (p *Provider) SetModel(model string) {
	p.model = model
}

func (p *Provider) SetSystemPrompt(prompt string) {
	// The external CLI might not support system prompt arg directly in one-shot mode easily
	// without flags, but we can prepend it to the prompt.
	// For now, we'll ignore it or handle it in Generate.
}

// Generate executes the gemini CLI with the given prompt.
func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	fullPrompt, err := p.buildPrompt(prompt, attachments)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	// Construct command: gemini --output-format json --prompt <full_prompt>
	// We use --prompt to pass it as an argument, avoiding interactive mode.
	args := []string{"--output-format", "json", "--prompt", fullPrompt}

	if p.model != "" {
		args = append(args, "--model", p.model)
	}

	cmd := exec.CommandContext(ctx, p.binaryPath, args...)
	cmd.Env = os.Environ()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("gemini execution failed: %v, stderr: %s", err, stderr.String())
	}

	// Parse JSON output
	// The CLI output is expected to be a JSON string or object.
	// Based on the user's experiment, providing --output-format json should make it structured.
	// However, if it returns a raw string in JSON format, we might need to unquote or parse.
	// Let's assume it returns the model's response directly, possibly wrapped.
	// We'll return the raw stdout for now, or attempt to parse if it wraps it in a known structure.
	// *Correction*: The CLI help says "--output-format json".
	// Let's try to unmarshal it as a generic wrapper if appropriate,
	// or return the raw text if it's just the content.
	// Given we don't know the exact schema, we will return the stdout as string unless it's obviously a wrapper.
	// But to be safe, we return string.

	output := strings.TrimSpace(stdout.String())
	// Try to unmarshal to see if there's a "content" field or similar,
	// or if the whole thing is the response.
	// Since we are replacing the internal API which returned text,
	// ensuring we return the actual text content is important.

	// If the output is a JSON string literal (e.g. "Hello"), Unquote it.
	if strings.HasPrefix(output, "\"") && strings.HasSuffix(output, "\"") {
		var text string
		if err := json.Unmarshal([]byte(output), &text); err == nil {
			return text, nil
		}
	}

	return output, nil
}

func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	// For streaming, we might need to check if the CLI supports streaming output to stdout.
	// CLI help mentions "--output-format stream-json".
	// Let's use that!

	fullPrompt, err := p.buildPrompt(prompt, attachments)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	args := []string{"--output-format", "stream-json", "--prompt", fullPrompt}
	if p.model != "" {
		args = append(args, "--model", p.model)
	}

	cmd := exec.CommandContext(ctx, p.binaryPath, args...)
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Parse stdout as stream
	// We expect a stream of JSON chunks.
	// We'll read line by line or decode JSON objects as they come.
	decoder := json.NewDecoder(stdout)
	for {
		var chunk struct {
			Content string `json:"content"`
			// Add other fields if known (e.g. error, done)
		}
		if err := decoder.Decode(&chunk); err != nil {
			break // EOF or error
		}
		if chunk.Content != "" {
			tokens <- chunk.Content
		}
	}

	return cmd.Wait()
}

func (p *Provider) buildPrompt(prompt string, attachments []adapter.Attachment) (string, error) {
	var sb strings.Builder

	// Append attachments first
	for _, att := range attachments {
		sb.WriteString(fmt.Sprintf("\n--- File: %s ---\n", att.Name))
		sb.WriteString(att.Data)
		sb.WriteString("\n--- End File ---\n")
	}

	sb.WriteString("\n")
	sb.WriteString(prompt)

	return sb.String(), nil
}
