package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// UserRepository is an interface that defines the methods for interacting with the User model
type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUser(userID uuid.UUID, user *models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error) // Fetch a user by email
}
