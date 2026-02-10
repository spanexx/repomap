package qwen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Code Map:
// - QwenToken: Token structure for Qwen
// - SaveTokenToFile: Persists token to disk
// - LoadTokenFromFile: Loads token from disk

// QwenToken represents the OAuth token and metadata
// CID: QWEN_TOKEN_STRUCT
type QwenToken struct {
	Access      string    `json:"access_token"`
	Refresh     string    `json:"refresh_token"`
	Expires     int64     `json:"expires_in"`
	Expiry      time.Time `json:"expiry"`
	ResourceURL string    `json:"resourceUrl"`
}

func (p *Provider) getTokenFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".qwen_token.json")
}

// SaveTokenToFile persists the token to disk
// CID: QWEN_TOKEN_SAVE
func (p *Provider) SaveTokenToFile(token *QwenToken) error {
	path := p.getTokenFilePath()
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadTokenFromFile reads the token from disk
// CID: QWEN_TOKEN_LOAD
func (p *Provider) LoadTokenFromFile() (*QwenToken, error) {
	path := p.getTokenFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var token QwenToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return &token, nil
}
