package copilot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp" // Added
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/auth"
)

// Code Map:
// - Login: Main entry to get valid token
// - StartAuth: Initiates GitHub device flow
// - PollAuth: Polls GitHub then exchanges for Copilot token
// - RefreshToken: Refreshes Copilot token using GitHub token

const (
	// Client ID for Copilot (base64 encoded: Iv1.b507a08c87ecfe98)
	CLIENT_ID_ENCODED = "SXYxLmI1MDdhMDhjODdlY2ZlOTg="
	DEFAULT_SCOPE     = "read:user"
)

// Copilot-specific headers
var COPILOT_HEADERS = map[string]string{
	"User-Agent":             "GitHubCopilotChat/0.35.0",
	"Editor-Version":         "vscode/1.107.0",
	"Editor-Plugin-Version":  "copilot-chat/0.35.0",
	"Copilot-Integration-Id": "vscode-chat",
}

// Internal response structs
type deviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	Interval        int    `json:"interval"`
	ExpiresIn       int    `json:"expires_in"`
}

type gitHubTokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Error       string `json:"error,omitempty"`
	Interval    int    `json:"interval,omitempty"`
}

type copilotTokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func decodeClientID() string {
	decoded, _ := base64.StdEncoding.DecodeString(CLIENT_ID_ENCODED)
	return string(decoded)
}

func getURLs(domain string) (deviceCodeURL, accessTokenURL, copilotTokenURL string) {
	return fmt.Sprintf("https://%s/login/device/code", domain),
		fmt.Sprintf("https://%s/login/oauth/access_token", domain),
		fmt.Sprintf("https://api.%s/copilot_internal/v2/token", domain)
}

// GetBaseURLFromToken extracts the API base URL from a Copilot token
// Token format: tid=...;exp=...;proxy-ep=proxy.individual.githubcopilot.com;...
func GetBaseURLFromToken(token string) string {
	re := regexp.MustCompile(`proxy-ep=([^;]+)`)
	match := re.FindStringSubmatch(token)
	if len(match) < 2 {
		return "https://api.individual.githubcopilot.com"
	}
	// Convert proxy.xxx to api.xxx
	apiHost := strings.Replace(match[1], "proxy.", "api.", 1)
	return "https://" + apiHost
}

// Login attempts to load a valid token. If interaction is required, it returns error.
func (p *Provider) Login(ctx context.Context) (*auth.Token, error) {
	// 1. Try Load
	token, err := p.LoadTokenFromFile()
	if err == nil {
		if token.Expires > time.Now().UnixMilli() {
			return p.convertToAuthToken(token), nil
		}
		// 2. Refresh
		if token.Refresh != "" {
			fmt.Println("🔄 Refreshing Copilot token...")
			newToken, err := p.RefreshToken(ctx, token.Refresh, token.EnterpriseDomain) // Pass domain for refresh
			if err == nil {
				cToken := p.convertFromAuthToken(newToken)
				p.SaveTokenToFile(cToken)
				return newToken, nil
			}
			fmt.Printf("⚠️ Refresh failed: %v.\n", err)
		}
	}

	return nil, fmt.Errorf("login required")
}

// StartAuth initiates device code flow and returns info for UI display.
// Returns (deviceCodeResp, domain, error). Domain is passed as verifier string.
func (p *Provider) StartAuth(ctx context.Context) (*auth.DeviceCode, string, error) {
	domain := "github.com"
	deviceCodeURL, _, _ := getURLs(domain)

	deviceResp, err := p.startDeviceFlow(ctx, deviceCodeURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to start device flow: %w", err)
	}

	return &auth.DeviceCode{
		DeviceCode:      deviceResp.DeviceCode,
		UserCode:        deviceResp.UserCode,
		VerificationURI: deviceResp.VerificationURI,
		ExpiresIn:       deviceResp.ExpiresIn,
		Interval:        deviceResp.Interval,
	}, domain, nil
}

