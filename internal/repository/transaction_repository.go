package repository

import "verve/internal/models"

// TransactionRepository abstracts coin transfer and ledger logging
// All operations are performed atomically in a DB transaction

type TransactionRepository interface {
	TransferCoins(senderWalletID, receiverWalletID, amount int64) (*models.Transaction, []*models.LedgerEntry, error)
}

// LedgerRepository abstracts append-only logging for anonymous transfers
type LedgerRepository interface {
	LogAnonymousTransfer(senderWalletID, receiverWalletID, amount int64, pubKey interface{}, signature string) error
}
