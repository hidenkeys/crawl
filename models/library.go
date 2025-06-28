package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SongPurchase struct {
	BaseModel
	UserID              uuid.UUID `gorm:"not null;index:idx_user_song,unique" json:"user_id"`
	SongID              uuid.UUID `gorm:"not null;index:idx_user_song,unique" json:"song_id"`
	PurchasePrice       float64   `gorm:"type:decimal(10,2);not null" json:"purchase_price"`
	Currency            string    `gorm:"size:3;default:'USD'" json:"currency"`
	PaymentStatus       string    `gorm:"size:20;default:'pending'" json:"payment_status"`
	StripeTransactionID string    `gorm:"size:255" json:"stripe_transaction_id"`
	User                User      `gorm:"foreignKey:UserID" json:"-"`
	Song                Song      `gorm:"foreignKey:SongID" json:"song"`
}

type AlbumPurchase struct {
	BaseModel
	UserID              uuid.UUID `gorm:"not null;index:idx_user_album,unique" json:"user_id"`
	AlbumID             uuid.UUID `gorm:"not null;index:idx_user_album,unique" json:"album_id"`
	PurchasePrice       float64   `gorm:"type:decimal(10,2);not null" json:"purchase_price"`
	Currency            string    `gorm:"size:3;default:'USD'" json:"currency"`
	PaymentStatus       string    `gorm:"size:20;default:'pending'" json:"payment_status"`
	StripeTransactionID string    `gorm:"size:255" json:"stripe_transaction_id"`
	User                User      `gorm:"foreignKey:UserID" json:"-"`
	Album               Album     `gorm:"foreignKey:AlbumID" json:"album"`
}

type UserFavorite struct {
	UserID    uuid.UUID      `gorm:"primaryKey" json:"user_id"`
	SongID    uuid.UUID      `gorm:"primaryKey" json:"song_id"`
	CreatedAt time.Time      `json:"created_at"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Song      Song           `gorm:"foreignKey:SongID" json:"song"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type Playlist struct {
	BaseModel
	UserID        uuid.UUID `gorm:"not null;index" json:"user_id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	CoverImageURL string    `gorm:"size:255" json:"cover_image_url"`
	IsPublic      bool      `gorm:"default:false" json:"is_public"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	Likes         *int      `gorm:"default:0" json:"likes"`
	Songs         []Song    `gorm:"many2many:playlist_songs;" json:"songs,omitempty"`
}

type PlaylistSong struct {
	PlaylistID uuid.UUID      `gorm:"primaryKey" json:"playlist_id"`
	SongID     uuid.UUID      `gorm:"primaryKey" json:"song_id"`
	Position   int            `gorm:"not null" json:"position"`
	AddedAt    time.Time      `json:"added_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
