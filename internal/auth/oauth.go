package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oktaOauthConfig   *oauth2.Config
)

type OAuthUserInfo struct {
	ID            string        `json:"id"`
	Email         string        `json:"email"`
	Name          string        `json:"name"`
	Picture       string        `json:"picture"`
	Provider      string        `json:"provider"`
	AccessToken   string        `json:"-"` // not exposed in JSON
	RefreshToken  string        `json:"-"` // not exposed in JSON
	TokenExpiry   time.Time     `json:"-"` // not exposed in JSON
	ProviderToken *oauth2.Token `json:"-"` // store full OAuth2 token
}

func InitOAuth() {
	oktaDomain := os.Getenv("OKTA_DOMAIN")
	if oktaDomain == "" {
		oktaDomain = "dev-123456.okta.com" // Default domain for development
	}

	// Initialize Google OAuth
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Initialize Okta OAuth with standard OAuth2 endpoints
	oktaOauthConfig = &oauth2.Config{
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

func GetGoogleAuthURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state)
}

func GetOktaAuthURL(state string) string {
	return oktaOauthConfig.AuthCodeURL(state)
}

func GetGoogleUserInfo(code string) (*OAuthUserInfo, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}
	userInfo.Provider = "google"

	return &userInfo, nil
}

func GetOktaUserInfo(code string) (*OAuthUserInfo, error) {
	token, err := oktaOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}

	client := oktaOauthConfig.Client(context.Background(), token)
	resp, err := client.Get(fmt.Sprintf("https://%s/oauth2/v1/userinfo", os.Getenv("OKTA_DOMAIN")))
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user info request failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}
	userInfo.Provider = "okta"

	return &userInfo, nil
}

func ValidateOAuthToken(provider, token string) (*OAuthUserInfo, error) {
	var userInfoEndpoint string
	var client *http.Client

	switch provider {
	case "google":
		userInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
		client = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	case "okta":
		userInfoEndpoint = fmt.Sprintf("https://%s/oauth2/v1/userinfo", os.Getenv("OKTA_DOMAIN"))
		client = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		))
	default:
		return nil, errors.New("unsupported OAuth provider")
	}

	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var userInfo OAuthUserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %v", err)
	}
	userInfo.Provider = provider

	return &userInfo, nil
}
