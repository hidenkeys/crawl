package models

import (
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"gorm.io/gorm"
	"time"
)

type Genre struct {
	BaseModel
	Name        string  `gorm:"size:50;uniqueIndex" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	ImageURL    string  `gorm:"size:255" json:"image_url"`
	Songs       []Song  `gorm:"foreignKey:GenreID" json:"songs,omitempty"`
	Albums      []Album `gorm:"foreignKey:GenreID" json:"albums,omitempty"`
}

type Song struct {
	BaseModel
	Title         string             `gorm:"size:255;not null" json:"title"`
	ArtistID      uuid.UUID          `gorm:"not null;index" json:"artist_id"`
	AlbumID       *uuid.UUID         `gorm:"index" json:"album_id,omitempty"`
	Duration      int                `gorm:"not null" json:"duration"` // in seconds
	Price         int                `gorm:"not null" json:"price"`
	AudioURL      string             `gorm:"size:255;not null" json:"audio_url"`
	PreviewURL    string             `gorm:"size:255" json:"preview_url"`
	ReleaseDate   openapi_types.Date `gorm:"type:date" json:"release_date"`
	CoverImageURL string             `gorm:"size:255" json:"cover_image_url"`
	GenreID       *uuid.UUID         `gorm:"index" json:"genre_id,omitempty"`
	PlaysCount    int                `gorm:"default:0" json:"plays_count"`
	Likes         *int               `gorm:"default:0" json:"likes"`
	IsFlagged     bool               `gorm:"default:false" json:"is_flagged"`
	Artist        Artist             `gorm:"foreignKey:ArtistID" json:"artist"`
	Album         *Album             `gorm:"foreignKey:AlbumID" json:"album,omitempty"`
	Genre         *Genre             `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	Contributors  []Artist           `gorm:"many2many:song_contributors;" json:"contributors,omitempty"`
}

type Album struct {
	BaseModel
	Title         string             `gorm:"size:255;not null" json:"title"`
	ArtistID      uuid.UUID          `gorm:"not null;index" json:"artist_id"`
	Description   string             `gorm:"type:text" json:"description"`
	Price         int                `gorm:"not null" json:"price"`
	CoverImageURL string             `gorm:"size:255" json:"cover_image_url"`
	ReleaseDate   openapi_types.Date `gorm:"type:date" json:"release_date"`
	GenreID       *uuid.UUID         `gorm:"index" json:"genre_id,omitempty"`
	IsFlagged     bool               `gorm:"default:false" json:"is_flagged"`
	Artist        Artist             `gorm:"foreignKey:ArtistID" json:"artist"`
	Genre         *Genre             `gorm:"foreignKey:GenreID" json:"genre,omitempty"`
	Songs         []Song             `gorm:"foreignKey:AlbumID" json:"songs,omitempty"`
	Contributors  []Artist           `gorm:"many2many:album_contributors;" json:"contributors,omitempty"`
}

type SongContributor struct {
	SongID            uuid.UUID      `gorm:"primaryKey" json:"song_id"`
	ArtistID          uuid.UUID      `gorm:"primaryKey" json:"artist_id"`
	ContributionType  string         `gorm:"size:100;primaryKey" json:"contribution_type"`
	RoyaltyPercentage int            `gorm:"not null;default:0" json:"royalty_percentage"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type AlbumContributor struct {
	AlbumID           uuid.UUID      `gorm:"primaryKey" json:"album_id"`
	ArtistID          uuid.UUID      `gorm:"primaryKey" json:"artist_id"`
	ContributionType  string         `gorm:"size:100;primaryKey" json:"contribution_type"`
	RoyaltyPercentage int            `gorm:"not null;default:0" json:"royalty_percentage"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
