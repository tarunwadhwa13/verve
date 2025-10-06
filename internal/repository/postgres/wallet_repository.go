package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresWalletRepository struct {
	DB *sql.DB
}

func NewPostgresWalletRepository(db *sql.DB) repository.WalletRepository {
	return &postgresWalletRepository{DB: db}
}

func (r *postgresWalletRepository) Create(wallet *models.Wallet) error {
	return r.DB.QueryRow(
		"INSERT INTO wallets (user_id, currency, balance) VALUES ($1, $2, $3) RETURNING id, created_at",
		wallet.UserID, wallet.Currency, wallet.Balance,
	).Scan(&wallet.ID, &wallet.CreatedAt)
}

func (r *postgresWalletRepository) FindByUserID(userID int) ([]models.Wallet, error) {
	rows, err := r.DB.Query("SELECT id, user_id, balance, currency, created_at FROM wallets WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []models.Wallet
	for rows.Next() {
		var wallet models.Wallet
		if err := rows.Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt); err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	return wallets, nil
}

func (r *postgresWalletRepository) FindByID(id int64) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.DB.QueryRow(
		"SELECT id, user_id, balance, currency, created_at FROM wallets WHERE id = $1",
		id,
	).Scan(&wallet.ID, &wallet.UserID, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
