package repository

import "verve/internal/models"

type UserRepository interface {
	Create(user *models.User, passwordHash, pinHash string) (int, error)
	FindByID(id int) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmailAndProvider(email, provider string) (*models.User, error)
	FindByProviderID(providerID string) (*models.User, error)
	SetPin(userID int, pinHash string) error
	Update(user *models.User) error
	FindAll() ([]*models.User, error)
}
