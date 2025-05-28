package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
)

// SongRepository defines methods for interacting with the Song model
type SongRepository interface {
	CreateSong(song *models.Song) (*models.Song, error)
	GetSongByID(id uuid.UUID) (*models.Song, error)
	GetSongsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Song, int64, error)
	UpdateSong(id uuid.UUID, song *models.Song) (*models.Song, error)
	DeleteSong(id uuid.UUID) error
}
