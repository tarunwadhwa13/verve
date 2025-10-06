package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresTransferRepository struct {
	DB *sql.DB
}

func NewPostgresTransferRepository(db *sql.DB) repository.TransferRepository {
	return &postgresTransferRepository{DB: db}
}

func (r *postgresTransferRepository) Create(transfer *models.Transfer) error {
	query := `
		INSERT INTO transfers (sender_wallet_id, receiver_wallet_id, amount, status, is_anonymous)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.DB.QueryRow(
		query,
		transfer.SenderWalletID,
		transfer.ReceiverWalletID,
		transfer.Amount,
		transfer.Status,
		transfer.IsAnonymous,
	).Scan(&transfer.ID, &transfer.CreatedAt, &transfer.UpdatedAt)

	return err
}

func (r *postgresTransferRepository) FindByID(id int64) (*models.Transfer, error) {
	transfer := &models.Transfer{}
	query := `
		SELECT id, sender_wallet_id, receiver_wallet_id, amount, status, is_anonymous, created_at, updated_at
		FROM transfers
		WHERE id = $1`

	err := r.DB.QueryRow(query, id).Scan(
		&transfer.ID,
		&transfer.SenderWalletID,
		&transfer.ReceiverWalletID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.IsAnonymous,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (r *postgresTransferRepository) UpdateStatus(id int64, status models.TransferStatus) error {
	query := "UPDATE transfers SET status = $1 WHERE id = $2"
	_, err := r.DB.Exec(query, status, id)
	return err
}
