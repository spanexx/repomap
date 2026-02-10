package antigravity

import (
	"context"
	"fmt"
	"os"
)

const (
	// Cloud Code Assist API endpoints - Antigravity uses sandbox with fallback
	DEFAULT_ENDPOINT        = "https://cloudcode-pa.googleapis.com"
	SANDBOX_ENDPOINT        = "https://daily-cloudcode-pa.sandbox.googleapis.com"
	DEFAULT_ANTIGRAVITY_VER = "1.15.8"
	DEFAULT_MODEL           = "gemini-1.5-flash-latest"
	MAX_RETRIES             = 3
	BASE_DELAY_MS           = 2000
)

// getAntigravityHeaders returns headers for Antigravity API calls
func getAntigravityHeaders() map[string]string {
	version := os.Getenv("PI_AI_ANTIGRAVITY_VERSION")
	if version == "" {
		version = DEFAULT_ANTIGRAVITY_VER
	}
	return map[string]string{
		"User-Agent":        fmt.Sprintf("antigravity/%s darwin/arm64", version),
		"X-Goog-Api-Client": "google-cloud-sdk vscode_cloudshelleditor/0.1",
		"Client-Metadata":   `{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}`,
	}
}

// Provider is the Antigravity LLM provider
type Provider struct {
	accessToken  string
	projectID    string
	model        string
	endpoints    []string // Fallback endpoints
	systemPrompt string
}

// New creates a new Antigravity provider
func New(ctx context.Context) (*Provider, error) {
	return &Provider{
		model: DEFAULT_MODEL,
		endpoints: []string{
			SANDBOX_ENDPOINT, // Try sandbox first for Antigravity
			DEFAULT_ENDPOINT,
		},
	}, nil
}

// SetAccessToken sets the OAuth access token
func (p *Provider) SetAccessToken(token string) {
	p.accessToken = token
}

// SetProjectID sets the Google Cloud project ID
func (p *Provider) SetProjectID(projectID string) {
	p.projectID = projectID
}

// Init validates the provider configuration
func (p *Provider) Init(ctx context.Context) error {
	if p.accessToken == "" {
		return fmt.Errorf("access token is required")
	}
	if p.projectID == "" {
		return fmt.Errorf("project ID is required")
	}
	return nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "antigravity"
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
