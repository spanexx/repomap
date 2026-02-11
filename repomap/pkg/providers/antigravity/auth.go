package antigravity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spanexx/agents-cli/repomap/pkg/auth"
)

// Code Map:
// - StartWebAuth: Initiates web flow for UI
// - ExchangeAndSave: Exchanges code for token and saves (for UI)
// - InteractiveLogin: Main entry point for login flow
// - Login: Silent login (load + refresh)
// - RefreshToken: Refresh access token
// - getOAuthConfig: Helper to build auth config

// StartWebAuth starts the local callback server and returns the auth URL, verifier, channel for code, and cleanup.
// CID: ANTIGRAVITY_AUTH_WEB_START
func (p *Provider) StartWebAuth() (string, string, <-chan string, func(), error) {
	config, err := p.getOAuthConfig()
	if err != nil {
		return "", "", nil, func() {}, err
	}
	return auth.StartWebFlow(config)
}

// ExchangeAndSave exchanges the code for a token and saves it.
// CID: ANTIGRAVITY_AUTH_EXCHANGE
func (p *Provider) ExchangeAndSave(code, verifier string) (*AntigravityToken, error) {
	config, err := p.getOAuthConfig()
	if err != nil {
		return nil, err
	}
	authToken, err := auth.ExchangeCode(context.Background(), config, code, verifier)
	if err != nil {
		return nil, err
	}

	antigravityToken := p.convertToken(authToken)

	// Fetch Email
	email, err := p.getUserEmail(antigravityToken.Access)
	if err == nil {
		antigravityToken.Email = email
	}

	// Discover project ID
	if antigravityToken.ProjectID == "" {
		pid, err := p.discoverProject(antigravityToken.Access)
		if err == nil {
			antigravityToken.ProjectID = pid
		} else {
			antigravityToken.ProjectID = DEFAULT_PROJECT_ID
		}
	}

	if err := p.SaveTokenToFile(antigravityToken); err != nil {
		return nil, err
	}

	return antigravityToken, nil
}

// - Login: Silent login (load + refresh)
// - RefreshToken: Refresh access token
// - getOAuthConfig: Helper to build auth config
// - getUserEmail: Fetches user email
// - discoverProject: Discovers the Google Cloud Project ID

// InteractiveLogin performs the full OAuth flow with user interaction if needed.
// CID: ANTIGRAVITY_AUTH_INTERACTIVE
func (p *Provider) InteractiveLogin(ctx context.Context) (*AntigravityToken, error) {
	// 1. Try silent login first
	token, err := p.Login(ctx)
	if err == nil {
		return token, nil
	}

	// 2. Start interactive flow
	config, err := p.getOAuthConfig()
	if err != nil {
		return nil, err
	}
	authURL, verifier, codeChan, cleanup, err := auth.StartWebFlow(config)
	if err != nil {
		return nil, fmt.Errorf("failed to start auth flow: %w", err)
	}
	defer cleanup()

	fmt.Printf("\n🚀 Antigravity OAuth Login\n")
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
	antigravityToken := p.convertToken(authToken)

	// Fetch Email
	email, err := p.getUserEmail(antigravityToken.Access)
	if err == nil {
		antigravityToken.Email = email
	}

	// Discover project ID
	if antigravityToken.ProjectID == "" {
		pid, err := p.discoverProject(antigravityToken.Access)
		if err == nil {
			antigravityToken.ProjectID = pid
		} else {
			fmt.Printf("Warning: failed to discover project: %v\n", err)
			antigravityToken.ProjectID = DEFAULT_PROJECT_ID
		}
	}

	// 5. Save token
	if err := p.SaveTokenToFile(antigravityToken); err != nil {
		fmt.Printf("Warning: failed to save token: %v\n", err)
	}

	fmt.Printf("✅ Login successful! (Project: %s)\n", antigravityToken.ProjectID)
	return antigravityToken, nil
}

