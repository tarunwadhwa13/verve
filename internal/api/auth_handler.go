package api

import (
	"fmt"
	"net/http"
	"verve/internal/api/middleware"
	"verve/internal/auth"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication routes
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authRoutes := router.Group("/api/auth")
	{
		// Local authentication
		authRoutes.POST("/login", LocalLoginHandler(authService))

		// OAuth routes
		authRoutes.GET("/google/login", GoogleLoginHandler())
		authRoutes.GET("/google/callback", handleOAuthCallback("google", authService))
		authRoutes.GET("/okta/login", OktaLoginHandler())
		authRoutes.GET("/okta/callback", handleOAuthCallback("okta", authService))

		// Account linking (protected routes)
		protected := authRoutes.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/link/oauth", LinkOAuthHandler(authService))
			protected.POST("/unlink/oauth", UnlinkOAuthHandler(authService))
		}
	}
}

// Removed duplicate types - using definitions from api_types.go

// LocalLoginHandler handles username/password login
// @Summary Login with username and password
// @Description Authenticate user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func LocalLoginHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, token, err := authService.AuthenticateLocal(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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

// LinkOAuthHandler links an OAuth account to the current user
// @Summary Link OAuth account
// @Description Link an OAuth account to the current user's account
// @Tags auth
// @Accept json
// @Produce json
// @Param provider query string true "OAuth provider (google/okta)"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /auth/link/oauth [post]
func LinkOAuthHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		provider := c.Query("provider")
		if provider == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Provider is required"})
			return
		}

		// Start OAuth flow for linking
		state := auth.GenerateState()
		c.SetCookie("oauth_state", state, 3600, "/", "", false, true)
		c.SetCookie("oauth_action", "link", 3600, "/", "", false, true)
		c.SetCookie("oauth_user_id", fmt.Sprintf("%d", userID), 3600, "/", "", false, true)

		var (
			url string
			err error
		)
		switch provider {
		case "google":
			url, err = auth.GetOAuth2Manager().GetAuthURL("google", state)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auth URL: " + err.Error()})
				return
			}
		case "okta":
			url, err = auth.GetOAuth2Manager().GetAuthURL("okta", state)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auth URL: " + err.Error()})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider"})
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// UnlinkOAuthHandler removes OAuth provider from user account
// @Summary Unlink OAuth account
// @Description Remove OAuth provider from user account
// @Tags auth
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Security ApiKeyAuth
// @Router /auth/unlink/oauth [post]
func UnlinkOAuthHandler(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt("userID")
		err := authService.UnlinkOAuth(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "OAuth provider unlinked successfully"})
	}
}
