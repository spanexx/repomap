package qwen

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/auth"
)

// Code Map:
// - Login: Main entry to get valid token
// - StartAuth: Initiates device flow (PKCE)
// - PollAuth: Polls for token (PKCE)
// - RefreshToken: Refreshes access token

const (
	QWEN_OAUTH_BASE_URL             = "https://chat.qwen.ai"
	QWEN_OAUTH_DEVICE_CODE_ENDPOINT = QWEN_OAUTH_BASE_URL + "/api/v1/oauth2/device/code"
	QWEN_OAUTH_TOKEN_ENDPOINT       = QWEN_OAUTH_BASE_URL + "/api/v1/oauth2/token"
	QWEN_OAUTH_CLIENT_ID            = "f0304373b74a44d2b584a3fb70ca9e56"
	QWEN_OAUTH_SCOPE                = "openid profile email model.completion"
	QWEN_OAUTH_GRANT_TYPE           = "urn:ietf:params:oauth:grant-type:device_code"
)

type tokenResponse struct {
	Access    string `json:"access_token"`
	Refresh   string `json:"refresh_token"`
	Expires   int64  `json:"expires_in"`
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

// Login attempts to load a valid token. If interaction is required, it returns error.
// CID: QWEN_AUTH_LOGIN
func (p *Provider) Login(ctx context.Context) (*auth.Token, error) {
	// 1. Try environment variable
	apiKey := os.Getenv("QWEN_ACCESS_TOKEN")
	if apiKey != "" {
		return &auth.Token{
			AccessToken: apiKey,
			Expiry:      time.Now().Add(24 * time.Hour),
		}, nil
	}

	// 2. Try Load
	token, err := p.LoadTokenFromFile()
	if err == nil {
		if time.Now().Before(token.Expiry.Add(-1 * time.Minute)) {
			return p.convertToAuthToken(token), nil
		}
		// 3. Refresh
		fmt.Println("🔄 Token expired, refreshing...")
		newToken, err := p.RefreshToken(ctx, token.Refresh)
		if err == nil {
			qToken := p.convertFromAuthToken(newToken)
			p.SaveTokenToFile(qToken)
			return newToken, nil
		}
		fmt.Printf("⚠️ Refresh failed: %v.\n", err)
	}

	return nil, fmt.Errorf("login required")
}

// StartAuth initiates the Device Code flow returns the authorization data and verifier
// CID: QWEN_AUTH_START
func (p *Provider) StartAuth(ctx context.Context) (*auth.DeviceCode, string, error) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return nil, "", err
	}

	authData, err := requestDeviceCode(ctx, challenge)
	if err != nil {
		return nil, "", err
	}
	return authData, verifier, nil
}

// PollAuth polls for the token after StartAuth
// CID: QWEN_AUTH_POLL
func (p *Provider) PollAuth(ctx context.Context, deviceCode, verifier string) (*auth.Token, error) {
	token, err := pollDeviceToken(ctx, deviceCode, verifier)
	if err != nil {
		return nil, err
	}

	// Save
	qToken := p.convertFromAuthToken(token)
	p.SaveTokenToFile(qToken)

	return token, nil
}

// RefreshToken exchanges a refresh token for a new access token
// CID: QWEN_AUTH_REFRESH
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*auth.Token, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", QWEN_OAUTH_CLIENT_ID)
	data.Set("refresh_token", refreshToken)

	req, _ := http.NewRequestWithContext(ctx, "POST", QWEN_OAUTH_TOKEN_ENDPOINT, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh failed: %d %s", resp.StatusCode, string(body))
	}

	var tr tokenResponse
	if err := json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}

	return &auth.Token{
		AccessToken:  tr.Access,
		RefreshToken: tr.Refresh,
		ExpiresIn:    int(tr.Expires),
		Expiry:       time.Now().Add(time.Duration(tr.Expires) * time.Second),
		ExtraData:    map[string]interface{}{"resourceUrl": "https://portal.qwen.ai/v1"},
	}, nil
}

// Helpers

func generatePKCE() (string, string, error) {
	verifier := make([]byte, 32)
	_, err := rand.Read(verifier)
	if err != nil {
		return "", "", err
	}
	verifierStr := base64.RawURLEncoding.EncodeToString(verifier)

	hash := sha256.Sum256([]byte(verifierStr))
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])

	return verifierStr, challenge, nil
}

func requestDeviceCode(ctx context.Context, challenge string) (*auth.DeviceCode, error) {
	data := url.Values{}
	data.Set("client_id", QWEN_OAUTH_CLIENT_ID)
	data.Set("scope", QWEN_OAUTH_SCOPE)
	data.Set("code_challenge", challenge)
	data.Set("code_challenge_method", "S256")

	req, _ := http.NewRequestWithContext(ctx, "POST", QWEN_OAUTH_DEVICE_CODE_ENDPOINT, strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to request device code: %d %s", resp.StatusCode, string(body))
	}

	var ad auth.DeviceCode
	if err := json.Unmarshal(body, &ad); err != nil {
		return nil, err
	}

	return &ad, nil
}

func pollDeviceToken(ctx context.Context, deviceCode, verifier string) (*auth.Token, error) {
	data := url.Values{}
	data.Set("grant_type", QWEN_OAUTH_GRANT_TYPE)
	data.Set("client_id", QWEN_OAUTH_CLIENT_ID)
	data.Set("device_code", deviceCode)
	data.Set("code_verifier", verifier)

	maxAttempts := 120
	interval := 2

	// Loop with context awareness
	timer := time.NewTimer(0)
	defer timer.Stop()

	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}

		req, _ := http.NewRequestWithContext(ctx, "POST", QWEN_OAUTH_TOKEN_ENDPOINT, strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var token tokenResponse
			if err := json.Unmarshal(body, &token); err != nil {
				return nil, err
			}

			if token.Error != "" {
				if token.Error == "authorization_pending" {
					timer.Reset(time.Duration(interval) * time.Second)
					continue
				}
				return nil, fmt.Errorf("oauth error: %s - %s", token.Error, token.ErrorDesc)
			}

			return &auth.Token{
				AccessToken:  token.Access,
				RefreshToken: token.Refresh,
				ExpiresIn:    int(token.Expires),
				Expiry:       time.Now().Add(time.Duration(token.Expires) * time.Second),
				ExtraData:    map[string]interface{}{"resourceUrl": "https://portal.qwen.ai/v1"},
			}, nil
		}

		timer.Reset(time.Duration(interval) * time.Second)
	}

	return nil, fmt.Errorf("device authorization timeout")
}

func (p *Provider) convertToAuthToken(t *QwenToken) *auth.Token {
	return &auth.Token{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
		ExpiresIn:    int(t.Expires),
		Expiry:       t.Expiry,
		ExtraData:    map[string]interface{}{"resourceUrl": t.ResourceURL},
	}
}

func (p *Provider) convertFromAuthToken(t *auth.Token) *QwenToken {
	resURL, _ := t.ExtraData["resourceUrl"].(string)
	return &QwenToken{
		Access:      t.AccessToken,
		Refresh:     t.RefreshToken,
		Expires:     int64(t.ExpiresIn),
		Expiry:      t.Expiry,
		ResourceURL: resURL,
	}
}
