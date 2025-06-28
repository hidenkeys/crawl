package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumPurchaseRepository struct {
	BaseRepository[models.AlbumPurchase]
}

func NewAlbumPurchaseRepository(db *gorm.DB) IAlbumPurchaseRepository {
	return &AlbumPurchaseRepository{
		BaseRepository: BaseRepository[models.AlbumPurchase]{DB: db},
	}
}

func (r *AlbumPurchaseRepository) FindByUserAndAlbum(userID, albumID uuid.UUID) (*models.AlbumPurchase, error) {
	var purchase models.AlbumPurchase
	err := r.DB.
		Where("user_id = ? AND album_id = ?", userID, albumID).
		First(&purchase).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &purchase, err
}

func (r *AlbumPurchaseRepository) GetUserAlbumPurchases(userID uuid.UUID) ([]models.AlbumPurchase, error) {
	var purchases []models.AlbumPurchase
	err := r.DB.
		Where("user_id = ?", userID).
		Preload("Album").
		Preload("Album.Artist").
		Find(&purchases).
		Error
	return purchases, err
}

func (r *AlbumPurchaseRepository) HasPurchasedAlbum(userID, albumID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(&models.AlbumPurchase{}).
		Where("user_id = ? AND album_id = ? payment_status = completed", userID, albumID).
		Count(&count).
		Error
	return count > 0, err
}
