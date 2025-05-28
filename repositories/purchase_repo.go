package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// PurchaseRepository defines methods for interacting with the Purchase model
type PurchaseRepository interface {
	CreatePurchase(purchase *models.Purchase) (*models.Purchase, error)
	GetPurchasesByUser(userID uuid.UUID, limit, offset int) ([]models.Purchase, int64, error)
}
