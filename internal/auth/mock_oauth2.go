package auth

import (
	"context"

	"golang.org/x/oauth2"
)

// MockOAuth2Manager is a test-specific implementation of OAuth2Provider
type MockOAuth2Manager struct {
	configs      map[string]*oauth2.Config
	mockUserInfo *OAuthUserInfo
}

func NewMockOAuth2Manager() *MockOAuth2Manager {
	return &MockOAuth2Manager{
		configs: make(map[string]*oauth2.Config),
	}
}

// SetMockUserInfo sets up mock user info for testing
func (m *MockOAuth2Manager) SetMockUserInfo(info *OAuthUserInfo) {
	m.mockUserInfo = info
}

// HandleCallback implements OAuth2Provider interface
func (m *MockOAuth2Manager) HandleCallback(ctx context.Context, provider, code string) (*OAuthUserInfo, error) {
	if m.mockUserInfo != nil {
		info := *m.mockUserInfo  // Return a copy to prevent modification
		info.Provider = provider // Set the provider from the request
		return &info, nil
	}
	return nil, nil
}

// GetAuthURL implements OAuth2Provider interface
func (m *MockOAuth2Manager) GetAuthURL(provider, state string) (string, error) {
	config, ok := m.configs[provider]
	if !ok {
		return "", nil
	}
	return config.AuthCodeURL(state), nil
}

// setupMockConfig configures the mock OAuth2 provider
func (m *MockOAuth2Manager) setupMockConfig(config *OAuth2Config) {
	m.configs["google"] = &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes:       config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
	}
}
