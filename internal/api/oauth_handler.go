package api

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"verve/internal/auth"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterOAuthRoutes(router *gin.Engine, userService *services.UserService) {
	authRoutes := router.Group("/api/auth")
	{
		// Google OAuth routes
		authRoutes.GET("/google/login", GoogleLoginHandler())
		authRoutes.GET("/google/callback", GoogleCallbackHandler(userService))

		// Okta OAuth routes
		authRoutes.GET("/okta/login", OktaLoginHandler())
		authRoutes.GET("/okta/callback", OktaCallbackHandler(userService))
	}
}

// generateState generates a random state string for OAuth CSRF protection
func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GoogleLoginHandler initiates Google OAuth login
func GoogleLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := generateState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate state"})
			return
		}

		c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
		url, err := auth.GetOAuth2Manager().GetAuthURL("google", state)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auth URL"})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GoogleCallbackHandler handles Google OAuth callback
func GoogleCallbackHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		storedState, _ := c.Cookie("oauth_state")
		if state != storedState {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid state"})
			return
		}

		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Code not provided"})
			return
		}

		userInfo, err := auth.GetOAuth2Manager().HandleCallback(c.Request.Context(), "google", code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get user info"})
			return
		}

		user, err := userService.UpsertOAuthUser(userInfo.Email, userInfo.Name, userInfo.Picture, "google")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create/update user"})
			return
		}

		token, err := auth.GenerateJWT(user.ID, user.Roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
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

// OktaLoginHandler initiates Okta OAuth login
func OktaLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		state, err := generateState()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate state"})
			return
		}

		c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
		url, err := auth.GetOAuth2Manager().GetAuthURL("okta", state)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auth URL"})
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// OktaCallbackHandler handles Okta OAuth callback
func OktaCallbackHandler(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := c.Query("state")
		storedState, _ := c.Cookie("oauth_state")
		if state != storedState {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid state"})
			return
		}

		code := c.Query("code")
		if code == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Code not provided"})
			return
		}

		userInfo, err := auth.GetOAuth2Manager().HandleCallback(c.Request.Context(), "okta", code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get user info"})
			return
		}

		user, err := userService.UpsertOAuthUser(userInfo.Email, userInfo.Name, userInfo.Picture, "okta")
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create/update user"})
			return
		}

		token, err := auth.GenerateJWT(user.ID, user.Roles)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to generate token"})
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
