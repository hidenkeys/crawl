package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseRepositoryImpl struct {
	DB *gorm.DB
}

func NewPurchaseRepository(db *gorm.DB) PurchaseRepository {
	return &PurchaseRepositoryImpl{DB: db}
}

func (r *PurchaseRepositoryImpl) CreatePurchase(purchase *models.Purchase) (*models.Purchase, error) {
	err := r.DB.Create(purchase).Error
	if err != nil {
		return nil, err
	}
	return purchase, nil
}

func (r *PurchaseRepositoryImpl) GetPurchasesByUser(userID uuid.UUID, limit, offset int) ([]models.Purchase, int64, error) {
	var purchases []models.Purchase
	var count int64
	err := r.DB.Model(&models.Purchase{}).Where("user_id = ?", userID).Count(&count).Offset(offset).Limit(limit).Find(&purchases).Error
	if err != nil {
		return nil, count, err
	}
	return purchases, count, nil
}
