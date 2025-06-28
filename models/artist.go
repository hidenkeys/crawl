package models

import "github.com/google/uuid"

type Artist struct {
	BaseModel
	UserID           uuid.UUID `gorm:"uniqueIndex;not null" json:"user_id"`
	ArtistName       string    `gorm:"size:100;not null" json:"artist_name"`
	Verified         bool      `gorm:"default:false" json:"verified"`
	WalletBalance    float64   `gorm:"type:decimal(10,2);default:0.00" json:"wallet_balance"`
	StripeAccountID  string    `gorm:"size:255" json:"-"`
	MonthlyListeners int       `gorm:"default:0" json:"monthly_listeners"`
	User             User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Songs            []Song    `gorm:"foreignKey:ArtistID" json:"songs,omitempty"`
	Albums           []Album   `gorm:"foreignKey:ArtistID" json:"albums,omitempty"`
}
