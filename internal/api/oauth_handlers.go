package api

import (
	"fmt"
	"net/http"
	"verve/internal/auth"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

// handleOAuthCallback creates a handler for OAuth provider callbacks
func handleOAuthCallback(provider string, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization code and state from the callback
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing authorization code"})
			return
		}

		// Validate state
		state := c.Query("state")
		expectedState, err := c.Cookie("oauth_state")
		if err != nil || state != expectedState {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid state parameter"})
			return
		}

		// Exchange code for token and get user info
		fmt.Printf("OAuth callback for provider %s with code %s\n", provider, code)
		userInfo, err := auth.GetOAuth2Provider().HandleCallback(c.Request.Context(), provider, code)
		fmt.Printf("HandleCallback result: userInfo: %+v, err: %v\n", userInfo, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Failed to complete OAuth flow: %v", err)})
			return
		}
		if userInfo == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "No user info returned from OAuth provider"})
			return
		}

		// Debug log
		fmt.Printf("OAuth userInfo: %+v\n", userInfo)

		// Authenticate or create user
		user, token, err := authService.AuthenticateOAuth(userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("Authentication failed: %v", err)})
			return
		}
		if user == nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "No user returned from authentication"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Token: token,
			User: struct {
				Username string `json:"username" example:"john.doe@example.com"`
			}{
				Username: user.Email,
			},
		})
	}
}
