package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongMetricsRepositoryImpl struct {
	DB *gorm.DB
}

func NewSongMetricsRepository(db *gorm.DB) SongMetricsRepository {
	return &SongMetricsRepositoryImpl{DB: db}
}

func (r *SongMetricsRepositoryImpl) CreateSongMetrics(songMetrics *models.SongMetrics) (*models.SongMetrics, error) {
	err := r.DB.Create(songMetrics).Error
	if err != nil {
		return nil, err
	}
	return songMetrics, nil
}

func (r *SongMetricsRepositoryImpl) GetMetricsBySongID(songID uuid.UUID) (*models.SongMetrics, error) {
	var metrics models.SongMetrics
	err := r.DB.First(&metrics, "song_id = ?", songID).Error
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}

func (r *SongMetricsRepositoryImpl) UpdateSongMetrics(songID uuid.UUID, metrics *models.SongMetrics) (*models.SongMetrics, error) {
	var existingMetrics models.SongMetrics
	err := r.DB.First(&existingMetrics, "song_id = ?", songID).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&existingMetrics).Updates(metrics).Error
	if err != nil {
		return nil, err
	}
	return &existingMetrics, nil
}
