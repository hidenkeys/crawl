package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepositoryImpl is the implementation of the UserRepository interface
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

// CreateUser inserts a new user into the database
func (r *UserRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	err := r.DB.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID fetches a user by their ID
func (r *UserRepositoryImpl) GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user by their ID
func (r *UserRepositoryImpl) UpdateUser(userID uuid.UUID, user *models.User) (*models.User, error) {
	var existingUser models.User
	err := r.DB.First(&existingUser, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}

	// Update the fields
	err = r.DB.Model(&existingUser).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return &existingUser, nil
}

// DeleteUser deletes a user by their ID
func (r *UserRepositoryImpl) DeleteUser(userID uuid.UUID) error {
	var user models.User
	err := r.DB.First(&user, "id = ?", userID).Error
	if err != nil {
		return err
	}
	return r.DB.Delete(&user).Error
}

// GetUserByUsername fetches a user by their username
func (r *UserRepositoryImpl) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail fetches a user by their email
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