// Login attempts to load a valid token. If expired, it attempts refresh.
// CID: ANTIGRAVITY_AUTH_LOGIN
func (p *Provider) Login(ctx context.Context) (*AntigravityToken, error) {
	// 1. Load from file
	token, err := p.LoadTokenFromFile()
	if err != nil {
		return nil, err
	}

	// 2. Check expiry
	if time.Now().Before(token.Expiry.Add(-5 * time.Minute)) {
		return token, nil
	}

	// 3. Refresh if needed
	fmt.Println("🔄 Token expired, refreshing...")
	newToken, err := p.RefreshToken(ctx, token.Refresh, token.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("refresh failed: %w", err)
	}

	// Preserve metadata
	if token.ProjectID != "" {
		newToken.ProjectID = token.ProjectID
	} else {
		pid, err := p.discoverProject(newToken.Access)
		if err == nil {
			newToken.ProjectID = pid
		} else {
			newToken.ProjectID = DEFAULT_PROJECT_ID // Fallback
		}
	}
	if token.Email != "" {
		newToken.Email = token.Email
	} else {
		newToken.Email, _ = p.getUserEmail(newToken.Access)
	}

	p.SaveTokenToFile(newToken)
	return newToken, nil
}

// RefreshToken exchanges a refresh token for a new access token.
// CID: ANTIGRAVITY_AUTH_REFRESH
func (p *Provider) RefreshToken(ctx context.Context, refreshToken string, projectID string) (*AntigravityToken, error) {
	config, err := p.getOAuthConfig()
	if err != nil {
		return nil, err
	}
	authToken, err := auth.RefreshToken(ctx, config, refreshToken)
	if err != nil {
		return nil, err
	}
	t := p.convertToken(authToken)
	if t.ProjectID == "" {
		t.ProjectID = projectID
	}
	return t, nil
}

func (p *Provider) convertToken(t *auth.Token) *AntigravityToken {
	return &AntigravityToken{
		Access:    t.AccessToken,
		Refresh:   t.RefreshToken,
		ExpiresIn: t.ExpiresIn,
		Expiry:    t.Expiry,
	}
}

func (p *Provider) getOAuthConfig() (auth.OAuthConfig, error) {
	clientId, clientSecret, err := getOAuthCredentials()
	if err != nil {
		return auth.OAuthConfig{}, err
	}
	return auth.OAuthConfig{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		AuthURL:      AUTH_URL,
		TokenURL:     TOKEN_URL,
		RedirectURI:  REDIRECT_URI,
		Scopes:       SCOPES,
		CallbackPort: 51121,
		CallbackPath: "/oauth-callback",
	}, nil
}

func getOAuthCredentials() (string, string, error) {
	clientId := os.Getenv("ANTIGRAVITY_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("ANTIGRAVITY_OAUTH_CLIENT_SECRET")

	if clientId == "" {
		clientId = DEFAULT_CLIENT_ID
	}
	if clientSecret == "" {
		return "", "", fmt.Errorf("ANTIGRAVITY_OAUTH_CLIENT_SECRET environment variable not set")
	}
	return clientId, clientSecret, nil
}

// getUserEmail logic retained
// CID: ANTIGRAVITY_USER_EMAIL
func (p *Provider) getUserEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data struct {
			Email string `json:"email"`
		}
		json.NewDecoder(resp.Body).Decode(&data)
		return data.Email, nil
	}

	return "", fmt.Errorf("failed to get email")
}

// discoverProject logic retained (with multi-endpoint support)
// CID: ANTIGRAVITY_PROJECT_DISCOVERY
func (p *Provider) discoverProject(accessToken string) (string, error) {
	// Try environment first
	if pid := os.Getenv("GOOGLE_CLOUD_PROJECT"); pid != "" {
		return pid, nil
	}
	if pid := os.Getenv("GOOGLE_CLOUD_PROJECT_ID"); pid != "" {
		return pid, nil
	}

	headers := map[string]string{
		"Authorization":     "Bearer " + accessToken,
		"Content-Type":      "application/json",
		"User-Agent":        "google-api-nodejs-client/9.15.1",
		"X-Goog-Api-Client": "google-cloud-sdk vscode_cloudshelleditor/0.1",
		"Client-Metadata":   `{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}`,
	}

	// Try endpoints in order: prod first, then sandbox
	endpoints := []string{
		"https://cloudcode-pa.googleapis.com",
		"https://daily-cloudcode-pa.sandbox.googleapis.com",
	}

	reqBody := `{"metadata":{"ideType":"IDE_UNSPECIFIED","platform":"PLATFORM_UNSPECIFIED","pluginType":"GEMINI"}}`

	client := &http.Client{Timeout: 5 * time.Second}

	for _, endpoint := range endpoints {
		req, _ := http.NewRequest("POST", endpoint+"/v1internal:loadCodeAssist", strings.NewReader(reqBody))
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		resp, err := client.Do(req)
		if err != nil {
			continue // Try next endpoint
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
	}

	// Use fallback project ID
	return DEFAULT_PROJECT_ID, nil
}
