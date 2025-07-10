package services

import (
	"context"
	"crawl/api"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumService interface {
	SearchAlbums(ctx context.Context, query *string, artist *string, genre *string, sort *string, page *int, limit *int) ([]models.Album, error)
	CreateAlbum(ctx context.Context, album models.Album) (*models.Album, error)
	GetAllAlbums(ctx context.Context, params api.GetAlbumsParams) ([]models.Album, error)
	GetAlbumByID(ctx context.Context, albumID uuid.UUID) (*models.Album, error)
	GetAllArtistAlbums(ctx context.Context, artistID uuid.UUID, page int, limit int) ([]models.Album, error)
	UpdateAlbum(ctx context.Context, albumID uuid.UUID, album *models.Album) (*models.Album, error)
	DeleteAlbum(ctx context.Context, albumID uuid.UUID) error
	GetAlbumContributors(ctx context.Context, albumID uuid.UUID) ([]models.AlbumContributor, error)
	AddAlbumContributor(ctx context.Context, albumID uuid.UUID, contributor *models.AlbumContributor) error
	GetAlbumSongs(ctx context.Context, albumID uuid.UUID) ([]models.Song, error)
	FlagAlbum(ctx context.Context, albumID uuid.UUID) error
	LikeAlbum(ctx context.Context, albumID uuid.UUID) error
	CreateAlbumWithSongs(ctx context.Context, album *models.Album, songs []models.Song) (*models.Album, error)
}

type albumService struct {
	albumRepo            repositories.IAlbumRepository
	albumContributorRepo repositories.IAlbumContributorRepository
	songRepo             repositories.ISongRepository
	artistRepo           repositories.IArtistRepository
}

func NewAlbumService(
	albumRepo repositories.IAlbumRepository,
	albumContributorRepo repositories.IAlbumContributorRepository,
	songRepo repositories.ISongRepository,
	artistRepo repositories.IArtistRepository,
) AlbumService {
	return &albumService{
		albumRepo:            albumRepo,
		albumContributorRepo: albumContributorRepo,
		songRepo:             songRepo,
		artistRepo:           artistRepo,
	}
}

func (s *albumService) SearchAlbums(ctx context.Context, query *string, artist *string, genre *string, sort *string, page *int, limit *int) ([]models.Album, error) {
	return s.albumRepo.SearchAlbums(query, artist, genre, sort, page, limit)
}

func (s *albumService) CreateAlbum(ctx context.Context, album models.Album) (*models.Album, error) {
	// Validate required fields
	if album.Title == "" {
		return nil, errors.New("album title is required")
	}
	if album.ArtistID == uuid.Nil {
		return nil, errors.New("artist ID is required")
	}

	return s.albumRepo.Create(&album)
}

func (s *albumService) GetAlbumByID(ctx context.Context, albumID uuid.UUID) (*models.Album, error) {
	album, err := s.albumRepo.GetWithSongs(albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("album not found")
		}
		return nil, err
	}
	return album, nil
}

func (s *albumService) GetAllAlbums(ctx context.Context, params api.GetAlbumsParams) ([]models.Album, error) {
	return s.albumRepo.GetAll(*params.Page, *params.Limit)
}

func (s *albumService) GetAllArtistAlbums(ctx context.Context, artist uuid.UUID, page int, limit int) ([]models.Album, error) {

	return s.albumRepo.GetByArtist(artist)
}

func (s *albumService) UpdateAlbum(ctx context.Context, albumID uuid.UUID, album *models.Album) (*models.Album, error) {
	// First check if album exists
	existingAlbum, err := s.albumRepo.GetByID(albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("album not found")
		}
		return nil, err
	}

	// Update the existing album with new values
	existingAlbum.Title = album.Title
	existingAlbum.Description = album.Description
	existingAlbum.Price = album.Price
	existingAlbum.CoverImageURL = album.CoverImageURL
	existingAlbum.ReleaseDate = album.ReleaseDate
	existingAlbum.GenreID = album.GenreID
	existingAlbum.IsFlagged = album.IsFlagged

	return s.albumRepo.Update(existingAlbum)
}

func (s *albumService) DeleteAlbum(ctx context.Context, albumID uuid.UUID) error {
	// First check if album exists
	_, err := s.albumRepo.GetByID(albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("album not found")
		}
		return err
	}

	return s.albumRepo.Delete(albumID)
}

func (s *albumService) GetAlbumContributors(ctx context.Context, albumID uuid.UUID) ([]models.AlbumContributor, error) {
	return s.albumContributorRepo.FindByAlbumID(albumID)
}

func (s *albumService) AddAlbumContributor(ctx context.Context, albumID uuid.UUID, contributor *models.AlbumContributor) error {
	// Validate the contributor
	if contributor.ArtistID == uuid.Nil {
		return errors.New("artist ID is required")
	}
	if contributor.ContributionType == "" {
		return errors.New("contribution type is required")
	}

	// Set the album ID
	contributor.AlbumID = albumID

	return s.albumContributorRepo.AddContributor(contributor)
}

func (s *albumService) GetAlbumSongs(ctx context.Context, albumID uuid.UUID) ([]models.Song, error) {
	album, err := s.albumRepo.GetWithSongs(albumID)
	if err != nil {
		return nil, err
	}
	return album.Songs, nil
}

func (s *albumService) FlagAlbum(ctx context.Context, albumID uuid.UUID) error {
	err := s.albumRepo.FlagContent(albumID)
	if err != nil {
		return errors.New("Album not found")
	}
	return nil
}

func (s *albumService) LikeAlbum(ctx context.Context, albumID uuid.UUID) error {
	err := s.albumRepo.AddLikes(albumID, 1)
	if err != nil {
		return errors.New("album not found")
	}
	return nil
}

func (s *albumService) CreateAlbumWithSongs(ctx context.Context, album *models.Album, songs []models.Song) (*models.Album, error) {
	// Basic validation for album
	if album.Title == "" {
		return nil, errors.New("album title is required")
	}
	if album.ArtistID == uuid.Nil {
		return nil, errors.New("artist ID is required")
	}

	// Basic validation for songs
	if len(songs) == 0 {
		return nil, errors.New("at least one song is required for the album")
	}
	for _, song := range songs {
		if song.Title == "" {
			return nil, errors.New("song title is required")
		}
		if song.AudioURL == "" {
			return nil, errors.New("song audio URL is required")
		}
		// Add more song validations as needed
	}

	createdAlbum, err := s.albumRepo.CreateAlbumWithSongs(album, songs)
	if err != nil {
		return nil, err
	}

	// Increment song upload count for the artist
	err = s.artistRepo.AddSongUploadCount(album.ArtistID, len(songs))
	if err != nil {
		// Log the error but don't return it, as the album and songs are already created
		// In a real application, you might want to handle this more robustly (e.g., a retry mechanism)
		println("Warning: Failed to increment song upload count for artist: " + err.Error())
	}

	return createdAlbum, nil
}
