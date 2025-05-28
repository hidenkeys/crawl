package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongRepositoryImpl struct {
	DB *gorm.DB
}

func NewSongRepository(db *gorm.DB) SongRepository {
	return &SongRepositoryImpl{DB: db}
}

func (r *SongRepositoryImpl) CreateSong(song *models.Song) (*models.Song, error) {
	err := r.DB.Create(song).Error
	if err != nil {
		return nil, err
	}
	return song, nil
}

func (r *SongRepositoryImpl) GetSongByID(id uuid.UUID) (*models.Song, error) {
	var song models.Song
	err := r.DB.First(&song, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (r *SongRepositoryImpl) GetSongsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Song, int64, error) {
	var songs []models.Song
	var count int64
	err := r.DB.Model(&models.Song{}).Where("artist_id = ?", artistID).Count(&count).Offset(offset).Limit(limit).Find(&songs).Error
	if err != nil {
		return nil, count, err
	}
	return songs, count, nil
}

func (r *SongRepositoryImpl) UpdateSong(id uuid.UUID, song *models.Song) (*models.Song, error) {
	var existingSong models.Song
	err := r.DB.First(&existingSong, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&existingSong).Updates(song).Error
	if err != nil {
		return nil, err
	}
	return &existingSong, nil
}

func (r *SongRepositoryImpl) DeleteSong(id uuid.UUID) error {
	var song models.Song
	err := r.DB.First(&song, "id = ?", id).Error
	if err != nil {
		return err
	}
	return r.DB.Delete(&song).Error
}
