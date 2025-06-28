package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type AlbumRepository struct {
	BaseRepository[models.Album]
}

func NewAlbumRepository(db *gorm.DB) IAlbumRepository {
	return &AlbumRepository{
		BaseRepository: BaseRepository[models.Album]{DB: db},
	}
}

func (r *AlbumRepository) GetWithSongs(id uuid.UUID) (*models.Album, error) {
	var album models.Album
	err := r.DB.Preload("Songs").First(&album, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &album, err
}

func (r *AlbumRepository) GetByArtist(artistID uuid.UUID) ([]models.Album, error) {
	var albums []models.Album
	err := r.DB.Where("artist_id = ?", artistID).Find(&albums).Error
	return albums, err
}

func (r *AlbumRepository) SearchAlbums(query *string, artist *string, genre *string, sort *string, page *int, limit *int) ([]models.Album, error) {
	// Implementation would depend on your specific search requirements
	// This is a basic example that would need to be expanded
	dbQuery := r.DB.Model(&models.Album{}).Preload("Artist").Preload("Genre")

	if query != nil && *query != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+*query+"%")
	}

	if artist != nil && *artist != "" {
		dbQuery = dbQuery.Joins("JOIN artists ON artists.id = albums.artist_id").
			Where("artists.name LIKE ?", "%"+*artist+"%")
	}

	if genre != nil && *genre != "" {
		dbQuery = dbQuery.Joins("JOIN genres ON genres.id = albums.genre_id").
			Where("genres.name LIKE ?", "%"+*genre+"%")
	}

	if sort != nil {
		switch *sort {
		case "newest":
			dbQuery = dbQuery.Order("release_date DESC")
		case "oldest":
			dbQuery = dbQuery.Order("release_date ASC")
		case "popular":
			// Would need a plays_count or similar field
			dbQuery = dbQuery.Order("plays_count DESC")
		}
	}

	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		dbQuery = dbQuery.Offset(offset).Limit(*limit)
	}

	var albums []models.Album
	err := dbQuery.Find(&albums).Error
	return albums, err
}
