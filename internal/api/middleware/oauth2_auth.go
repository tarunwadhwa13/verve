package middleware

import (
	"net/http"
	"strings"
	"verve/internal/auth"

	"github.com/gin-gonic/gin"
)

// OAuth2AuthMiddleware handles both JWT and OAuth2 token authentication
func OAuth2AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenType := parts[0]
		tokenString := parts[1]

		switch tokenType {
		case "Bearer":
			// Handle JWT token
			claims, err := auth.ValidateToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT token"})
				c.Abort()
				return
			}
			c.Set("userID", claims.UserID)
			c.Set("roles", claims.Roles)

		case "OAuth":
			// Handle OAuth2 token
			provider := c.GetHeader("X-OAuth-Provider")
			if provider == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "OAuth provider not specified"})
				c.Abort()
				return
			}

			userInfo, err := auth.GetOAuth2Manager().ValidateToken(c.Request.Context(), provider, tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OAuth token"})
				c.Abort()
				return
			}

			// Set user info in context
			c.Set("userID", userInfo.ID)
			c.Set("email", userInfo.Email)
			c.Set("provider", userInfo.Provider)
			// Default role for OAuth users
			c.Set("roles", []string{"user"})

		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unsupported token type"})
			c.Abort()
			return
		}

		c.Next()
	}
}
