package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongPurchaseRepository struct {
	BaseRepository[models.SongPurchase]
}

func NewSongPurchaseRepository(db *gorm.DB) ISongPurchaseRepository {
	return &SongPurchaseRepository{
		BaseRepository: BaseRepository[models.SongPurchase]{DB: db},
	}
}

func (r *SongPurchaseRepository) FindByUserAndSong(userID, songID uuid.UUID) (*models.SongPurchase, error) {
	var purchase models.SongPurchase
	err := r.DB.
		Where("user_id = ? AND song_id = ?", userID, songID).
		First(&purchase).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &purchase, err
}

func (r *SongPurchaseRepository) GetUserSongPurchases(userID uuid.UUID) ([]models.SongPurchase, error) {
	var purchases []models.SongPurchase
	err := r.DB.
		Preload("Song").
		Preload("Song.Artist").
		Where("user_id = ?", userID).
		Find(&purchases).
		Error
	return purchases, err
}

func (r *SongPurchaseRepository) HasPurchasedSong(userID, songID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.SongPurchase{}).
		Where("user_id = ? AND song_id = ? payment_status = completed", userID, songID).
		Count(&count).
		Error
	return count > 0, err
}

func (r *SongPurchaseRepository) GetUserPurchases(userID uuid.UUID) ([]models.SongPurchase, error) {
	var purchases []models.SongPurchase
	err := r.DB.
		Where("user_id = ?", userID).
		Preload("Song").
		Preload("Song.Artist").
		Find(&purchases).
		Error
	return purchases, err
}
