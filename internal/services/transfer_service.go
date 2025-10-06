package services

import (
	"errors"
	"verve/internal/models"
	"verve/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

// TransferService orchestrates the creation and execution of transfers.
type TransferService struct {
	transferRepo repository.TransferRepository
	txRepo       repository.TransactionRepository
	ledgerRepo   repository.LedgerRepository
	userRepo     repository.UserRepository
}

// NewTransferService creates a new TransferService.
func NewTransferService(
	transferRepo repository.TransferRepository,
	txRepo repository.TransactionRepository,
	ledgerRepo repository.LedgerRepository,
	userRepo repository.UserRepository,
) *TransferService {
	return &TransferService{
		transferRepo: transferRepo,
		txRepo:       txRepo,
		ledgerRepo:   ledgerRepo,
		userRepo:     userRepo,
	}
}

// InitiateTransfer creates a new transfer record and processes it.
func (s *TransferService) InitiateTransfer(
	userID int,
	senderWalletID, receiverWalletID, amount int64,
	isAnonymous bool,
	pin string,
) (*models.Transfer, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.PinRequiredForTransfer {
		if err := bcrypt.CompareHashAndPassword([]byte(user.PinHash), []byte(pin)); err != nil {
			return nil, errors.New("invalid PIN")
		}
	}
	// ... existing code ...
	transfer := &models.Transfer{
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		Amount:           amount,
		Status:           models.TransferStatusPending,
		IsAnonymous:      isAnonymous,
	}
	// ... existing code ...
	if isAnonymous && s.ledgerRepo != nil {
		_ = s.ledgerRepo.LogAnonymousTransfer(senderWalletID, receiverWalletID, amount, nil, "")
	}

	return transfer, nil
}

// GetTransferStatus retrieves the status of a transfer by its ID.
func (s *TransferService) GetTransferStatus(id int64) (*models.Transfer, error) {
	return s.transferRepo.FindByID(id)
}
