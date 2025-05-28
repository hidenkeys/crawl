package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"github.com/google/uuid"
)

type TipService struct {
	tipRepo repositories.TipRepository
}

func NewTipService(tipRepo repositories.TipRepository) *TipService {
	return &TipService{tipRepo: tipRepo}
}

// CreateTip creates a new tip record
func (s *TipService) CreateTip(ctx context.Context, tip *models.Tip) (*models.Tip, error) {
	return s.tipRepo.CreateTip(tip)
}

// GetTipsByArtist retrieves all tips received by an artist
func (s *TipService) GetTipsByArtist(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]models.Tip, int64, error) {
	return s.tipRepo.GetTipsByArtist(artistID, limit, offset)
}

// GetTipsByUser retrieves all tips given by a user
func (s *TipService) GetTipsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Tip, int64, error) {
	return s.tipRepo.GetTipsByUser(userID, limit, offset)
}
