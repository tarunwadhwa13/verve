package models

import "time"

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
)

type Transfer struct {
	ID               int64          `json:"id"`
	SenderWalletID   int64          `json:"sender_wallet_id"`
	ReceiverWalletID int64          `json:"receiver_wallet_id"`
	Amount           int64          `json:"amount"`
	Status           TransferStatus `json:"status"`
	IsAnonymous      bool           `json:"is_anonymous"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
