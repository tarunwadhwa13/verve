package services

import (
	"verve/internal/models"
)

func (s *UserService) UpsertOAuthUser(email, name, picture, provider string) (*models.User, error) {
	// Try to find existing user by email and provider
	user, err := s.userRepo.FindByEmailAndProvider(email, provider)
	if err == nil {
		// User exists, update profile if needed
		needsUpdate := false
		if user.DisplayName != name {
			user.DisplayName = name
			needsUpdate = true
		}
		if user.ProfilePhotoURL != picture {
			user.ProfilePhotoURL = picture
			needsUpdate = true
		}

		if needsUpdate {
			err = s.userRepo.Update(user)
			if err != nil {
				return nil, err
			}
		}
		return user, nil
	}

	// Create new user
	user = &models.User{
		Email:           email,
		DisplayName:     name,
		ProfilePhotoURL: picture,
		Provider:        provider,
		Roles:           []string{"user"}, // Default role
	}

	_, err = s.userRepo.Create(user, "", "") // No password/pin for OAuth users
	if err != nil {
		return nil, err
	}

	return user, nil
}
