package services

import (
	"context"
	"crawl/api"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type TipService interface {
	SendTip(ctx context.Context, tip api.PostTipsJSONBody, senderID uuid.UUID) (*models.ArtistTip, error)
	GetArtistTips(ctx context.Context, artistID uuid.UUID, limit int) ([]models.ArtistTip, error)
	GetUserTips(ctx context.Context, userID uuid.UUID, limit int) ([]models.ArtistTip, error)
}
type tipService struct {
	tipRepo    repositories.ITipRepository
	userRepo   repositories.IUserRepository
	artistRepo repositories.IArtistRepository
}

func NewTipService(
	tipRepo repositories.ITipRepository,
	userRepo repositories.IUserRepository,
	artistRepo repositories.IArtistRepository,
) TipService {
	return &tipService{
		tipRepo:    tipRepo,
		userRepo:   userRepo,
		artistRepo: artistRepo,
	}
}

func (s *tipService) SendTip(ctx context.Context, tip api.PostTipsJSONBody, senderID uuid.UUID) (*models.ArtistTip, error) {
	// Validate required fields
	if senderID == uuid.Nil {
		return nil, errors.New("sender ID is required")
	}
	if tip.ArtistId == uuid.Nil {
		return nil, errors.New("artist ID is required")
	}
	if tip.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	// Verify sender exists
	_, err := s.userRepo.GetByID(senderID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("sender not found")
		}
		return nil, err
	}

	// Verify artist exists
	_, err = s.artistRepo.GetByID(tip.ArtistId)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}

	// Create new tip
	newTip := &models.ArtistTip{
		SenderID:            senderID,
		ArtistID:            tip.ArtistId,
		Amount:              tip.Amount,
		Message:             *tip.Message,
		Currency:            "NGN",
		PaymentStatus:       "completed",
		StripeTransactionID: tip.PaymentMethodId,
	}

	return s.tipRepo.Create(newTip)
}

func (s *tipService) GetArtistTips(ctx context.Context, artistID uuid.UUID, limit int) ([]models.ArtistTip, error) {
	// Verify artist exists
	_, err := s.artistRepo.GetByID(artistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}

	return s.tipRepo.GetArtistTips(artistID)
}

func (s *tipService) GetUserTips(ctx context.Context, userID uuid.UUID, limit int) ([]models.ArtistTip, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return s.tipRepo.GetUserTips(userID)
}

func (s *tipService) GetTotalTipsReceived(ctx context.Context, artistID uuid.UUID) (int64, error) {
	// Verify artist exists
	_, err := s.artistRepo.GetByID(artistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return 0, errors.New("artist not found")
		}
		return 0, err
	}

	return s.tipRepo.GetTotalTipsReceived(artistID)
}

func (s *tipService) GetTotalTipsSent(ctx context.Context, userID uuid.UUID) (int64, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}

	return s.tipRepo.GetTotalTipsSent(userID)
}
