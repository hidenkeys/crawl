package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type StreamService interface {
	RecordStream(ctx context.Context, stream models.Stream) error
	GetStreamCount(ctx context.Context, songID uuid.UUID) (int64, error)
	GetArtistStreams(ctx context.Context, artistID uuid.UUID) ([]models.Stream, error)
	GetStreamBySong(ctx context.Context, songID uuid.UUID) (*models.Stream, error)
}

type streamService struct {
	streamRepo repositories.IStreamRepository
	songRepo   repositories.ISongRepository
}

func NewStreamService(
	streamRepo repositories.IStreamRepository,
	songRepo repositories.ISongRepository,
) StreamService {
	return &streamService{
		streamRepo: streamRepo,
		songRepo:   songRepo,
	}
}

func (s *streamService) RecordStream(ctx context.Context, stream models.Stream) error {
	// Validate required fields
	if stream.SongID == uuid.Nil {
		return errors.New("song ID is required")
	}

	// Verify song exists
	_, err := s.songRepo.GetByID(stream.SongID)
	if err != nil {
		return errors.New("song not found")
	}

	// Set timestamp if not provided
	if stream.CreatedAt.IsZero() {
		stream.CreatedAt = time.Now()
	}

	// Record the stream
	_, err = s.streamRepo.Create(&stream)
	if err != nil {
		return err
	}

	// Increment play count in songs table
	err = s.songRepo.AddPlayCount(stream.SongID, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *streamService) GetStreamCount(ctx context.Context, songID uuid.UUID) (int64, error) {
	// Get count from last 30 days
	since := time.Now().AddDate(0, 0, -30)
	return s.streamRepo.GetStreamCount(songID, since)
}

func (s *streamService) GetArtistStreams(ctx context.Context, artistID uuid.UUID) ([]models.Stream, error) {
	// Get streams from last 30 days
	end := time.Now()
	start := end.AddDate(0, 0, -30)
	return s.streamRepo.GetArtistStreams(artistID, start, end)
}

func (s *streamService) GetStreamBySong(ctx context.Context, songID uuid.UUID) (*models.Stream, error) {
	// Verify song exists first
	_, err := s.songRepo.GetByID(songID)
	if err != nil {
		return nil, errors.New("song not found")
	}

	return s.streamRepo.GetStreamBySong(songID)
}
