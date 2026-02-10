package gemini_cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Code Map:
// - GeminiCliToken: Token structure for Gemini CLI
// - SaveTokenToFile: Persists token to disk
// - LoadTokenFromFile: Loads token from disk

const (
	REDIRECT_URI = "http://localhost:8085/oauth2callback"
	AUTH_URL     = "https://accounts.google.com/o/oauth2/v2/auth"
	TOKEN_URL    = "https://oauth2.googleapis.com/token"

	DEFAULT_CLIENT_ID     = "681255809395-oo8ft2oprdrnp9e3aqf6av3hmdib135j.apps.googleusercontent.com"
	DEFAULT_CLIENT_SECRET = "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"
)

var SCOPES = []string{
	"openid",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/userinfo.profile",
}

// GeminiCliToken structure compatible with auth.Token but with specific JSON tags
// CID: GEMINI_TOKEN_STRUCT
type GeminiCliToken struct {
	Access    string    `json:"access_token"`
	Refresh   string    `json:"refresh_token"`
	ExpiresIn int       `json:"expires_in"`
	Expiry    time.Time `json:"expiry"`
	ProjectID string    `json:"project_id"`
}

func (p *Provider) getTokenFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".gemini_cli_token.json")
}

// SaveTokenToFile persists the token to disk
// CID: GEMINI_TOKEN_SAVE
func (p *Provider) SaveTokenToFile(token *GeminiCliToken) error {
	path := p.getTokenFilePath()
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTokenFromFile reads the token from disk
// CID: GEMINI_TOKEN_LOAD
func (p *Provider) LoadTokenFromFile() (*GeminiCliToken, error) {
	path := p.getTokenFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token GeminiCliToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
