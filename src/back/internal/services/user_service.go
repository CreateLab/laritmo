package services

import (
	"errors"

	"github.com/CreateLab/laritmo/internal/models"
	"github.com/CreateLab/laritmo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCannotModifyOwner     = errors.New("cannot modify owner account")
	ErrCannotDeleteOwner     = errors.New("cannot delete owner account")
	ErrInvalidOldPassword    = errors.New("invalid old password")
	ErrCannotResetOwnPass    = errors.New("use password change endpoint to update your own password")
	ErrCannotResetOwnerPass  = errors.New("cannot reset password for owner account")
	ErrCannotDeactivateOwner = errors.New("cannot deactivate owner account")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAllAdmins returns all users with owner or admin role
func (s *UserService) GetAllAdmins() ([]*models.User, error) {
	return s.repo.GetAllAdmins()
}

// GetUserByID returns a user by their ID
func (s *UserService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         req.Role,
		IsActive:     true,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user fields (except for owner accounts)
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user.Role == models.RoleOwner {
		return ErrCannotModifyOwner
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return s.repo.Update(user)
}

// DeleteUser deletes a user (except for owner accounts)
func (s *UserService) DeleteUser(id int) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user.Role == models.RoleOwner {
		return ErrCannotDeleteOwner
	}

	return s.repo.Delete(id)
}

// ChangePassword allows a user to change their own password
func (s *UserService) ChangePassword(userID int, oldPass, newPass string) error {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPass))
	if err != nil {
		return ErrInvalidOldPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, string(passwordHash))
}

// ResetPassword allows an admin to reset another user's password
func (s *UserService) ResetPassword(adminID, targetUserID int, newPass string) error {
	if adminID == targetUserID {
		return ErrCannotResetOwnPass
	}

	targetUser, err := s.repo.GetByID(targetUserID)
	if err != nil {
		return err
	}

	if targetUser.Role == models.RoleOwner {
		return ErrCannotResetOwnerPass
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(targetUserID, string(passwordHash))
}

// DeactivateUser deactivates a user account (except for owner accounts)
func (s *UserService) DeactivateUser(id int) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if user.Role == models.RoleOwner {
		return ErrCannotDeactivateOwner
	}

	user.IsActive = false
	return s.repo.Update(user)
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(id int) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.repo.Update(user)
}
