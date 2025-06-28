package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserFavoriteRepository struct {
	DB *gorm.DB
}

func NewUserFavoriteRepository(db *gorm.DB) IUserFavoriteRepository {
	return &UserFavoriteRepository{DB: db}
}

func (r *UserFavoriteRepository) AddFavorite(userID, songID uuid.UUID) error {
	return r.DB.Create(&models.UserFavorite{
		UserID: userID,
		SongID: songID,
	}).Error
}

func (r *UserFavoriteRepository) RemoveFavorite(userID, songID uuid.UUID) error {
	return r.DB.
		Where("user_id = ? AND song_id = ?", userID, songID).
		Delete(&models.UserFavorite{}).
		Error
}

func (r *UserFavoriteRepository) GetUserFavorites(userID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	err := r.DB.
		Joins("JOIN user_favorites ON songs.id = user_favorites.song_id").
		Where("user_favorites.user_id = ?", userID).
		Preload("Artist").
		Find(&songs).
		Error
	return songs, err
}

func (r *UserFavoriteRepository) IsFavorite(userID, songID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.UserFavorite{}).
		Where("user_id = ? AND song_id = ?", userID, songID).
		Count(&count).
		Error
	return count > 0, err
}
