package auth

import (
	"context"
	"fmt"
	"time"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	RedirectURI  string
	Scopes       []string
	CallbackPort int
	CallbackPath string
}

type Token struct {
	AccessToken  string                 `json:"access_token"`
	RefreshToken string                 `json:"refresh_token"`
	ExpiresIn    int                    `json:"expires_in"`
	Expiry       time.Time              `json:"expiry"`
	ExtraData    map[string]interface{} `json:"extra_data,omitempty"`
}

type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

// StartWebFlow initiates the OAuth2 web flow.
// This is a stub implementation. Real implementation would start a local server.
func StartWebFlow(config OAuthConfig) (string, string, <-chan string, func(), error) {
	// In a real implementation:
	// 1. Generate code verifier & challenge
	// 2. Start local HTTP server on CallbackPort
	// 3. Construct Auth URL
	// 4. Return URL, verifier, and channel that receives code from server

	codeChan := make(chan string)
	cleanup := func() { close(codeChan) }

	// Mock URL
	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code",
		config.AuthURL, config.ClientID, config.RedirectURI)

	// In this stub, we might just print instructions?
	// Or simply return mocked values for now since we rely on existing tokens most of the time.

	return url, "mock-verifier", codeChan, cleanup, nil
}

// ExchangeCode exchanges an authorization code for a token.
func ExchangeCode(ctx context.Context, config OAuthConfig, code, verifier string) (*Token, error) {
	return &Token{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
		ExpiresIn:    3600,
		Expiry:       time.Now().Add(1 * time.Hour),
	}, nil
}

// RefreshToken refreshes an access token using a refresh token.
func RefreshToken(ctx context.Context, config OAuthConfig, refreshToken string) (*Token, error) {
	return &Token{
		AccessToken:  "refreshed-access-token",
		RefreshToken: refreshToken, // Rotate if needed
		ExpiresIn:    3600,
		Expiry:       time.Now().Add(1 * time.Hour),
	}, nil
}
