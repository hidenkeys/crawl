package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Album represents the album model
type Album struct {
	gorm.Model
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string         `gorm:"not null;index:idx_album_title,priority:1" json:"title"` // Indexed for search
	ArtistID    uuid.UUID      `gorm:"not null;index:idx_artist_id" json:"artist_id"`          // Foreign key to User (Artist)
	ArtistName  string         `gorm:"type:text[]" json:"artists_name"`                        // Array of artist names
	Price       float64        `gorm:"not null" json:"price"`
	ReleaseDate time.Time      `gorm:"not null" json:"release_date"`
	CoverArt    string         `gorm:"type:varchar(255)" json:"cover_art"` // URL for album cover art
	IsAvailable bool           `gorm:"default:true" json:"is_available"`   // If the album is available for purchase
	Songs       []Song         `gorm:"foreignKey:AlbumID" json:"songs"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
