package models

import "time"

// User represents a user in the system
// @Description User information including display name and profile settings
type User struct {
	ID                     int       `json:"id" example:"1"`
	Username               string    `json:"username" example:"john.doe"`
	Email                  string    `json:"email" example:"john.doe@example.com"`
	PasswordHash           string    `json:"-"` // Sensitive data, not exposed in API
	PinHash                string    `json:"-"` // Sensitive data, not exposed in API
	DisplayName            string    `json:"display_name" example:"John Doe"`
	ProfilePhotoURL        string    `json:"profile_photo_url" example:"https://example.com/photo.jpg"`
	PinRequiredForTransfer bool      `json:"pin_required_for_transfer" example:"true"`
	Provider               string    `json:"provider" example:"google"`  // OAuth provider (google, okta, or local)
	ProviderUserID         string    `json:"provider_user_id,omitempty"` // ID from the OAuth provider
	Roles                  []string  `json:"roles" example:"['user','admin']"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Permission struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
