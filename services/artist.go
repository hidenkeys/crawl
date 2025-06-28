package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArtistService interface {
	SearchArtistsByName(ctx context.Context, query string, page int, limit int) ([]models.Artist, error)
	CreateArtist(ctx context.Context, artist *models.Artist) (*models.Artist, error)
	GetArtistByID(ctx context.Context, artistID uuid.UUID) (*models.Artist, error)
	GetAllArtists(ctx context.Context, page *int, limit *int) ([]models.Artist, error)
	UpdateArtist(ctx context.Context, artistID uuid.UUID, artist *models.Artist) (*models.Artist, error)
	GetArtistSongs(ctx context.Context, artistID uuid.UUID, page *int, limit *int) ([]models.Song, error)
}

type artistService struct {
	artistRepo repositories.IArtistRepository
	songRepo   repositories.ISongRepository
	userRepo   repositories.IUserRepository
}

func NewArtistService(
	artistRepo repositories.IArtistRepository,
	songRepo repositories.ISongRepository,
	userRepo repositories.IUserRepository,
) ArtistService {
	return &artistService{
		artistRepo: artistRepo,
		songRepo:   songRepo,
		userRepo:   userRepo,
	}
}

func (s *artistService) SearchArtistsByName(ctx context.Context, query string, page int, limit int) ([]models.Artist, error) {
	return s.artistRepo.SearchByName(query, limit, page)
}

func (s *artistService) CreateArtist(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	if artist.ArtistName == "" {
		return nil, errors.New("artist name is required")
	}

	newArtist, err := s.artistRepo.Create(artist)
	if err == nil {
		err := s.userRepo.SetUserAsArtist(newArtist.UserID)
		if err != nil {
			log.Infof("Failed to change user isArtist to true")
			return nil, err
		}
		return newArtist, nil
	}
	return nil, err
}

func (s *artistService) GetArtistByID(ctx context.Context, artistID uuid.UUID) (*models.Artist, error) {
	artist, err := s.artistRepo.GetWithAlbums(artistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return artist, nil
}

func (s *artistService) GetAllArtists(ctx context.Context, page *int, limit *int) ([]models.Artist, error) {
	var offset int
	if page != nil && limit != nil {
		offset = (*page - 1) * *limit
	}

	where := make([]interface{}, 0) // Empty where clause for all records
	return s.artistRepo.GetAll(offset, *limit, where...)
}

func (s *artistService) UpdateArtist(ctx context.Context, artistID uuid.UUID, artist *models.Artist) (*models.Artist, error) {
	existingArtist, err := s.artistRepo.GetByID(artistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}

	// Update fields
	existingArtist.ArtistName = artist.ArtistName
	existingArtist.WalletBalance = artist.WalletBalance
	existingArtist.StripeAccountID = artist.StripeAccountID
	existingArtist.Verified = artist.Verified
	existingArtist.MonthlyListeners = artist.MonthlyListeners

	return s.artistRepo.Update(existingArtist)
}

func (s *artistService) GetArtistSongs(ctx context.Context, artistID uuid.UUID, page *int, limit *int) ([]models.Song, error) {
	// First verify artist exists
	_, err := s.artistRepo.GetByID(artistID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}

	//var offset int
	/*if page != nil && limit != nil {
		offset = (*page - 1) * *limit
	}*/

	songs, err := s.songRepo.GetByArtist(artistID)
	if err != nil {
		return nil, err
	}

	return songs, nil
}
