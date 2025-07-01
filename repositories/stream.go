package repositories

import (
	"context"
	"crawl/models"
	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
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

func (r *StreamRepository) GetUserRecentStreams(ctx context.Context, userID types.UUID, limit int, offset int) ([]models.Song, error) {
	var songs []models.Song

	// Subquery to get distinct song IDs from user's streams, ordered by most recent
	subQuery := r.DB.WithContext(ctx).
		Model(&models.Stream{}).
		Where("user_id = ?", userID).
		Select("DISTINCT ON (song_id) song_id, created_at").
		Order("song_id, created_at DESC")

	// Main query to get songs joined with the subquery
	err := r.DB.WithContext(ctx).
		Model(&models.Song{}).
		Select("songs.*").
		Joins("INNER JOIN (?) AS user_recent_streams ON songs.id = user_recent_streams.song_id", subQuery).
		Order("user_recent_streams.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&songs).Error

	if err != nil {
		return nil, err
	}

	return songs, nil
}
