package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumRepositoryImpl struct {
	DB *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &AlbumRepositoryImpl{DB: db}
}

func (r *AlbumRepositoryImpl) CreateAlbum(album *models.Album) (*models.Album, error) {
	err := r.DB.Create(album).Error
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (r *AlbumRepositoryImpl) GetAlbumByID(id uuid.UUID) (*models.Album, error) {
	var album models.Album
	err := r.DB.First(&album, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (r *AlbumRepositoryImpl) GetAlbumsByArtist(artistID uuid.UUID, limit, offset int) ([]models.Album, int64, error) {
	var albums []models.Album
	var count int64
	err := r.DB.Model(&models.Album{}).Where("artist_id = ?", artistID).Count(&count).Offset(offset).Limit(limit).Find(&albums).Error
	if err != nil {
		return nil, count, err
	}
	return albums, count, nil
}

func (r *AlbumRepositoryImpl) UpdateAlbum(id uuid.UUID, album *models.Album) (*models.Album, error) {
	var existingAlbum models.Album
	err := r.DB.First(&existingAlbum, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	err = r.DB.Model(&existingAlbum).Updates(album).Error
	if err != nil {
		return nil, err
	}
	return &existingAlbum, nil
}

func (r *AlbumRepositoryImpl) DeleteAlbum(id uuid.UUID) error {
	var album models.Album
	err := r.DB.First(&album, "id = ?", id).Error
	if err != nil {
		return err
	}
	return r.DB.Delete(&album).Error
}
