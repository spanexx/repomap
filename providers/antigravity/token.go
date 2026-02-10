package antigravity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Code Map:
// - AntigravityToken: Token structure for Antigravity
// - SaveTokenToFile: Persists token to disk
// - LoadTokenFromFile: Loads token from disk

const (
	REDIRECT_URI = "http://localhost:51121/oauth-callback"
	AUTH_URL     = "https://accounts.google.com/o/oauth2/v2/auth"
	TOKEN_URL    = "https://oauth2.googleapis.com/token"

	// Antigravity OAuth credentials (different from Gemini CLI)
	DEFAULT_CLIENT_ID     = "1071006060591-tmhssin2h21lcre235vtolojh4g403ep.apps.googleusercontent.com"
	DEFAULT_CLIENT_SECRET = "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"

	// Fallback project ID when discovery fails
	DEFAULT_PROJECT_ID = "rising-fact-p41fc"
)

// Antigravity requires additional scopes
var SCOPES = []string{
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
	"https://www.googleapis.com/auth/cclog",
	"https://www.googleapis.com/auth/experimentsandconfigs",
}

// AntigravityToken represents the OAuth token and metadata
// CID: ANTIGRAVITY_TOKEN_STRUCT
type AntigravityToken struct {
	Access    string    `json:"access_token"`
	Refresh   string    `json:"refresh_token"`
	ExpiresIn int       `json:"expires_in"`
	Expiry    time.Time `json:"expiry"`
	ProjectID string    `json:"project_id"`
	Email     string    `json:"email,omitempty"`
}

func (p *Provider) getTokenFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".antigravity_token.json")
}

// SaveTokenToFile persists the token to disk
// CID: ANTIGRAVITY_TOKEN_SAVE
func (p *Provider) SaveTokenToFile(token *AntigravityToken) error {
	path := p.getTokenFilePath()
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTokenFromFile reads the token from disk
// CID: ANTIGRAVITY_TOKEN_LOAD
func (p *Provider) LoadTokenFromFile() (*AntigravityToken, error) {
	path := p.getTokenFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token AntigravityToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