// PollAuth polls for auth completion and exchanges for Copilot token.
// deviceCode is from StartAuth, verifier is valid domain (default github.com).
func (p *Provider) PollAuth(ctx context.Context, deviceCode, verifier string) (*auth.Token, error) {
	domain := verifier
	if domain == "" {
		domain = "github.com"
	}
	_, accessTokenURL, copilotTokenURL := getURLs(domain) // Ignored deviceCodeURL

	// Poll for GitHub token
	githubToken, err := p.pollForGitHubToken(ctx, accessTokenURL, deviceCode, 5, 900) // 15 min default timeout
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub token: %w", err)
	}

	// Exchange for Copilot token
	authToken, err := p.getCopilotToken(ctx, copilotTokenURL, githubToken, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to get Copilot token: %w", err)
	}

	// Save
	cToken := p.convertFromAuthToken(authToken)
	p.SaveTokenToFile(cToken)

	return authToken, nil
}

// RefreshToken implementation wrapper
func (p *Provider) RefreshToken(ctx context.Context, githubToken, domain string) (*auth.Token, error) {
	if domain == "" {
		domain = "github.com"
	}
	_, _, copilotTokenURL := getURLs(domain)

	return p.getCopilotToken(ctx, copilotTokenURL, githubToken, domain)
}

// Internal helpers

func (p *Provider) startDeviceFlow(ctx context.Context, url string) (*deviceCodeResponse, error) {
	reqBody := fmt.Sprintf(`{"client_id":"%s","scope":"%s"}`, decodeClientID(), DEFAULT_SCOPE)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", COPILOT_HEADERS["User-Agent"])

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("device code request failed (%d): %s", resp.StatusCode, string(body))
	}

	var deviceResp deviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return nil, err
	}

	return &deviceResp, nil
}

func (p *Provider) pollForGitHubToken(ctx context.Context, url, deviceCode string, interval, expiresIn int) (string, error) {
	if interval < 5 {
		interval = 5
	}
	deadline := time.Now().Add(time.Duration(expiresIn) * time.Second)

	timer := time.NewTimer(0)
	defer timer.Stop()

	for time.Now().Before(deadline) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-timer.C:
		}

		reqBody := fmt.Sprintf(`{"client_id":"%s","device_code":"%s","grant_type":"urn:ietf:params:oauth:grant-type:device_code"}`,
			decodeClientID(), deviceCode)

		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(reqBody))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", COPILOT_HEADERS["User-Agent"])

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			timer.Reset(time.Duration(interval) * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var tokenResp gitHubTokenResponse
		if err := json.Unmarshal(body, &tokenResp); err != nil {
			timer.Reset(time.Duration(interval) * time.Second)
			continue
		}

		if tokenResp.AccessToken != "" {
			return tokenResp.AccessToken, nil
		}

		switch tokenResp.Error {
		case "authorization_pending":
			timer.Reset(time.Duration(interval) * time.Second)
			continue
		case "slow_down":
			if tokenResp.Interval > 0 {
				interval = tokenResp.Interval
			} else {
				interval += 5
			}
			timer.Reset(time.Duration(interval) * time.Second)
			continue
		default:
			return "", fmt.Errorf("GitHub OAuth error: %s", tokenResp.Error)
		}
	}

	return "", fmt.Errorf("device flow timed out")
}

func (p *Provider) getCopilotToken(ctx context.Context, url, githubToken, domain string) (*auth.Token, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+githubToken)
	for k, v := range COPILOT_HEADERS {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Copilot token request failed (%d): %s", resp.StatusCode, string(body))
	}

	var copilotResp copilotTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&copilotResp); err != nil {
		return nil, err
	}

	// Calculate expiry time (Copilot token expiry is seconds from API)
	expiry := time.Unix(copilotResp.ExpiresAt, 0)

	return &auth.Token{
		AccessToken:  copilotResp.Token,
		RefreshToken: githubToken,
		Expiry:       expiry,
		ExtraData:    map[string]interface{}{"enterprise_domain": domain},
	}, nil
}

func (p *Provider) convertToAuthToken(t *CopilotToken) *auth.Token {
	return &auth.Token{
		AccessToken:  t.Access,
		RefreshToken: t.Refresh,
		Expiry:       time.Unix(t.Expires/1000, 0), // Convert ms to time
		ExtraData:    map[string]interface{}{"enterprise_domain": t.EnterpriseDomain},
	}
}

func (p *Provider) convertFromAuthToken(t *auth.Token) *CopilotToken {
	domain, _ := t.ExtraData["enterprise_domain"].(string)
	return &CopilotToken{
		Access:           t.AccessToken,
		Refresh:          t.RefreshToken,
		Expires:          t.Expiry.UnixMilli(),
		EnterpriseDomain: domain,
	}
}
