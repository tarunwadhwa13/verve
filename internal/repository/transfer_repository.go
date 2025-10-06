package repository

import "verve/internal/models"

// TransferRepository defines the interface for managing transfer requests.
type TransferRepository interface {
	Create(transfer *models.Transfer) error
	FindByID(id int64) (*models.Transfer, error)
	UpdateStatus(id int64, status models.TransferStatus) error
}
