package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type SongRepository struct {
	BaseRepository[models.Song]
}

func NewSongRepository(db *gorm.DB) ISongRepository {
	return &SongRepository{
		BaseRepository: BaseRepository[models.Song]{DB: db},
	}
}

func (r *SongRepository) GetWithArtist(id uuid.UUID) (*models.Song, error) {
	var song models.Song
	err := r.DB.Preload("Artist").First(&song, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &song, err
}

func (r *SongRepository) GetByArtist(artistID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	err := r.DB.Where("artist_id = ?", artistID).Find(&songs).Error
	return songs, err
}

func (r *SongRepository) GetTrending(limit int, since time.Time) ([]models.Song, error) {
	var songs []models.Song
	err := r.DB.
		Where("created_at >= ?", since).
		Order("plays_count DESC").
		Limit(limit).
		Find(&songs).
		Error
	return songs, err
}

func (r *SongRepository) AddPlayCount(id uuid.UUID, count int) error {
	return r.DB.Model(&models.Song{}).
		Where("id = ?", id).
		Update("plays_count", gorm.Expr("plays_count + ?", count)).
		Error
}

func (r *SongRepository) GetByAlbum(albumID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	err := r.DB.Where("album_id = ?", albumID).Find(&songs).Error
	return songs, err
}

func (r *SongRepository) Search(query, artist, genre *string, sort, order *string, offset, limit int) ([]models.Song, error) {
	var songs []models.Song
	db := r.DB.Model(&models.Song{}).Preload("Artist")

	if query != nil && *query != "" {
		db = db.Where("title ILIKE ?", "%"+*query+"%")
	}

	if artist != nil && *artist != "" {
		db = db.Joins("JOIN artists ON songs.artist_id = artists.id").
			Where("artists.name ILIKE ?", "%"+*artist+"%")
	}

	if genre != nil && *genre != "" {
		db = db.Joins("JOIN genres ON songs.genre_id = genres.id").
			Where("genres.name ILIKE ?", "%"+*genre+"%")
	}

	if sort != nil {
		sortField := *sort
		if order != nil && *order == "desc" {
			sortField += " DESC"
		}
		db = db.Order(sortField)
	}

	if limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&songs).Error
	return songs, err
}
