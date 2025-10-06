package models

import "time"

// Wallet represents a user's wallet for storing and transferring funds
// @Description A digital wallet that can hold a specific currency
type Wallet struct {
	ID          int64      `json:"id" example:"1"`
	UserID      int        `json:"user_id" example:"1"`
	Currency    string     `json:"currency" example:"USD"`
	Balance     int64      `json:"balance" example:"10000"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	IsActive    bool       `json:"is_active" example:"true"`
	CanTransfer bool       `json:"can_transfer" example:"true"`
}
