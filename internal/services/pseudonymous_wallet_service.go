package services

import (
	"verve/internal/repository"
	"verve/pkg/utils"
)

type PseudonymousWalletService struct {
	repo repository.PseudonymousWalletRepository
}

func NewPseudonymousWalletService(repo repository.PseudonymousWalletRepository) *PseudonymousWalletService {
	return &PseudonymousWalletService{repo: repo}
}

// CreatePseudonymousWallet generates a keypair, stores public key, returns wallet and private key
func (s *PseudonymousWalletService) CreatePseudonymousWallet() (walletID int64, publicKey string, privateKey string, err error) {
	priv, err := utils.GenerateECDSAKeyPair()
	if err != nil {
		return 0, "", "", err
	}
	pubStr := utils.PublicKeyToString(&priv.PublicKey)
	w, err := s.repo.CreatePseudonymousWallet(pubStr)
	if err != nil {
		return 0, "", "", err
	}
	// TODO: Securely handle private key
	return w.ID, pubStr, "", nil
}
