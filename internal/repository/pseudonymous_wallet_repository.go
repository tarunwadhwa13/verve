package repository

import "verve/internal/models"

// PseudonymousWalletRepository abstracts creation and lookup of pseudonymous wallets

type PseudonymousWalletRepository interface {
	CreatePseudonymousWallet(publicKey string) (*models.PseudonymousWallet, error)
	FindByID(id int64) (*models.PseudonymousWallet, error)
	FindByPublicKey(publicKey string) (*models.PseudonymousWallet, error)
}
