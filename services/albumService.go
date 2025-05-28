package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"github.com/google/uuid"
)

type AlbumService struct {
	albumRepo repositories.AlbumRepository
}

func NewAlbumService(albumRepo repositories.AlbumRepository) *AlbumService {
	return &AlbumService{albumRepo: albumRepo}
}

// CreateAlbum creates a new album
func (s *AlbumService) CreateAlbum(ctx context.Context, album *models.Album) (*models.Album, error) {
	return s.albumRepo.CreateAlbum(album)
}

// GetAlbumByID fetches a specific album by ID
func (s *AlbumService) GetAlbumByID(ctx context.Context, id uuid.UUID) (*models.Album, error) {
	return s.albumRepo.GetAlbumByID(id)
}

// GetAlbumsByArtist fetches albums by artist ID
func (s *AlbumService) GetAlbumsByArtist(ctx context.Context, artistID uuid.UUID, limit, offset int) ([]models.Album, int64, error) {
	return s.albumRepo.GetAlbumsByArtist(artistID, limit, offset)
}

// UpdateAlbum updates an album by ID
func (s *AlbumService) UpdateAlbum(ctx context.Context, id uuid.UUID, album *models.Album) (*models.Album, error) {
	return s.albumRepo.UpdateAlbum(id, album)
}

// DeleteAlbum deletes an album by ID
func (s *AlbumService) DeleteAlbum(ctx context.Context, id uuid.UUID) error {
	return s.albumRepo.DeleteAlbum(id)
}
