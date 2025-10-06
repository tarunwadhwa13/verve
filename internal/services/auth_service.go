package services

import (
	"errors"
	"verve/internal/auth"
	"verve/internal/models"
	"verve/internal/repository"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// AuthenticateLocal handles username/password authentication
func (s *AuthService) AuthenticateLocal(username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Only allow local auth for users with local provider
	if user.Provider != "" && user.Provider != "local" {
		return nil, "", errors.New("this account uses " + user.Provider + " authentication")
	}

	if !auth.ValidatePassword(password, user.PasswordHash) {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Roles)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// AuthenticateOAuth handles OAuth authentication
func (s *AuthService) AuthenticateOAuth(userInfo *auth.OAuthUserInfo) (*models.User, string, error) {
	// Try to find existing user by email and provider
	user, err := s.userRepo.FindByEmailAndProvider(userInfo.Email, userInfo.Provider)
	if err != nil {
		// Create new user if not found
		user = &models.User{
			Email:           userInfo.Email,
			DisplayName:     userInfo.Name,
			ProfilePhotoURL: userInfo.Picture,
			Provider:        userInfo.Provider,
			ProviderUserID:  userInfo.ID,
			Roles:           []string{"user"}, // Default role for OAuth users
		}

		_, err = s.userRepo.Create(user, "", "") // No password/pin for OAuth users
		if err != nil {
			return nil, "", err
		}
	} else {
		// Update existing user's OAuth info if needed
		needsUpdate := false
		if user.DisplayName != userInfo.Name {
			user.DisplayName = userInfo.Name
			needsUpdate = true
		}
		if user.ProfilePhotoURL != userInfo.Picture {
			user.ProfilePhotoURL = userInfo.Picture
			needsUpdate = true
		}
		if user.ProviderUserID != userInfo.ID {
			user.ProviderUserID = userInfo.ID
			needsUpdate = true
		}

		if needsUpdate {
			err = s.userRepo.Update(user)
			if err != nil {
				return nil, "", err
			}
		}
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Roles)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// CreateAdminUser creates a new admin user with local authentication
func (s *AuthService) CreateAdminUser(username, email, password string) (*models.User, error) {
	// Hash password
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Provider:     "local",
		Roles:        []string{"admin"},
	}

	_, err = s.userRepo.Create(user, passwordHash, "") // Create with password but no pin
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LinkOAuthToLocal links an OAuth account to an existing local account
func (s *AuthService) LinkOAuthToLocal(userID int, oauthProvider, oauthUserID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if user.Provider != "local" {
		return errors.New("can only link OAuth to local accounts")
	}

	// Check if OAuth account is already linked to another user
	exists, err := s.userRepo.FindByProviderID(oauthUserID)
	if err == nil && exists.ID != userID {
		return errors.New("oauth account already linked to another user")
	}

	user.Provider = oauthProvider
	user.ProviderUserID = oauthUserID
	return s.userRepo.Update(user)
}

// UnlinkOAuth removes OAuth provider from a user account
func (s *AuthService) UnlinkOAuth(userID int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	// Only allow unlinking if user has a password set
	if user.PasswordHash == "" {
		return errors.New("cannot unlink OAuth without setting a password first")
	}

	user.Provider = "local"
	user.ProviderUserID = ""
	return s.userRepo.Update(user)
}
