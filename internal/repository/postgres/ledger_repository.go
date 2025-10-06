package postgres

import (
	"database/sql"
	"verve/internal/repository"
)

type postgresLedgerRepository struct {
	DB *sql.DB
}

func NewPostgresLedgerRepository(db *sql.DB) repository.LedgerRepository {
	return &postgresLedgerRepository{DB: db}
}

// LogAnonymousTransfer logs an anonymous transfer to the ledger (append-only)
func (r *postgresLedgerRepository) LogAnonymousTransfer(senderWalletID, receiverWalletID, amount int64, pubKey interface{}, signature string) error {
	_, err := r.DB.Exec(
		`INSERT INTO ledger_entries (transaction_id, wallet_id, entry_type, amount, created_at, extra) VALUES (NULL, $1, 'debit', $3, CURRENT_TIMESTAMP, $5), (NULL, $2, 'credit', $3, CURRENT_TIMESTAMP, $6)`,
		senderWalletID, receiverWalletID, amount, senderWalletID, signature, pubKey,
	)
	return err
}
