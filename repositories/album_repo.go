package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// AlbumRepository defines methods for interacting with the Album model
type AlbumRepository interface {
	CreateAlbum(album *models.Album) (*models.Album, error)
	GetAlbumByID(id uuid.UUID) (*models.Album, error)
	GetAlbumsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Album, int64, error)
	UpdateAlbum(id uuid.UUID, album *models.Album) (*models.Album, error)
	DeleteAlbum(id uuid.UUID) error
}
