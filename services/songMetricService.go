package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"github.com/google/uuid"
)

type SongMetricsService struct {
	songMetricsRepo repositories.SongMetricsRepository
}

func NewSongMetricsService(songMetricsRepo repositories.SongMetricsRepository) *SongMetricsService {
	return &SongMetricsService{songMetricsRepo: songMetricsRepo}
}

// CreateSongMetrics creates a new song metrics record
func (s *SongMetricsService) CreateSongMetrics(ctx context.Context, songMetrics *models.SongMetrics) (*models.SongMetrics, error) {
	return s.songMetricsRepo.CreateSongMetrics(songMetrics)
}

// GetMetricsBySongID retrieves song metrics by song ID
func (s *SongMetricsService) GetMetricsBySongID(ctx context.Context, songID uuid.UUID) (*models.SongMetrics, error) {
	return s.songMetricsRepo.GetMetricsBySongID(songID)
}

// UpdateSongMetrics updates the song metrics by song ID
func (s *SongMetricsService) UpdateSongMetrics(ctx context.Context, songID uuid.UUID, metrics *models.SongMetrics) (*models.SongMetrics, error) {
	return s.songMetricsRepo.UpdateSongMetrics(songID, metrics)
}
