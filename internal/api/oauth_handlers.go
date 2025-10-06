package api

import (
	"net/http"
	"verve/internal/auth"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

// handleOAuthCallback creates a handler for OAuth provider callbacks
func handleOAuthCallback(provider string, authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization code from the callback
		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Missing authorization code"})
			return
		}

		// Exchange code for token and get user info
		userInfo, err := auth.GetOAuth2Manager().HandleCallback(c.Request.Context(), provider, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to complete OAuth flow: " + err.Error()})
			return
		}

		// Authenticate or create user
		user, token, err := authService.AuthenticateOAuth(userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Authentication failed: " + err.Error()})
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
