package repositories

import (
	"crawl/models"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"gorm.io/gorm"
)

type PlaylistRepository struct {
	BaseRepository[models.Playlist]
}

func NewPlaylistRepository(db *gorm.DB) IPlaylistRepository {
	return &PlaylistRepository{
		BaseRepository: BaseRepository[models.Playlist]{DB: db},
	}
}

func (r *PlaylistRepository) GetWithSongs(id uuid.UUID) (*models.Playlist, error) {
	var playlist models.Playlist
	err := r.DB.
		Preload("Songs").
		Preload("Songs.Artist").
		First(&playlist, id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &playlist, err
}

func (r *PlaylistRepository) GetUserPlaylists(userID uuid.UUID) ([]models.Playlist, error) {
	var playlists []models.Playlist
	err := r.DB.
		Where("user_id = ?", userID).
		Find(&playlists).
		Error
	return playlists, err
}

func (r *PlaylistRepository) GetUserPublicPlaylists(userID uuid.UUID) ([]models.Playlist, error) {
	var playlists []models.Playlist
	err := r.DB.
		Where("user_id = ? and is_public = true", userID).
		Find(&playlists).
		Error
	return playlists, err
}

func (r *PlaylistRepository) AddSongToPlaylist(playlistID, songID uuid.UUID) error {
	// Get max position
	var maxPos int
	err := r.DB.Model(&models.PlaylistSong{}).
		Where("playlist_id = ?", playlistID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&maxPos).
		Error
	if err != nil {
		return err
	}

	return r.DB.Create(&models.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
		Position:   maxPos + 1,
		AddedAt:    time.Now(),
	}).Error
}

func (r *PlaylistRepository) SearchPlaylists(
	query *string,
	owner *string,
	isPublic *bool,
	sort *string,
	page int,
	limit int,
) ([]models.Playlist, int64, error) {
	var playlists []models.Playlist
	var total int64

	// Base query with user preload
	q := r.DB.Model(&models.Playlist{}).Preload("User")

	// Apply search query if provided
	if query != nil && *query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(*query))
		q = q.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Filter by owner (username or ID)
	if owner != nil && *owner != "" {
		q = q.Joins("JOIN users ON users.id = playlists.user_id").
			Where("LOWER(users.username) LIKE ? OR users.id::text = ?",
				fmt.Sprintf("%%%s%%", strings.ToLower(*owner)), *owner)
	}

	// Filter by public/private status
	if isPublic != nil {
		q = q.Where("is_public = ?", *isPublic)
	}

	// Apply sorting
	if sort != nil {
		switch *sort {
		case "title":
			q = q.Order("title ASC")
		case "created_at":
			q = q.Order("created_at DESC")
		case "updated_at":
			q = q.Order("updated_at DESC")
		case "popularity":
			// Assuming you track playlist popularity (e.g., through likes or plays)
			//q = q.Order("(SELECT COUNT(*) FROM playlist_likes WHERE playlist_likes.playlist_id = playlists.id) DESC")
		default:
			q = q.Order("created_at DESC") // Default sort
		}
	}

	// Count total before pagination
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count playlists: %w", err)
	}

	// Apply pagination
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		q = q.Offset(offset).Limit(limit)
	}

	// Execute query
	if err := q.Find(&playlists).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search playlists: %w", err)
	}

	return playlists, total, nil
}
