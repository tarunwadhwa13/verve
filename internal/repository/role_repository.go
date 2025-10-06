package repository

import "verve/internal/models"

type RoleRepository interface {
	FindByName(name string) (*models.Role, error)
	AssignToUser(userID, roleID int) error
	GetForUser(userID int) ([]string, error)
}
