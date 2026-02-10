package gemini_cli

import (
	"context"
	"fmt"
)

const (
	// Cloud Code Assist API endpoint (same as OpenClaw/pi-ai SDK)
	DEFAULT_ENDPOINT = "https://cloudcode-pa.googleapis.com"
	DEFAULT_MODEL    = "gemini-1.5-flash-latest"
	MAX_RETRIES      = 3
	BASE_DELAY_MS    = 2000
)

// Headers for Gemini CLI (same as OpenClaw/pi-ai SDK)
var GEMINI_CLI_HEADERS = map[string]string{
	"User-Agent":        "google-cloud-sdk vscode_cloudshelleditor/0.1",
	"X-Goog-Api-Client": "gl-node/22.17.0",
	"Client-Metadata":   `{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}`,
}

type Provider struct {
	accessToken  string
	projectID    string
	model        string
	systemPrompt string
}

func New(ctx context.Context) (*Provider, error) {
	return &Provider{
		model: DEFAULT_MODEL,
	}, nil
}

func (p *Provider) SetAccessToken(token string) {
	p.accessToken = token
}

func (p *Provider) SetProjectID(projectID string) {
	p.projectID = projectID
}

func (p *Provider) Init(ctx context.Context) error {
	if p.accessToken == "" {
		return fmt.Errorf("access token is required")
	}
	if p.projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	return nil
}

func (p *Provider) Name() string {
	return "gemini-cli"
}

// SetModel sets the model to use for this provider.
func (p *Provider) SetModel(model string) {
	if model != "" {
		p.model = model
	}
}

func (p *Provider) SetSystemPrompt(prompt string) {
	p.systemPrompt = prompt
}
