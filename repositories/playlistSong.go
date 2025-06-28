package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PlaylistSongRepository struct {
	DB *gorm.DB
}

func NewPlaylistSongRepository(db *gorm.DB) IPlaylistSongRepository {
	return &PlaylistSongRepository{DB: db}
}

func (r *PlaylistSongRepository) AddSongToPlaylist(playlistID, songID uuid.UUID, position int) error {
	return r.DB.Create(&models.PlaylistSong{
		PlaylistID: playlistID,
		SongID:     songID,
		Position:   position,
		AddedAt:    time.Now(),
	}).Error
}

func (r *PlaylistSongRepository) RemoveSongFromPlaylist(playlistID, songID uuid.UUID) error {
	return r.DB.
		Where("playlist_id = ? AND song_id = ?", playlistID, songID).
		Delete(&models.PlaylistSong{}).
		Error
}

func (r *PlaylistSongRepository) ReorderSongs(playlistID uuid.UUID, songOrder map[uuid.UUID]int) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for songID, position := range songOrder {
		if err := tx.Model(&models.PlaylistSong{}).
			Where("playlist_id = ? AND song_id = ?", playlistID, songID).
			Update("position", position).
			Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *PlaylistSongRepository) GetPlaylistSongs(playlistID uuid.UUID) ([]models.Song, error) {
	var songs []models.Song
	err := r.DB.
		Joins("JOIN playlist_songs ON songs.id = playlist_songs.song_id").
		Where("playlist_songs.playlist_id = ?", playlistID).
		Order("playlist_songs.position").
		Preload("Artist").
		Find(&songs).
		Error
	return songs, err
}

func (r *PlaylistSongRepository) SongInPlaylistExists(playlistID uuid.UUID, songID uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Where("playlist_id = ? AND song_id = ?", playlistID, songID).Count(&count).Error
	return count > 0, err
}
