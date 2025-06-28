package services

import (
	"context"
	"crawl/models"
	"crawl/repositories"
	"errors"
	"github.com/google/uuid"
)

type PlaylistService interface {
	GetPlaylistByID(ctx context.Context, playlistID uuid.UUID) (*models.Playlist, error)
	GetAllPlaylists(ctx context.Context) ([]models.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlistID uuid.UUID, playlist *models.Playlist) (*models.Playlist, error)
	DeletePlaylist(ctx context.Context, playlistID uuid.UUID) error
	GetPlaylistSongs(ctx context.Context, playlistID uuid.UUID) ([]models.Song, error)
	AddSongToPlaylist(ctx context.Context, playlistID uuid.UUID, songID uuid.UUID) error
	RemoveSongFromPlaylist(ctx context.Context, playlistID uuid.UUID, songID uuid.UUID) error
	SearchPlaylists(
		ctx context.Context,
		query *string,
		owner *string,
		isPublic *bool,
		sort *string,
		page int,
		limit int,
	) ([]models.Playlist, int64, error)
}

type playlistService struct {
	playlistRepo     repositories.IPlaylistRepository
	playlistSongRepo repositories.IPlaylistSongRepository
	songRepo         repositories.ISongRepository
}

func NewPlaylistService(
	playlistRepo repositories.IPlaylistRepository,
	playlistSongRepo repositories.IPlaylistSongRepository,
	songRepo repositories.ISongRepository,
) PlaylistService {
	return &playlistService{
		playlistRepo:     playlistRepo,
		playlistSongRepo: playlistSongRepo,
		songRepo:         songRepo,
	}
}

func (s *playlistService) GetPlaylistByID(ctx context.Context, playlistID uuid.UUID) (*models.Playlist, error) {
	playlist, err := s.playlistRepo.GetWithSongs(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("playlist not found")
		}
		return nil, err
	}
	return playlist, nil
}

func (s *playlistService) GetAllPlaylists(ctx context.Context) ([]models.Playlist, error) {
	// Get all playlists without pagination (using 0, 0 for offset/limit)
	playlists, err := s.playlistRepo.GetAll(0, 0, "is_public = true")
	if err != nil {
		return nil, err
	}

	// If no playlists found, return empty slice rather than nil
	if playlists == nil {
		return []models.Playlist{}, nil
	}

	return playlists, nil
}

func (s *playlistService) UpdatePlaylist(ctx context.Context, playlistID uuid.UUID, playlist *models.Playlist) (*models.Playlist, error) {
	// First check if playlist exists
	existingPlaylist, err := s.playlistRepo.GetByID(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("playlist not found")
		}
		return nil, err
	}

	// Update the existing playlist with new values
	existingPlaylist.Title = playlist.Title
	existingPlaylist.Description = playlist.Description
	existingPlaylist.CoverImageURL = playlist.CoverImageURL
	existingPlaylist.IsPublic = playlist.IsPublic

	return s.playlistRepo.Update(existingPlaylist)
}

func (s *playlistService) DeletePlaylist(ctx context.Context, playlistID uuid.UUID) error {
	// First check if playlist exists
	_, err := s.playlistRepo.GetByID(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("playlist not found")
		}
		return err
	}

	return s.playlistRepo.Delete(playlistID)
}

func (s *playlistService) GetPlaylistSongs(ctx context.Context, playlistID uuid.UUID) ([]models.Song, error) {
	// Verify playlist exists first
	_, err := s.playlistRepo.GetByID(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return nil, errors.New("playlist not found")
		}
		return nil, err
	}

	return s.playlistSongRepo.GetPlaylistSongs(playlistID)
}

func (s *playlistService) AddSongToPlaylist(ctx context.Context, playlistID uuid.UUID, songID uuid.UUID) error {
	// Verify playlist exists
	_, err := s.playlistRepo.GetByID(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("playlist not found")
		}
		return err
	}

	// Verify song exists
	_, err = s.songRepo.GetByID(songID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("song not found")
		}
		return err
	}

	// Check if song already exists in playlist
	exists, err := s.playlistSongRepo.SongInPlaylistExists(playlistID, songID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("song already exists in playlist")
	}

	return s.playlistSongRepo.AddSongToPlaylist(playlistID, songID, 0)
}

func (s *playlistService) RemoveSongFromPlaylist(ctx context.Context, playlistID uuid.UUID, songID uuid.UUID) error {
	// Verify playlist exists
	_, err := s.playlistRepo.GetByID(playlistID)
	if err != nil {
		if errors.Is(err, repositories.ErrRecordNotFound) {
			return errors.New("playlist not found")
		}
		return err
	}

	// Verify song exists in playlist
	exists, err := s.playlistSongRepo.SongInPlaylistExists(playlistID, songID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("song not found in playlist")
	}

	return s.playlistSongRepo.RemoveSongFromPlaylist(playlistID, songID)
}

func (s *playlistService) SearchPlaylists(
	ctx context.Context,
	query *string,
	owner *string,
	isPublic *bool,
	sort *string,
	page int,
	limit int,
) ([]models.Playlist, int64, error) {
	return s.playlistRepo.SearchPlaylists(query, owner, isPublic, sort, page, limit)
}
