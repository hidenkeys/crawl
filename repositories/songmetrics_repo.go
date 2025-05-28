package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// SongMetricsRepository defines methods for interacting with the SongMetrics model
type SongMetricsRepository interface {
	CreateSongMetrics(songMetrics *models.SongMetrics) (*models.SongMetrics, error)
	GetMetricsBySongID(songID uuid.UUID) (*models.SongMetrics, error)
	UpdateSongMetrics(songID uuid.UUID, metrics *models.SongMetrics) (*models.SongMetrics, error)
}
