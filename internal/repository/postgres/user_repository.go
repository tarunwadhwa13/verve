package postgres

import (
	"database/sql"
	"verve/internal/models"
	"verve/internal/repository"
)

type postgresUserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{DB: db}
}

func (r *postgresUserRepository) Create(user *models.User, passwordHash, pinHash string) (int, error) {
	var userID int
	err := r.DB.QueryRow(`
		INSERT INTO users (username, email, password_hash, pin_hash, display_name, profile_photo_url, 
			provider, provider_user_id, pin_required_for_transfer) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id`,
		user.Username, user.Email, passwordHash, pinHash, user.DisplayName,
		user.ProfilePhotoURL, user.Provider, user.ProviderUserID,
		user.PinRequiredForTransfer).Scan(&userID)
	return userID, err
}

func (r *postgresUserRepository) FindByID(id int) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, pin_hash, display_name, profile_photo_url, 
			provider, provider_user_id, pin_required_for_transfer, created_at, updated_at 
		FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.PinHash,
		&user.DisplayName, &user.ProfilePhotoURL, &user.Provider, &user.ProviderUserID,
		&user.PinRequiredForTransfer, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) FindByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, pin_hash, display_name, profile_photo_url, 
			provider, provider_user_id, pin_required_for_transfer, created_at, updated_at 
		FROM users WHERE username = $1`, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.PinHash,
		&user.DisplayName, &user.ProfilePhotoURL, &user.Provider, &user.ProviderUserID,
		&user.PinRequiredForTransfer, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) FindByEmailAndProvider(email, provider string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, pin_hash, display_name, profile_photo_url, 
			provider, provider_user_id, pin_required_for_transfer, created_at, updated_at 
		FROM users WHERE email = $1 AND provider = $2`, email, provider).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.PinHash,
		&user.DisplayName, &user.ProfilePhotoURL, &user.Provider, &user.ProviderUserID,
		&user.PinRequiredForTransfer, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) FindByProviderID(providerID string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.QueryRow(`
		SELECT id, username, email, password_hash, pin_hash, display_name, profile_photo_url, 
			provider, provider_user_id, pin_required_for_transfer, created_at, updated_at 
		FROM users WHERE provider_user_id = $1`, providerID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.PinHash,
		&user.DisplayName, &user.ProfilePhotoURL, &user.Provider, &user.ProviderUserID,
		&user.PinRequiredForTransfer, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) SetPin(userID int, pinHash string) error {
	_, err := r.DB.Exec("UPDATE users SET pin_hash = $1 WHERE id = $2", pinHash, userID)
	return err
}

func (r *postgresUserRepository) Update(user *models.User) error {
	_, err := r.DB.Exec(`
		UPDATE users SET 
			username = $1, email = $2, display_name = $3, profile_photo_url = $4,
			provider = $5, provider_user_id = $6, pin_required_for_transfer = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $8`,
		user.Username, user.Email, user.DisplayName, user.ProfilePhotoURL,
		user.Provider, user.ProviderUserID, user.PinRequiredForTransfer,
		user.ID)
	return err
}

func (r *postgresUserRepository) FindAll() ([]*models.User, error) {
	rows, err := r.DB.Query(`
		SELECT id, username, email, password_hash, pin_hash, display_name, profile_photo_url,
			provider, provider_user_id, pin_required_for_transfer, created_at, updated_at 
		FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.PinHash,
			&user.DisplayName, &user.ProfilePhotoURL, &user.Provider, &user.ProviderUserID,
			&user.PinRequiredForTransfer, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
