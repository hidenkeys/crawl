package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"fmt"
	"github.com/google/uuid"
)

type SongService struct {
	songRepo repositories.SongRepository
}

func NewSongService(songRepo repositories.SongRepository) *SongService {
	return &SongService{songRepo: songRepo}
}

// CreateSong creates a new song
func (s *SongService) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	fmt.Println(song.ArtistsNames)
	fmt.Println("in service CreateSong")
	return s.songRepo.CreateSong(song)
}

// GetSongByID fetches a song by ID
func (s *SongService) GetSongByID(ctx context.Context, id uuid.UUID) (*models.Song, error) {
	return s.songRepo.GetSongByID(id)
}

// GetSongsByArtist fetches songs by artist ID
func (s *SongService) GetSongsByArtist(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]models.Song, int64, error) {
	return s.songRepo.GetSongsByArtist(artistID, limit, offset)
}

// UpdateSong updates a song by ID
func (s *SongService) UpdateSong(ctx context.Context, id uuid.UUID, song *models.Song) (*models.Song, error) {
	return s.songRepo.UpdateSong(id, song)
}

// DeleteSong deletes a song by ID
func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) error {
	return s.songRepo.DeleteSong(id)
}
