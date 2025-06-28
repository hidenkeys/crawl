package repositories

import (
	"crawl/models"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type GenreRepository struct {
	BaseRepository[models.Genre]
}

func NewGenreRepository(db *gorm.DB) IGenreRepository {
	return &GenreRepository{
		BaseRepository: BaseRepository[models.Genre]{DB: db},
	}
}

func (r *GenreRepository) GetPopular(limit int) ([]models.Genre, error) {
	var genres []models.Genre
	err := r.DB.
		Model(&models.Genre{}).
		Select("genres.*, COUNT(songs.id) as song_count").
		Joins("LEFT JOIN songs ON songs.genre_id = genres.id").
		Group("genres.id").
		Order("song_count DESC").
		Limit(limit).
		Find(&genres).Error
	return genres, err
}

func (r *GenreRepository) SearchGenres(query *string, sort *string) ([]models.Genre, error) {
	var genres []models.Genre

	// Base query
	q := r.DB.Model(&models.Genre{})

	// Apply search query if provided
	if query != nil && *query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(*query))
		q = q.Where("LOWER(name) LIKE ?", searchTerm)
	}

	// Apply sorting if provided
	if sort != nil {
		switch *sort {
		case "name":
			q = q.Order("name ASC")
		case "popularity":
			// Assuming you have a popularity field or need to join with related data
			q = q.Order("popularity DESC")
		default:
			q = q.Order("name ASC") // Default sort
		}
	}

	// Execute query
	if err := q.Find(&genres).Error; err != nil {
		return nil, fmt.Errorf("failed to search genres: %w", err)
	}

	return genres, nil
}
