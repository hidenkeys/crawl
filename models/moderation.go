package models

import "github.com/google/uuid"

type ContentFlag struct {
	BaseModel
	ReporterUserID uuid.UUID `gorm:"not null" json:"reporter_user_id"`    // user Id
	TargetID       uuid.UUID `gorm:"not null" json:"target_id"`           // song or album
	TargetType     string    `gorm:"size:20;not null" json:"target_type"` // "song" or "album"
	Reason         string    `gorm:"size:100;not null" json:"reason"`
	Description    string    `gorm:"type:text" json:"description"`
	Status         string    `gorm:"size:20;default:'pending'" json:"status"`
	Reporter       User      `gorm:"foreignKey:ReporterUserID" json:"reporter"`
}
