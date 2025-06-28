package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type ModerationService interface {
	FlagContent(ctx context.Context, flag *models.ContentFlag) (*models.ContentFlag, error)
	GetFlaggedContent(ctx context.Context) ([]models.ContentFlag, error)
	ReviewFlag(ctx context.Context, flagID uuid.UUID, status string) error
	GetFlagByID(ctx context.Context, flagID uuid.UUID) (*models.ContentFlag, error)
}

type moderationService struct {
	moderationRepo repositories.IModerationRepository
}

func NewModerationService(moderationRepo repositories.IModerationRepository) ModerationService {
	return &moderationService{
		moderationRepo: moderationRepo,
	}
}

func (s *moderationService) FlagContent(ctx context.Context, flag *models.ContentFlag) (*models.ContentFlag, error) {
	// Validate required fields
	if flag.ReporterUserID == uuid.Nil {
		return nil, errors.New("reporter user ID is required")
	}
	if flag.TargetID == uuid.Nil {
		return nil, errors.New("target ID is required")
	}
	if flag.TargetType == "" {
		return nil, errors.New("target type is required")
	}
	if flag.Reason == "" {
		return nil, errors.New("reason is required")
	}

	// Set default status if not provided
	if flag.Status == "" {
		flag.Status = "pending"
	}

	// Create the flag
	return s.moderationRepo.Create(flag)
}

func (s *moderationService) GetFlaggedContent(ctx context.Context) ([]models.ContentFlag, error) {
	return s.moderationRepo.GetFlaggedContent()
}

func (s *moderationService) ReviewFlag(ctx context.Context, flagID uuid.UUID, status string) error {
	// Validate status
	if status != "approved" && status != "rejected" && status != "pending" {
		return errors.New("invalid status, must be 'approved', 'rejected', or 'pending'")
	}

	// Check if flag exists
	_, err := s.moderationRepo.GetByID(flagID)
	if err != nil {
		return errors.New("flag not found")
	}

	return s.moderationRepo.UpdateFlagStatus(flagID, status)
}

func (s *moderationService) GetFlagByID(ctx context.Context, flagID uuid.UUID) (*models.ContentFlag, error) {
	return s.moderationRepo.GetByID(flagID)
}
