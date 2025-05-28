package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// User represents the user model
type User struct {
	gorm.Model
	ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
	Username           string         `gorm:"unique;not null" json:"username"`
	Email              string         `gorm:"unique;not null" json:"email"`
	Password           string         `gorm:"not null" json:"password"`
	Role               string         `gorm:"type:varchar(10);default:'listener';not null" json:"role"` // 'artist' or 'listener'
	RecentlyPlayed     []uuid.UUID    `gorm:"type:uuid[]" json:"recently_played"`                       // Array of UUIDs for recently played songs
	ProfilePicture     string         `json:"profile_picture,omitempty"`
	TippedArtists      []uuid.UUID    `gorm:"type:uuid[]" json:"tipped_artists"`                          // Array of artists tipped
	SubscriptionStatus string         `gorm:"type:varchar(20);default:'free'" json:"subscription_status"` // 'free', 'premium', etc.
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft delete
}
