package repositories

import (
	"crawl/models"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		BaseRepository: BaseRepository[models.User]{DB: db},
	}
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *UserRepository) GetWithRoles(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Roles").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &user, err
}

func (r *UserRepository) Search(query string) ([]models.User, error) {
	var users []models.User
	err := r.DB.
		Where("username ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Find(&users).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return users, err
}

func (r *UserRepository) SetUserAsArtist(userID uuid.UUID) error {
	result := r.DB.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_artist", true)

	if result.Error != nil {
		return fmt.Errorf("failed to update user to artist: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
