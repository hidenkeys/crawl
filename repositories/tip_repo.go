package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// TipRepository defines methods for interacting with the Tip model
type TipRepository interface {
	CreateTip(tip *models.Tip) (*models.Tip, error)
	GetTipsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Tip, int64, error)
	GetTipsByUser(userID uuid.UUID, limit, offset int) ([]models.Tip, int64, error)
}
