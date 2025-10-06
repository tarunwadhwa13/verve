package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOAuth2Mock(t *testing.T) {
	// Initialize test config
	config := &OAuth2Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{"email", "profile"},
		AuthURL:      "https://example.com/auth",
		TokenURL:     "https://example.com/token",
	}
	InitializeTestOAuth2Config(config)

	// Test the mock provider
	provider := GetOAuth2Provider()
	assert.NotNil(t, provider, "Provider should not be nil")

	// Test GetAuthURL
	url, err := provider.GetAuthURL("google", "test-state")
	assert.NoError(t, err)
	assert.Contains(t, url, "test-client-id")

	// Test HandleCallback
	userInfo, err := provider.HandleCallback(context.Background(), "google", "test-code")
	assert.NoError(t, err)
	assert.NotNil(t, userInfo)
	assert.Equal(t, "test@example.com", userInfo.Email)
	assert.Equal(t, "google", userInfo.Provider)
}
