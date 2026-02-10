package gemini_cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"llm-adapter/pkg/auth"
)

// Code Map:
// - StartWebAuth: Initiates web flow for UI
// - ExchangeAndSave: Exchanges code for token and saves (for UI)
// - InteractiveLogin: Main entry point for CLI login flow
// - Login: Silent login (load + refresh)
// - RefreshToken: Refresh access token
// - getOAuthConfig: Helper to build auth config

// StartWebAuth starts the local callback server and returns the auth URL, verifier, channel for code, and cleanup.
// CID: GEMINI_AUTH_WEB_START
func (p *Provider) StartWebAuth() (string, string, <-chan string, func(), error) {
	return auth.StartWebFlow(p.getOAuthConfig())
}

// ExchangeAndSave exchanges the code for a token and saves it.
// CID: GEMINI_AUTH_EXCHANGE
func (p *Provider) ExchangeAndSave(code, verifier string) (*GeminiCliToken, error) {
	config := p.getOAuthConfig()
	// Context is needed for request, using Background as this might be called from async handler
	authToken, err := auth.ExchangeCode(context.Background(), config, code, verifier)
	if err != nil {
		return nil, err
	}

	geminiToken := p.convertToken(authToken)
	if geminiToken.ProjectID == "" {
		pid, err := p.discoverProject(geminiToken.Access)
		if err == nil {
			geminiToken.ProjectID = pid
		}
	}

	if err := p.SaveTokenToFile(geminiToken); err != nil {
		return nil, err
	}

	return geminiToken, nil
}

// - Login: Silent login (load + refresh)
// - RefreshToken: Refresh access token
// - getOAuthConfig: Helper to build auth config
// - discoverProject: Discovers the Google Cloud Project ID

// InteractiveLogin performs the full OAuth flow with user interaction if needed.
// CID: GEMINI_AUTH_INTERACTIVE
func (p *Provider) InteractiveLogin(ctx context.Context) (*GeminiCliToken, error) {
	// 1. Try silent login first
	token, err := p.Login(ctx)
	if err == nil {
		return token, nil
	}

	// 2. Start interactive flow
	config := p.getOAuthConfig()
	authURL, verifier, codeChan, cleanup, err := auth.StartWebFlow(config)
	if err != nil {
		return nil, fmt.Errorf("failed to start auth flow: %w", err)
	}
	defer cleanup()

	fmt.Printf("\n🔐 Gemini CLI OAuth Login\n")
	fmt.Println(strings.Repeat("=", 40))
	fmt.Printf("Please open the following URL in your browser:\n\n%s\n\n", authURL)
	fmt.Println("⏳ Waiting for authorization...")

	var code string
	select {
	case code = <-codeChan:
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("timeout waiting for authorization")
	}

	// 3. Exchange code for token
	authToken, err := auth.ExchangeCode(ctx, config, code, verifier)
	if err != nil {
		return nil, fmt.Errorf("token exchange failed: %w", err)
	}

	// 4. Convert and enhance token
	geminiToken := p.convertToken(authToken)

	// Discover project ID if not present
	if geminiToken.ProjectID == "" {
		pid, err := p.discoverProject(geminiToken.Access)
		if err == nil {
			geminiToken.ProjectID = pid
		} else {
			fmt.Printf("Warning: failed to discover project: %v\n", err)
		}
	}

	// 5. Save token
	if err := p.SaveTokenToFile(geminiToken); err != nil {
		fmt.Printf("Warning: failed to save token: %v\n", err)
	}

	if geminiToken.ProjectID != "" {
		fmt.Printf("✅ Login successful! (Project: %s)\n", geminiToken.ProjectID)
	} else {
		fmt.Println("✅ Login successful! (Project ID unknown)")
	}

	return geminiToken, nil
}

// Login attempts to load a valid token. If expired, it attempts refresh.
// CID: GEMINI_AUTH_LOGIN
func (p *Provider) Login(ctx context.Context) (*GeminiCliToken, error) {
	// 1. Load from file
	token, err := p.LoadTokenFromFile()
	if err != nil {
		return nil, err
	}

	// 2. Check expiry (with 5 min buffer)
	if time.Now().Before(token.Expiry.Add(-5 * time.Minute)) {
		return token, nil
	}

	// 3. Refresh if needed
	fmt.Println("🔄 Token expired, refreshing...")
	newToken, err := p.RefreshToken(ctx, token.Refresh)
	if err != nil {
		return nil, fmt.Errorf("refresh failed: %w", err)
	}

	// Preserve project ID
	if token.ProjectID != "" {
		newToken.ProjectID = token.ProjectID
	} else {
		pid, err := p.discoverProject(newToken.Access)
		if err == nil {
			newToken.ProjectID = pid
		}
	}

	p.SaveTokenToFile(newToken)
	return newToken, nil
}

// RefreshToken exchanges a refresh token for a new access token.
// CID: GEMINI_AUTH_REFRESH
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string) (*GeminiCliToken, error) {
	config := p.getOAuthConfig()
	authToken, err := auth.RefreshToken(ctx, config, refreshToken)
	if err != nil {
		return nil, err
	}
	return p.convertToken(authToken), nil
}

func (p *Provider) convertToken(t *auth.Token) *GeminiCliToken {
	return &GeminiCliToken{
		Access:    t.AccessToken,
		Refresh:   t.RefreshToken,
		ExpiresIn: t.ExpiresIn,
		Expiry:    t.Expiry,
	}
}

func (p *Provider) getOAuthConfig() auth.OAuthConfig {
	clientId, clientSecret := getOAuthCredentials()
	return auth.OAuthConfig{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		AuthURL:      AUTH_URL,
		TokenURL:     TOKEN_URL,
		RedirectURI:  REDIRECT_URI,
		Scopes:       SCOPES,
		CallbackPort: 8085,
		CallbackPath: "/oauth2callback",
	}
}

func getOAuthCredentials() (string, string) {
	clientId := os.Getenv("GEMINI_CLI_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GEMINI_CLI_OAUTH_CLIENT_SECRET")

	if clientId == "" {
		clientId = DEFAULT_CLIENT_ID
	}
	if clientSecret == "" {
		clientSecret = DEFAULT_CLIENT_SECRET
	}
	return clientId, clientSecret
}

// CID: GEMINI_PROJECT_DISCOVERY
func (p *Provider) discoverProject(accessToken string) (string, error) {
	// Try environment first
	if pid := os.Getenv("GOOGLE_CLOUD_PROJECT"); pid != "" {
		return pid, nil
	}
	if pid := os.Getenv("GOOGLE_CLOUD_PROJECT_ID"); pid != "" {
		return pid, nil
	}

	// Try to discover via loadCodeAssist (similar to OpenClaw)
	const endpoint = "https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist"
	reqBody := `{"metadata":{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}}`

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(reqBody))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data struct {
			Project string `json:"cloudaicompanionProject"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err == nil {
			if data.Project != "" {
				return data.Project, nil
			}
		}
	}

	return "", fmt.Errorf("no project discovered")
}
