package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TipRepositoryImpl struct {
	DB *gorm.DB
}

func NewTipRepository(db *gorm.DB) TipRepository {
	return &TipRepositoryImpl{DB: db}
}

func (r *TipRepositoryImpl) CreateTip(tip *models.Tip) (*models.Tip, error) {
	err := r.DB.Create(tip).Error
	if err != nil {
		return nil, err
	}
	return tip, nil
}

func (r *TipRepositoryImpl) GetTipsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Tip, int64, error) {
	var tips []models.Tip
	var count int64
	err := r.DB.Model(&models.Tip{}).Where("artist_id = ?", artistID).Count(&count).Offset(offset).Limit(limit).Find(&tips).Error
	if err != nil {
		return nil, count, err
	}
	return tips, count, nil
}

func (r *TipRepositoryImpl) GetTipsByUser(userID uuid.UUID, limit, offset int) ([]models.Tip, int64, error) {
	var tips []models.Tip
	var count int64
	err := r.DB.Model(&models.Tip{}).Where("user_id = ?", userID).Count(&count).Offset(offset).Limit(limit).Find(&tips).Error
	if err != nil {
		return nil, count, err
	}
	return tips, count, nil
}
