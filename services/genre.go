package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenreService interface {
	GetAllGenres(ctx context.Context) ([]models.Genre, error)
	GetGenreByID(ctx context.Context, genreID uuid.UUID) (*models.Genre, error)
	GetPopularGenres(ctx context.Context, limit int) ([]models.Genre, error)
	SearchGenres(ctx context.Context, query *string, sort *string) ([]models.Genre, error)
}
type genreService struct {
	genreRepo repositories.IGenreRepository
}

func NewGenreService(genreRepo repositories.IGenreRepository) GenreService {
	return &genreService{
		genreRepo: genreRepo,
	}
}

func (s *genreService) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	// Get all genres without pagination (using 0, 0 for offset/limit)
	genres, err := s.genreRepo.GetAll(0, 0)
	if err != nil {
		return nil, err
	}

	// If no genres found, return empty slice rather than nil
	if genres == nil {
		return []models.Genre{}, nil
	}

	return genres, nil
}

func (s *genreService) GetGenreByID(ctx context.Context, genreID uuid.UUID) (*models.Genre, error) {
	genre, err := s.genreRepo.GetByID(genreID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("genre not found")
		}
		return nil, err
	}

	return genre, nil
}

func (s *genreService) GetPopularGenres(ctx context.Context, limit int) ([]models.Genre, error) {
	return s.genreRepo.GetPopular(limit)
}

func (s *genreService) SearchGenres(ctx context.Context, query *string, sort *string) ([]models.Genre, error) {
	return s.genreRepo.SearchGenres(query, sort)
}
