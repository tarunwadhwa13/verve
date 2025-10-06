package auth

// InitializeTestOAuth2Config sets up OAuth2 config for testing
func InitializeTestOAuth2Config(config *OAuth2Config) {
	mock := NewMockOAuth2Manager()
	mock.setupMockConfig(config)

	// Setup default mock user info
	mock.SetMockUserInfo(&OAuthUserInfo{
		ID:       "12345",
		Email:    "test@example.com",
		Provider: "google",
		Name:     "Test User",
		Picture:  "https://example.com/photo.jpg",
	})

	SetTestProvider(mock)
}

// SetupTestOAuthUserInfo allows customizing the mock user info
func SetupTestOAuthUserInfo(info *OAuthUserInfo) {
	if provider, ok := GetOAuth2Provider().(*MockOAuth2Manager); ok {
		provider.SetMockUserInfo(info)
	}
}
