package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"verve/internal/auth"
	"verve/internal/models"
	"verve/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Initialize auth service with mock repository
	authService := services.NewAuthService(newMockUserRepo())

	// Setup auth routes
	auth.RegisterAuthRoutes(r, authService)

	t.Run("Local Authentication", func(t *testing.T) {
		// Test local login
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{
			"username": "admin",
			"password": "password"
		}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "admin", resp.User.Username)
	})

	t.Run("OAuth2 Flow", func(t *testing.T) {
		// Test OAuth login redirect
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/auth/google/login", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		assert.Contains(t, location, "accounts.google.com")

		// Test OAuth callback
		state := w.Header().Get("Set-Cookie") // Extract state from cookie
		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/api/auth/google/callback?state="+state+"&code=test_code", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Token string      `json:"token"`
			User  models.User `json:"user"`
		}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "google", resp.User.Provider)
	})

	t.Run("Account Linking", func(t *testing.T) {
		// Login first to get token
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(`{
			"username": "admin",
			"password": "password"
		}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		var loginResp struct {
			Token string `json:"token"`
		}
		err := json.NewDecoder(w.Body).Decode(&loginResp)
		assert.NoError(t, err)

		// Try to link OAuth account
		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/api/auth/link/oauth?provider=google", nil)
		req.Header.Set("Authorization", "Bearer "+loginResp.Token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		location := w.Header().Get("Location")
		assert.Contains(t, location, "accounts.google.com")
	})
}

// Mock user repository for testing
type mockUserRepo struct {
	users map[string]*models.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: map[string]*models.User{
			"admin": {
				ID:           1,
				Username:     "admin",
				PasswordHash: "$2a$10$xxx", // pre-hashed "password"
				Provider:     "local",
				Roles:        []string{"admin"},
			},
		},
	}
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	if user, ok := m.users[username]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

// Implement other repository methods as needed for testing
