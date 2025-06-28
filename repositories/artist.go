package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type ArtistRepository struct {
	BaseRepository[models.Artist]
}

func NewArtistRepository(db *gorm.DB) IArtistRepository {
	return &ArtistRepository{
		BaseRepository: BaseRepository[models.Artist]{DB: db},
	}
}

func (r *ArtistRepository) GetWithSongs(id uuid.UUID) (*models.Artist, error) {
	var artist models.Artist
	err := r.DB.Preload("Songs").First(&artist, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &artist, err
}

func (r *ArtistRepository) GetWithAlbums(id uuid.UUID) (*models.Artist, error) {
	var artist models.Artist
	err := r.DB.Preload("Albums").First(&artist, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &artist, err
}

func (r *ArtistRepository) GetWithUserId(userID uuid.UUID) (*models.Artist, error) {
	var artist models.Artist
	err := r.DB.Where("user_id = ?", userID).Find(&artist).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &artist, err
}

func (r *ArtistRepository) SearchByName(query string, limit int, offset int) ([]models.Artist, error) {
	var artists []models.Artist
	err := r.DB.
		Where("artist_name ILIKE ?", "%"+query+"%").
		Limit(limit).
		Offset(offset).
		Find(&artists).
		Error
	return artists, err
}
