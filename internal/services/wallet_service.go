package services

import (
	"verve/internal/models"
	"verve/internal/repository"
)

type WalletService struct {
	walletRepo repository.WalletRepository
}

func NewWalletService(walletRepo repository.WalletRepository) *WalletService {
	return &WalletService{walletRepo: walletRepo}
}

func (s *WalletService) CreateWallet(userID int, currency string) (*models.Wallet, error) {
	wallet := &models.Wallet{
		UserID:   userID,
		Currency: currency,
		Balance:  0,
	}
	err := s.walletRepo.Create(wallet)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *WalletService) GetWalletsByUserID(userID int) ([]models.Wallet, error) {
	return s.walletRepo.FindByUserID(userID)
}

func (s *WalletService) GetWalletByID(id int64) (*models.Wallet, error) {
	return s.walletRepo.FindByID(id)
}
