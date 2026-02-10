package copilot

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Code Map:
// - CopilotToken: Token structure for Copilot
// - SaveTokenToFile: Persists token to disk
// - LoadTokenFromFile: Loads token from disk

const (
	TOKEN_FILE = ".copilot_token.json"
)

// CopilotToken represents the authentication token for GitHub Copilot.
// CID: COPILOT_TOKEN_STRUCT
type CopilotToken struct {
	Access           string `json:"access"`
	Refresh          string `json:"refresh"` // Stores GitHub token
	Expires          int64  `json:"expires"` // Milliseconds
	EnterpriseDomain string `json:"enterprise_domain,omitempty"`
}

func (p *Provider) tokenPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, TOKEN_FILE)
}

// SaveTokenToFile persists the token to disk
// CID: COPILOT_TOKEN_SAVE
func (p *Provider) SaveTokenToFile(token *CopilotToken) error {
	path := p.tokenPath()
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTokenFromFile reads the token from disk
// CID: COPILOT_TOKEN_LOAD
func (p *Provider) LoadTokenFromFile() (*CopilotToken, error) {
	path := p.tokenPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token CopilotToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
