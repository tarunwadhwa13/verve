package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresPseudonymousWalletRepository struct {
	DB *sql.DB
}

func NewPostgresPseudonymousWalletRepository(db *sql.DB) repository.PseudonymousWalletRepository {
	return &postgresPseudonymousWalletRepository{DB: db}
}

func (r *postgresPseudonymousWalletRepository) CreatePseudonymousWallet(publicKey string) (*models.PseudonymousWallet, error) {
	var id int64
	isPseudonymous := true
	if err := r.DB.QueryRow(
		"INSERT INTO wallets (public_key, is_pseudonymous, balance) VALUES ($1, $2, 0) RETURNING id",
		publicKey, isPseudonymous,
	).Scan(&id); err != nil {
		return nil, err
	}
	return &models.PseudonymousWallet{ID: id, PublicKey: publicKey, IsPseudonymous: true}, nil
}

func (r *postgresPseudonymousWalletRepository) FindByID(id int64) (*models.PseudonymousWallet, error) {
	w := &models.PseudonymousWallet{}
	if err := r.DB.QueryRow(
		"SELECT id, public_key, is_pseudonymous FROM wallets WHERE id = $1 AND is_pseudonymous = TRUE", id,
	).Scan(&w.ID, &w.PublicKey, &w.IsPseudonymous); err != nil {
		return nil, err
	}
	return w, nil
}

func (r *postgresPseudonymousWalletRepository) FindByPublicKey(publicKey string) (*models.PseudonymousWallet, error) {
	w := &models.PseudonymousWallet{}
	if err := r.DB.QueryRow(
		"SELECT id, public_key, is_pseudonymous FROM wallets WHERE public_key = $1 AND is_pseudonymous = TRUE", publicKey,
	).Scan(&w.ID, &w.PublicKey, &w.IsPseudonymous); err != nil {
		return nil, err
	}
	return w, nil
}
