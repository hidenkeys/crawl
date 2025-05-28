package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"github.com/google/uuid"
)

type PurchaseService struct {
	purchaseRepo repositories.PurchaseRepository
}

func NewPurchaseService(purchaseRepo repositories.PurchaseRepository) *PurchaseService {
	return &PurchaseService{purchaseRepo: purchaseRepo}
}

// CreatePurchase creates a new purchase record
func (s *PurchaseService) CreatePurchase(ctx context.Context, purchase *models.Purchase) (*models.Purchase, error) {
	return s.purchaseRepo.CreatePurchase(purchase)
}

// GetPurchasesByUser retrieves all purchases by a user
func (s *PurchaseService) GetPurchasesByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Purchase, int64, error) {
	return s.purchaseRepo.GetPurchasesByUser(userID, limit, offset)
}
