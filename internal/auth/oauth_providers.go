package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getGoogleUserInfo(ctx context.Context, client *http.Client) (*OAuthUserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Google: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info from Google: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Google user info: %v", err)
	}

	return &OAuthUserInfo{
		ID:       result.Sub,
		Email:    result.Email,
		Name:     result.Name,
		Picture:  result.Picture,
		Provider: "google",
	}, nil
}

func getOktaUserInfo(ctx context.Context, client *http.Client) (*OAuthUserInfo, error) {
	// Get the user info endpoint from environment variable or configuration
	userInfoEndpoint := fmt.Sprintf("https://%s/oauth2/v1/userinfo", os.Getenv("OKTA_DOMAIN"))

	resp, err := client.Get(userInfoEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from Okta: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info from Okta: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode Okta user info: %v", err)
	}

	return &OAuthUserInfo{
		ID:       result.Sub,
		Email:    result.Email,
		Name:     result.Name,
		Picture:  result.Picture,
		Provider: "okta",
	}, nil
}
