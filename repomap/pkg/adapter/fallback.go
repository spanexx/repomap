package adapter

import (
	"fmt"
	"strings"
)

type FallbackProvider struct {
	providers      []Provider
	verbose        bool
	activeProvider Provider // tracks which provider actually responded
}

func NewFallbackProvider(providers []Provider, verbose bool) *FallbackProvider {
	clean := make([]Provider, 0, len(providers))
	for _, p := range providers {
		if p == nil {
			continue
		}
		clean = append(clean, p)
	}
	return &FallbackProvider{providers: clean, verbose: verbose}
}

func (f *FallbackProvider) Generate(prompt string, attachments []Attachment) (string, error) {
	if len(f.providers) == 0 {
		return "", fmt.Errorf("no providers configured")
	}

	var errs []string
	for i, p := range f.providers {
		resp, err := p.Generate(prompt, attachments)
		if err == nil {
			f.activeProvider = p
			return resp, nil
		}
		errs = append(errs, fmt.Sprintf("%s: %v", p.Name(), err))
		if f.verbose && i < len(f.providers)-1 {
			fmt.Printf("WARN: provider %s failed, trying next: %v\n", p.Name(), err)
		}
	}
	return "", fmt.Errorf("all providers failed: %s", strings.Join(errs, "; "))
}

func (f *FallbackProvider) Name() string {
	if len(f.providers) == 0 {
		return "fallback"
	}
	parts := make([]string, 0, len(f.providers))
	for _, p := range f.providers {
		parts = append(parts, p.Name())
	}
	return "Fallback(" + strings.Join(parts, ",") + ")"
}

// ActiveProviderName returns the name of the provider that actually responded
func (f *FallbackProvider) ActiveProviderName() string {
	if f.activeProvider != nil {
		return f.activeProvider.Name()
	}
	return f.Name()
}

func (f *FallbackProvider) SetModel(model string) {
	if model == "" {
		return // Don't wipe out specific lineup overrides if global model is unset
	}
	for _, p := range f.providers {
		p.SetModel(model)
	}
}

func (f *FallbackProvider) SetSystemPrompt(prompt string) {
	for _, p := range f.providers {
		p.SetSystemPrompt(prompt)
	}
}

// GenerateStream implements streaming for FallbackProvider.
func (f *FallbackProvider) GenerateStream(prompt string, attachments []Attachment, tokens chan<- string) error {
	if len(f.providers) == 0 {
		return fmt.Errorf("no providers configured")
	}

	var errs []string
	for i, p := range f.providers {
		// Try to stream with this provider
		err := p.GenerateStream(prompt, attachments, tokens)
		if err == nil {
			f.activeProvider = p
			return nil
		}

		errMsg := fmt.Sprintf("provider %s failed: %v", p.Name(), err)
		errs = append(errs, errMsg)

		if f.verbose && i < len(f.providers)-1 {
			// If streaming fails mid-stream, we might have already sent partial tokens.
			// Ideally we verify if ANY tokens were sent before retrying, or we restart stream.
			// Current interface assumes atomic failure or full success for simplicity, or client handles partial.
			// Real fallback for streaming is tricky if it fails mid-stream.
			// For MVP, we catch error returned by GenerateStream.
			fmt.Printf("WARN: streaming with %s failed, trying next: %v\n", p.Name(), err)
		}
	}

	return fmt.Errorf("all streaming providers failed: %s", strings.Join(errs, "; "))
}
