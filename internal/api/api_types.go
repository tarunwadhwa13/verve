package api

import "verve/internal/models"

type (
	// Common Response Types
	ErrorResponse struct {
		Error string `json:"error" example:"Invalid request"`
	}

	SuccessResponse struct {
		Message string `json:"message" example:"Operation successful"`
	}

	// Badge Related Types
	CreateBadgeRequest struct {
		Name        string                   `json:"name" binding:"required" example:"Achievement Master"`
		Description string                   `json:"description" binding:"required" example:"Awarded to users who complete all achievements"`
		IconURL     string                   `json:"icon_url" example:"https://example.com/badges/master.png"`
		Points      int                      `json:"points" example:"100"`
		Rules       []models.AchievementRule `json:"rules"`
	}

	UpdateBadgeRequest struct {
		Name        *string `json:"name" example:"Achievement Master"`
		Description *string `json:"description" example:"Awarded to users who complete all achievements"`
		IconURL     *string `json:"icon_url" example:"https://example.com/badges/master.png"`
		Points      *int    `json:"points" example:"100"`
		IsActive    *bool   `json:"is_active" example:"true"`
	}

	// User Related Types
	UpdateUserRequest struct {
		DisplayName            *string `json:"display_name" example:"John Doe"`
		ProfilePhotoURL        *string `json:"profile_photo_url" example:"https://example.com/photo.jpg"`
		PinRequiredForTransfer *bool   `json:"pin_required_for_transfer" example:"true"`
	}

	CreateUserRequest struct {
		Username string   `json:"username" binding:"required" example:"john.doe"`
		Password string   `json:"password" binding:"required" example:"secure_password"`
		Pin      string   `json:"pin" example:"1234"`
		Roles    []string `json:"roles" example:"['user','admin']"`
	}

	LoginRequest struct {
		Username string `json:"email" binding:"required" example:"john.doe@example.com"`
		Password string `json:"password" binding:"required" example:"secure_password"`
	}

	LoginResponse struct {
		Token        string `json:"token" example:"eyJhbGciOiJS..."`
		RefreshToken string `json:"refresh_token" example:""`
		User         struct {
			Username string `json:"username" example:"john.doe@example.com"`
		} `json:"user"`
	}

	// Wallet Related Types
	CreateWalletRequest struct {
		Currency string `json:"currency" binding:"required" example:"USD"`
	}

	// Transfer Related Types
	TransferRequest struct {
		SenderWalletID   int64  `json:"sender_wallet_id" binding:"required" example:"1"`
		ReceiverWalletID int64  `json:"receiver_wallet_id" binding:"required" example:"2"`
		Amount           int64  `json:"amount" binding:"required" example:"1000"`
		IsAnonymous      bool   `json:"is_anonymous" example:"false"`
		Pin              string `json:"pin" example:"1234"`
	}

	TransferResponse struct {
		Transfer *models.Transfer `json:"transfer"`
	}

	// Badge Related Types
	BadgeResponse struct {
		Badge models.Badge             `json:"badge"`
		Rules []models.AchievementRule `json:"rules"`
	}

	// Pseudonymous Wallet Related Types
	CreatePseudonymousWalletResponse struct {
		WalletID   int64  `json:"wallet_id"`
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"` // Only returned once
	}
)
