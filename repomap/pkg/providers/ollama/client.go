package ollama

import (
	"os"

	"github.com/spanexx/agents-cli/repomap/pkg/adapter"
	"github.com/spanexx/agents-cli/repomap/pkg/providers/generic"
)

// Provider wraps the generic provider with Ollama-specific defaults.
type Provider struct {
	*generic.Provider
}

// New creates a new Ollama provider with the given base URL and model.
// If baseURL is empty, it defaults to OLLAMA_BASE_URL env var or localhost:11434.
// If model is empty, it defaults to OLLAMA_MODEL env var or "llama3.2".
func New(baseURL, model string) *Provider {
	if baseURL == "" {
		baseURL = os.Getenv("OLLAMA_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
	}
	if model == "" {
		model = os.Getenv("OLLAMA_MODEL")
		if model == "" {
			model = "llama3.2"
		}
	}
	return &Provider{
		Provider: generic.New("ollama", baseURL, "", model),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "ollama"
}

// GenerateStream wraps the generic provider's GenerateStream
// GenerateStream implements streaming response (simulated).
func (p *Provider) GenerateStream(prompt string, attachments []adapter.Attachment, tokens chan<- string) error {
	resp, err := p.Generate(prompt, attachments)
	if err != nil {
		return err
	}
	tokens <- resp
	return nil
}

func (p *Provider) Generate(prompt string, attachments []adapter.Attachment) (string, error) {
	return p.Provider.Generate(prompt, attachments)
}

func (p *Provider) SetSystemPrompt(prompt string) {
	p.Provider.SetSystemPrompt(prompt)
}
