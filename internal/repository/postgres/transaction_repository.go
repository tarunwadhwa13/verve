package postgres

import (
	"database/sql"
	"errors"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresTransactionRepository struct {
	DB *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &postgresTransactionRepository{DB: db}
}

func (r *postgresTransactionRepository) TransferCoins(senderWalletID, receiverWalletID, amount int64) (*models.Transaction, []*models.LedgerEntry, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var senderBalance int64
	if err = tx.QueryRow("SELECT balance FROM wallets WHERE id = $1 FOR UPDATE", senderWalletID).Scan(&senderBalance); err != nil {
		return nil, nil, err
	}
	if senderBalance < amount {
		return nil, nil, errors.New("insufficient funds")
	}

	if _, err = tx.Exec("UPDATE wallets SET balance = balance - $1 WHERE id = $2", amount, senderWalletID); err != nil {
		return nil, nil, err
	}
	if _, err = tx.Exec("UPDATE wallets SET balance = balance + $1 WHERE id = $2", amount, receiverWalletID); err != nil {
		return nil, nil, err
	}

	var transactionID int64
	if err = tx.QueryRow(
		"INSERT INTO transactions (sender_wallet_id, receiver_wallet_id, amount) VALUES ($1, $2, $3) RETURNING id, created_at",
		senderWalletID, receiverWalletID, amount,
	).Scan(&transactionID, new(string)); err != nil {
		return nil, nil, err
	}

	// Ledger entries
	debitEntry := &models.LedgerEntry{
		TransactionID: transactionID,
		WalletID:      senderWalletID,
		EntryType:     "debit",
		Amount:        amount,
	}
	creditEntry := &models.LedgerEntry{
		TransactionID: transactionID,
		WalletID:      receiverWalletID,
		EntryType:     "credit",
		Amount:        amount,
	}

	if _, err = tx.Exec(
		"INSERT INTO ledger_entries (transaction_id, wallet_id, entry_type, amount) VALUES ($1, $2, $3, $4)",
		debitEntry.TransactionID, debitEntry.WalletID, debitEntry.EntryType, debitEntry.Amount,
	); err != nil {
		return nil, nil, err
	}
	if _, err = tx.Exec(
		"INSERT INTO ledger_entries (transaction_id, wallet_id, entry_type, amount) VALUES ($1, $2, $3, $4)",
		creditEntry.TransactionID, creditEntry.WalletID, creditEntry.EntryType, creditEntry.Amount,
	); err != nil {
		return nil, nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, err
	}

	transaction := &models.Transaction{
		ID:               transactionID,
		SenderWalletID:   senderWalletID,
		ReceiverWalletID: receiverWalletID,
		Amount:           amount,
	}

	return transaction, []*models.LedgerEntry{debitEntry, creditEntry}, nil
}
