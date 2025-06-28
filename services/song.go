package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type SongService interface {
	SearchSongs(ctx context.Context, query *string, artist *string, genre *string, sort *string, order *string, page *int, limit *int) ([]models.Song, error)
	CreateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	GetSongByID(ctx context.Context, songID uuid.UUID) (*models.Song, error)
	GetAllSongs(ctx context.Context, page *int, limit *int, genre *string, artistID *string, albumID *string) ([]models.Song, error)
	UpdateSong(ctx context.Context, songID uuid.UUID, song *models.Song) (*models.Song, error)
	DeleteSong(ctx context.Context, songID uuid.UUID) error
	GetSongContributors(ctx context.Context, songID uuid.UUID) ([]models.SongContributor, error)
	AddSongContributor(ctx context.Context, songID uuid.UUID, contributor *models.SongContributor) error
	RecordStream(ctx context.Context, stream *models.Stream) error
}

type songService struct {
	songRepo        repositories.ISongRepository
	artistRepo      repositories.IArtistRepository
	genreRepo       repositories.IGenreRepository
	albumRepo       repositories.IAlbumRepository
	streamRepo      repositories.IStreamRepository
	contributorRepo repositories.ISongContributorRepository
}

func NewSongService(
	songRepo repositories.ISongRepository,
	artistRepo repositories.IArtistRepository,
	genreRepo repositories.IGenreRepository,
	albumRepo repositories.IAlbumRepository,
	streamRepo repositories.IStreamRepository,
	contributorRepo repositories.ISongContributorRepository,
) SongService {
	return &songService{
		songRepo:        songRepo,
		artistRepo:      artistRepo,
		genreRepo:       genreRepo,
		albumRepo:       albumRepo,
		streamRepo:      streamRepo,
		contributorRepo: contributorRepo,
	}
}

func (s *songService) SearchSongs(ctx context.Context, query *string, artist *string, genre *string, sort *string, order *string, page *int, limit *int) ([]models.Song, error) {
	var offset int
	if page != nil && limit != nil {
		offset = (*page - 1) * *limit
	} else {
		limit = new(int)
		*limit = 20
	}

	return s.songRepo.Search(query, artist, genre, sort, order, offset, *limit)
}

func (s *songService) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	// Validate required fields
	if song.Title == "" {
		return nil, errors.New("song title is required")
	}
	if song.ArtistID == uuid.Nil {
		return nil, errors.New("artist ID is required")
	}
	if song.Duration <= 0 {
		return nil, errors.New("duration must be positive")
	}

	// Verify artist exists
	_, err := s.artistRepo.GetByID(song.ArtistID)
	if err != nil {
		return nil, errors.New("artist not found")
	}

	// Verify album exists if provided
	if song.AlbumID != nil && *song.AlbumID != uuid.Nil {
		_, err := s.albumRepo.GetByID(*song.AlbumID)
		if err != nil {
			return nil, errors.New("album not found")
		}
	}

	// Verify genre exists if provided
	if song.GenreID != nil {
		_, err := s.genreRepo.GetByID(*song.GenreID)
		if err != nil {
			return nil, errors.New("genre not found")
		}
	}

	return s.songRepo.Create(song)
}

func (s *songService) GetSongByID(ctx context.Context, songID uuid.UUID) (*models.Song, error) {
	song, err := s.songRepo.GetWithArtist(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("song not found")
		}
		return nil, err
	}
	return song, nil
}

func (s *songService) GetAllSongs(ctx context.Context, page *int, limit *int, genre *string, artistID *string, albumID *string) ([]models.Song, error) {
	var offset int
	if page != nil && limit != nil {
		offset = (*page - 1) * *limit
	} else {
		limit = new(int)
		*limit = 20
	}

	// Build where conditions
	var conditions []interface{}
	if genre != nil {
		conditions = append(conditions, "genre_id = ?", *genre)
	}
	if artistID != nil {
		id, err := uuid.Parse(*artistID)
		if err != nil {
			return nil, errors.New("invalid artist ID format")
		}
		conditions = append(conditions, "artist_id = ?", id)
	}
	if albumID != nil {
		id, err := uuid.Parse(*albumID)
		if err != nil {
			return nil, errors.New("invalid album ID format")
		}
		conditions = append(conditions, "album_id = ?", id)
	}

	if len(conditions) > 0 {
		return s.songRepo.GetAll(offset, *limit, conditions...)
	}
	return s.songRepo.GetAll(offset, *limit)
}

func (s *songService) UpdateSong(ctx context.Context, songID uuid.UUID, song *models.Song) (*models.Song, error) {
	// Verify song exists
	existingSong, err := s.songRepo.GetByID(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("song not found")
		}
		return nil, err
	}

	// Update fields
	existingSong.Title = song.Title
	existingSong.Duration = song.Duration
	existingSong.Price = song.Price
	existingSong.AudioURL = song.AudioURL
	existingSong.PreviewURL = song.PreviewURL
	existingSong.ReleaseDate = song.ReleaseDate
	existingSong.CoverImageURL = song.CoverImageURL
	existingSong.GenreID = song.GenreID
	existingSong.IsFlagged = song.IsFlagged

	// Verify artist exists if changing artist
	if song.ArtistID != uuid.Nil && existingSong.ArtistID != song.ArtistID {
		_, err := s.artistRepo.GetByID(song.ArtistID)
		if err != nil {
			return nil, errors.New("artist not found")
		}
		existingSong.ArtistID = song.ArtistID
	}

	// Verify album exists if changing album
	if song.AlbumID != nil {
		if *song.AlbumID == uuid.Nil {
			existingSong.AlbumID = nil
		} else {
			_, err := s.albumRepo.GetByID(*song.AlbumID)
			if err != nil {
				return nil, errors.New("album not found")
			}
			existingSong.AlbumID = song.AlbumID
		}
	}

	return s.songRepo.Update(existingSong)
}

func (s *songService) DeleteSong(ctx context.Context, songID uuid.UUID) error {
	// Verify song exists
	_, err := s.songRepo.GetByID(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("song not found")
		}
		return err
	}

	return s.songRepo.Delete(songID)
}

func (s *songService) GetSongContributors(ctx context.Context, songID uuid.UUID) ([]models.SongContributor, error) {
	// Verify song exists
	_, err := s.songRepo.GetByID(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("song not found")
		}
		return nil, err
	}

	return s.contributorRepo.FindBySongID(songID)
}

func (s *songService) AddSongContributor(ctx context.Context, songID uuid.UUID, contributor *models.SongContributor) error {
	// Verify song exists
	_, err := s.songRepo.GetByID(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("song not found")
		}
		return err
	}

	// Verify artist exists
	_, err = s.artistRepo.GetByID(contributor.ArtistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("artist not found")
		}
		return err
	}

	// Set song ID
	contributor.SongID = songID

	return s.contributorRepo.AddContributor(contributor)
}

func (s *songService) RecordStream(ctx context.Context, stream *models.Stream) error {
	// Verify song exists
	_, err := s.songRepo.GetByID(stream.SongID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("song not found")
		}
		return err
	}

	// Record the stream
	_, err = s.streamRepo.Create(stream)
	if err != nil {
		return err
	}

	// Increment play count
	err = s.songRepo.AddPlayCount(stream.SongID, 1)
	if err != nil {
		return err
	}

	return nil
}
