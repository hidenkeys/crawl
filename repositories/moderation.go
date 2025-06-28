package repositories

import (
	"crawl/models"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ModerationRepository struct {
	BaseRepository[models.ContentFlag]
}

func NewModerationRepository(db *gorm.DB) IModerationRepository {
	return &ModerationRepository{
		BaseRepository: BaseRepository[models.ContentFlag]{DB: db},
	}
}

func (r *ModerationRepository) GetFlaggedContent() ([]models.ContentFlag, error) {
	var flags []models.ContentFlag
	err := r.DB.
		Where("status = 'pending'").
		Preload("Reporter").
		Find(&flags).
		Error
	return flags, err
}

func (r *ModerationRepository) UpdateFlagStatus(id uuid.UUID, status string) error {
	return r.DB.Model(&models.ContentFlag{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}
