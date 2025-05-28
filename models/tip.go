package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tip represents a tip given by a listener to an artist
type Tip struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID   uuid.UUID `gorm:"not null" json:"user_id"`   // User who gave the tip
	ArtistID uuid.UUID `gorm:"not null" json:"artist_id"` // Artist who received the tip
	Amount   float64   `gorm:"not null" json:"amount"`    // Amount tipped
	Message  string    `json:"message,omitempty"`         // Optional message with the tip
}
