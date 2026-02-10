package lmstudio

import (
	"os"

	"llm-adapter/pkg/providers/generic"
)

// Provider wraps the generic provider with LMStudio-specific defaults.
type Provider struct {
	*generic.Provider
}

// New creates a new LMStudio provider with the given base URL and model.
// If baseURL is empty, it defaults to LMSTUDIO_BASE_URL env var or localhost:1234.
// If model is empty, it defaults to LMSTUDIO_MODEL env var or "default".
func New(baseURL, model string) *Provider {
	if baseURL == "" {
		baseURL = os.Getenv("LMSTUDIO_BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:1234/v1"
		}
	}
	if model == "" {
		model = os.Getenv("LMSTUDIO_MODEL")
		if model == "" {
			model = "default"
		}
	}
	return &Provider{
		Provider: generic.New("lmstudio", baseURL, "", model),
	}
}

// Name returns the provider name.
func (p *Provider) Name() string {
	return "lmstudio"
}
