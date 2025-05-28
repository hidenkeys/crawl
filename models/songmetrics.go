package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SongMetrics tracks the performance of each song
type SongMetrics struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SongID      uuid.UUID `gorm:"type:uuid" json:"song_id"`
	ListenCount int       `gorm:"default:0" json:"listen_count"` // Number of listens
	PlayCount   int       `gorm:"default:0" json:"play_count"`   // Number of plays
}
