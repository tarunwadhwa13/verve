package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresRoleRepository struct {
	DB *sql.DB
}

func NewPostgresRoleRepository(db *sql.DB) repository.RoleRepository {
	return &postgresRoleRepository{DB: db}
}

func (r *postgresRoleRepository) FindByName(name string) (*models.Role, error) {
	role := &models.Role{}
	err := r.DB.QueryRow("SELECT id FROM roles WHERE name = $1", name).Scan(&role.ID)
	if err != nil {
		return nil, err
	}
	role.Name = name
	return role, nil
}

func (r *postgresRoleRepository) AssignToUser(userID, roleID int) error {
	_, err := r.DB.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
	return err
}

func (r *postgresRoleRepository) GetForUser(userID int) ([]string, error) {
	rows, err := r.DB.Query("SELECT r.name FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}
