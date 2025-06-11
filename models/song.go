package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Song represents the song model
type Song struct {
	gorm.Model
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title        string         `gorm:"not null;index:idx_title,priority:1" json:"title"` // Indexed for search
	ArtistID     uuid.UUID      `gorm:"not null;index:idx_artist_id" json:"artist_id"`    // Foreign key to User (Artist)
	ArtistsNames []string       `gorm:"type:text[]" json:"artists_names"`                 // Array of artist names
	Genre        string         `gorm:"not null;index:idx_genre,priority:2" json:"genre"` // Indexed for genre search
	Price        float32        `gorm:"not null" json:"price"`
	Duration     int            `gorm:"not null" json:"duration"`           // Duration of the song in seconds
	AudioURL     string         `gorm:"type:varchar(255)" json:"audio_url"` // URL for the song's audio file
	ReleaseDate  time.Time      `gorm:"not null" json:"release_date"`
	AlbumID      uuid.UUID      `gorm:"index" json:"album_id,omitempty"`   // Foreign key to Album, nullable if it's a single
	IsPurchased  bool           `gorm:"default:false" json:"is_purchased"` // If the song is purchased
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
