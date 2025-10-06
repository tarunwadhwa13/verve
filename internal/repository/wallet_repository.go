package repository

import "verve/internal/models"

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByUserID(userID int) ([]models.Wallet, error)
	FindByID(id int64) (*models.Wallet, error)
}
