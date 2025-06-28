package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TipRepository struct {
	BaseRepository[models.ArtistTip]
}

func NewTipRepository(db *gorm.DB) ITipRepository {
	return &TipRepository{
		BaseRepository: BaseRepository[models.ArtistTip]{DB: db},
	}
}

func (r *TipRepository) GetArtistTips(artistID uuid.UUID) ([]models.ArtistTip, error) {
	var tips []models.ArtistTip
	err := r.DB.
		Where("artist_id = ?", artistID).
		Order("created_at DESC").
		Find(&tips).Error
	return tips, err
}

func (r *TipRepository) GetUserTips(userID uuid.UUID) ([]models.ArtistTip, error) {
	var tips []models.ArtistTip
	err := r.DB.
		Where("sender_id = ?", userID).
		Order("created_at DESC").
		Find(&tips).Error
	return tips, err
}

func (r *TipRepository) GetTotalTipsReceived(artistID uuid.UUID) (int64, error) {
	var total int64
	err := r.DB.Model(&models.ArtistTip{}).
		Where("artist_id = ?", artistID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *TipRepository) GetTotalTipsSent(userID uuid.UUID) (int64, error) {
	var total int64
	err := r.DB.Model(&models.ArtistTip{}).
		Where("sender_id = ?", userID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
