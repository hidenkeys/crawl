package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"github.com/google/uuid"
)

// UserService provides the logic for handling users
type UserService struct {
	UserRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return s.UserRepo.CreateUser(user)
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.UserRepo.GetUserByID(userID)
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, userID uuid.UUID, user *models.User) (*models.User, error) {
	return s.UserRepo.UpdateUser(userID, user)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.UserRepo.DeleteUser(userID)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.UserRepo.GetUserByUsername(username)
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.UserRepo.GetUserByEmail(email)
}
