package auth

import (
	"context"
)

// OAuth2Provider defines the interface for OAuth2 operations
type OAuth2Provider interface {
	GetAuthURL(provider, state string) (string, error)
	HandleCallback(ctx context.Context, provider, code string) (*OAuthUserInfo, error)
}

// For testing purposes only
var (
	testProvider OAuth2Provider
)

// SetTestProvider allows tests to inject a mock provider
func SetTestProvider(provider OAuth2Provider) {
	testProvider = provider
}

// GetOAuth2Provider returns either the test provider during testing or the real provider
func GetOAuth2Provider() OAuth2Provider {
	if testProvider != nil {
		return testProvider
	}
	return GetOAuth2Manager()
}
