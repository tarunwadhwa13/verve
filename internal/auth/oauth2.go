package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuth2Manager handles OAuth2 configuration and token management
type OAuth2Manager struct {
	configs map[string]*oauth2.Config
	tokens  sync.Map // thread-safe map for storing tokens
}

// OAuth2Config holds the configuration for supported OAuth providers
type OAuth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	Scopes       []string `json:"scopes"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
}

// OAuthUserInfo is defined in oauth.go

var (
	manager *OAuth2Manager
	once    sync.Once
)

// GetOAuth2Manager returns a singleton instance of OAuth2Manager
func GetOAuth2Manager() *OAuth2Manager {
	once.Do(func() {
		manager = &OAuth2Manager{
			configs: make(map[string]*oauth2.Config),
		}
		manager.initializeConfigs()
	})
	return manager
}

// GetAuthURL returns the authorization URL for the specified provider
func (m *OAuth2Manager) GetAuthURL(provider, state string) (string, error) {
	config, ok := m.configs[provider]
	if !ok {
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
	return config.AuthCodeURL(state), nil
}

// HandleCallback processes the OAuth callback and returns user information
func (m *OAuth2Manager) HandleCallback(ctx context.Context, provider, code string) (*OAuthUserInfo, error) {
	config, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	return m.getUserInfo(ctx, provider, token)
}

// getUserInfo fetches user information from the OAuth provider
func (m *OAuth2Manager) getUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))

	switch provider {
	case "google":
		return getGoogleUserInfo(ctx, client)
	case "okta":
		return getOktaUserInfo(ctx, client)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// initializeConfigs sets up OAuth configurations for supported providers
func (m *OAuth2Manager) initializeConfigs() {
	// Google OAuth2 Config
	m.configs["google"] = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Okta OAuth2 Config
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	if oktaDomain == "" {
		oktaDomain = "dev-123456.okta.com" // Default domain for development
	}
	m.configs["okta"] = &oauth2.Config{
		ClientID:     os.Getenv("OKTA_CLIENT_ID"),
		ClientSecret: os.Getenv("OKTA_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OKTA_REDIRECT_URL"),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   fmt.Sprintf("https://%s/oauth2/v1/authorize", oktaDomain),
			TokenURL:  fmt.Sprintf("https://%s/oauth2/v1/token", oktaDomain),
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

// getAuthorizationURL is an internal method to get the authorization URL with offline access
func (m *OAuth2Manager) getAuthorizationURL(provider, state string) (string, error) {
	config, ok := m.configs[provider]
	if !ok {
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// Exchange exchanges an authorization code for OAuth2 tokens
func (m *OAuth2Manager) Exchange(ctx context.Context, provider, code string) (*oauth2.Token, error) {
	config, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	return config.Exchange(ctx, code)
}

// GetUserInfo fetches user information from the OAuth provider
func (m *OAuth2Manager) GetUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	config, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	client := config.Client(ctx, token)
	userInfo := &OAuthUserInfo{
		Provider:      provider,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		TokenExpiry:   token.Expiry,
		ProviderToken: token,
	}

	var err error
	switch provider {
	case "google":
		err = m.fetchGoogleUserInfo(client, userInfo)
	case "okta":
		err = m.fetchOktaUserInfo(client, userInfo)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if err != nil {
		return nil, err
	}

	// Store token in memory
	m.tokens.Store(userInfo.ID, token)
	return userInfo, nil
}

func (m *OAuth2Manager) fetchGoogleUserInfo(client *http.Client, userInfo *OAuthUserInfo) error {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return fmt.Errorf("failed to get user info from Google: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get user info from Google: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Google user info response: %v", err)
	}

	if err := json.Unmarshal(data, userInfo); err != nil {
		return fmt.Errorf("failed to parse Google user info: %v", err)
	}

	return nil
}

func (m *OAuth2Manager) fetchOktaUserInfo(client *http.Client, userInfo *OAuthUserInfo) error {
	resp, err := client.Get(fmt.Sprintf("https://%s/oauth2/v1/userinfo", os.Getenv("OKTA_DOMAIN")))
	if err != nil {
		return fmt.Errorf("failed to get user info from Okta: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get user info from Okta: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read Okta user info response: %v", err)
	}

	if err := json.Unmarshal(data, userInfo); err != nil {
		return fmt.Errorf("failed to parse Okta user info: %v", err)
	}

	return nil
}

// RefreshToken refreshes an expired OAuth2 token
func (m *OAuth2Manager) RefreshToken(ctx context.Context, provider string, token *oauth2.Token) (*oauth2.Token, error) {
	config, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	if !token.Valid() && token.RefreshToken != "" {
		newToken, err := config.TokenSource(ctx, token).Token()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		return newToken, nil
	}

	return token, nil
}

// ValidateToken validates an OAuth token and refreshes it if necessary
func (m *OAuth2Manager) ValidateToken(ctx context.Context, provider, accessToken string) (*OAuthUserInfo, error) {
	// Create a token for validation
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	config, ok := m.configs[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	client := config.Client(ctx, token)
	userInfo := &OAuthUserInfo{Provider: provider}

	switch provider {
	case "google":
		if err := m.fetchGoogleUserInfo(client, userInfo); err != nil {
			return nil, err
		}
	case "okta":
		if err := m.fetchOktaUserInfo(client, userInfo); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported provider")
	}

	return userInfo, nil
}

// RemoveToken removes a stored token
func (m *OAuth2Manager) RemoveToken(userID string) {
	m.tokens.Delete(userID)
}

// GetStoredToken retrieves a stored token
func (m *OAuth2Manager) GetStoredToken(userID string) *oauth2.Token {
	if token, ok := m.tokens.Load(userID); ok {
		return token.(*oauth2.Token)
	}
	return nil
}
