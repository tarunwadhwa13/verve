package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"verve/internal/api"
	"verve/internal/models"
	"verve/internal/services"

	"verve/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Hash test password
	testPassword := "password"
	hash, err := auth.HashPassword(testPassword)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Initialize auth service with mock repository
	repo := newMockUserRepo()
	repo.users["admin@example.com"].PasswordHash = hash
	authService := services.NewAuthService(repo)

	// Initialize test OAuth config
	auth.InitializeTestOAuth2Config(&auth.OAuth2Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		Scopes:       []string{"email", "profile"},
		AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:     "https://oauth2.googleapis.com/token",
	})

	// Setup mock user info
	auth.SetupTestOAuthUserInfo(&auth.OAuthUserInfo{
		ID:       "123",
		Email:    "test@example.com",
		Name:     "Test User",
		Picture:  "https://example.com/photo.jpg",
		Provider: "google",
	})

	// Setup auth routes
	api.RegisterAuthRoutes(r, authService)

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
			Token string `json:"token"`
			User  struct {
				Username string `json:"username"`
			} `json:"user"`
		}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "admin@example.com", resp.User.Username)
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
		// Extract state from previous response
		location = w.Header().Get("Location")
		stateStart := strings.Index(location, "state=") + 6
		stateEnd := strings.Index(location[stateStart:], "&")
		var state string
		if stateEnd == -1 {
			state = location[stateStart:]
		} else {
			state = location[stateStart : stateStart+stateEnd]
		}

		w = httptest.NewRecorder()
		req = httptest.NewRequest("GET", fmt.Sprintf("/api/auth/google/callback?state=%s&code=test_code", state), nil)
		req.AddCookie(&http.Cookie{
			Name:  "oauth_state",
			Value: state,
			Path:  "/",
		})
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp struct {
			Token string `json:"token"`
			User  struct {
				Username string `json:"username"`
			} `json:"user"`
		}
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "test@example.com", resp.User.Username)
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
			"admin@example.com": {
				ID:       1,
				Username: "admin",
				Email:    "admin@example.com",
				Provider: "local",
				Roles:    []string{"admin"},
			},
		},
	}
}

func (m *mockUserRepo) Create(user *models.User, passwordHash, pinHash string) (int, error) {
	user.ID = len(m.users) + 1
	user.PasswordHash = passwordHash
	m.users[user.Email] = user // Store by email since that's how OAuth looks up users
	return user.ID, nil
}

func (m *mockUserRepo) FindByID(id int) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindByUsername(username string) (*models.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindByEmailAndProvider(email, provider string) (*models.User, error) {
	if user, ok := m.users[email]; ok && user.Provider == provider {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindByProviderID(providerID string) (*models.User, error) {
	for _, user := range m.users {
		if user.ProviderUserID == providerID {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserRepo) SetPin(userID int, pinHash string) error {
	for _, user := range m.users {
		if user.ID == userID {
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

func (m *mockUserRepo) Update(user *models.User) error {
	if _, ok := m.users[user.Email]; ok {
		m.users[user.Email] = user
		return nil
	}
	return fmt.Errorf("user not found")
}

func (m *mockUserRepo) FindAll() ([]*models.User, error) {
	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}
