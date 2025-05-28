package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Purchase represents the purchase model
type Purchase struct {
	gorm.Model
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID      `gorm:"not null" json:"user_id"`         // Foreign key to User
	SongID     *uuid.UUID     `gorm:"index" json:"song_id,omitempty"`  // Foreign key to Song, nullable if it's an album purchase
	AlbumID    *uuid.UUID     `gorm:"index" json:"album_id,omitempty"` // Foreign key to Album, nullable if it's a song purchase
	TotalPrice float64        `gorm:"not null" json:"total_price"`     // Total price for the purchase
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
