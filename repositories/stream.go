package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type StreamRepository struct {
	BaseRepository[models.Stream]
}

func NewStreamRepository(db *gorm.DB) IStreamRepository {
	return &StreamRepository{
		BaseRepository: BaseRepository[models.Stream]{DB: db},
	}
}

func (r *StreamRepository) GetStreamCount(songID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Stream{}).
		Where("song_id = ? AND created_at >= ?", songID, since).
		Count(&count).
		Error
	return count, err
}

func (r *StreamRepository) GetArtistStreams(artistID uuid.UUID, start, end time.Time) ([]models.Stream, error) {
	var streams []models.Stream
	err := r.DB.
		Joins("JOIN songs ON streams.song_id = songs.id").
		Where("songs.artist_id = ? AND streams.created_at BETWEEN ? AND ?", artistID, start, end).
		Find(&streams).
		Error
	return streams, err
}

func (r *StreamRepository) GetStreamBySong(songID uuid.UUID) (*models.Stream, error) {
	var stream *models.Stream
	err := r.DB.Model(&models.Stream{}).
		Where("song_id = ?", songID).
		Find(&stream).
		Error
	return stream, err
}
