package services

import (
	"verve/internal/models"
	"verve/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserService) CreateUser(username, password, pin string, roleNames []string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var hashedPin []byte
	if pin != "" {
		hashedPin, err = bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
	}

	user := &models.User{
		Username:               username,
		DisplayName:            username,
		ProfilePhotoURL:        "",
		PinRequiredForTransfer: false,
	}
	userID, err := s.userRepo.Create(user, string(hashedPassword), string(hashedPin))
	if err != nil {
		return nil, err
	}
	user.ID = userID

	for _, roleName := range roleNames {
		role, err := s.roleRepo.FindByName(roleName)
		if err != nil {
			return nil, err // Or handle more gracefully
		}
		if err := s.roleRepo.AssignToUser(userID, role.ID); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *UserService) SetPin(userID int, pin string) error {
	hashedPin, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.SetPin(userID, string(hashedPin))
}

func (s *UserService) VerifyPin(userID int, pin string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(user.PinHash), []byte(pin))
}

func (s *UserService) UpdateUser(userID int, displayName, profilePhotoURL *string, pinRequiredForTransfer *bool) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if displayName != nil {
		user.DisplayName = *displayName
	}
	if profilePhotoURL != nil {
		user.ProfilePhotoURL = *profilePhotoURL
	}
	if pinRequiredForTransfer != nil {
		user.PinRequiredForTransfer = *pinRequiredForTransfer
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetAllUsers() ([]*models.User, error) {
	return s.userRepo.FindAll()
}

func (s *UserService) Authenticate(username, password string) (int, []string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return 0, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return 0, nil, err
	}

	roles, err := s.roleRepo.GetForUser(user.ID)
	if err != nil {
		return 0, nil, err
	}

	return user.ID, roles, nil
}
