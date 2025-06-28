package models

import "github.com/google/uuid"

type ArtistTip struct {
	BaseModel
	SenderID            uuid.UUID `gorm:"not null" json:"sender_id"`
	ArtistID            uuid.UUID `gorm:"not null" json:"artist_id"`
	Amount              int       `gorm:"not null;default:0" json:"amount"`
	Message             string    `gorm:"type:text" json:"message"`
	Currency            string    `gorm:"size:3;default:'NGN'" json:"currency"`
	PaymentStatus       string    `gorm:"size:20;default:'completed'" json:"payment_status"`
	StripeTransactionID string    `gorm:"size:255" json:"stripe_transaction_id"`
	Sender              User      `gorm:"foreignKey:SenderID" json:"sender"`
	Artist              Artist    `gorm:"foreignKey:ArtistID" json:"artist"`
}

type Stream struct {
	BaseModel
	UserID      *uuid.UUID `gorm:"index" json:"user_id,omitempty"`
	SongID      uuid.UUID  `gorm:"not null;index" json:"song_id"`
	IsPreview   bool       `gorm:"default:false" json:"is_preview"`
	DeviceType  string     `gorm:"size:50" json:"device_type"`
	CountryCode string     `gorm:"size:2" json:"country_code"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Song        Song       `gorm:"foreignKey:SongID" json:"song"`
}

type MonthlyRoyalty struct {
	BaseModel
	ArtistID   uuid.UUID `gorm:"not null;index:idx_artist_month,unique" json:"artist_id"`
	Year       int       `gorm:"not null;index:idx_artist_month,unique" json:"year"`
	Month      int       `gorm:"not null;index:idx_artist_month,unique" json:"month"`
	Amount     int       `gorm:"not null;default:0" json:"amount"`
	Currency   string    `gorm:"size:3;default:'NGN'" json:"currency"`
	PaidStatus bool      `gorm:"default:false" json:"paid_status"`
	Artist     Artist    `gorm:"foreignKey:ArtistID" json:"artist"`
}
