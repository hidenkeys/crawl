package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlbumContributorRepository struct {
	DB *gorm.DB
}

func NewAlbumContributorRepository(db *gorm.DB) IAlbumContributorRepository {
	return &AlbumContributorRepository{DB: db}
}

func (r *AlbumContributorRepository) FindByAlbumID(albumID uuid.UUID) ([]models.AlbumContributor, error) {
	var contributors []models.AlbumContributor
	err := r.DB.
		Where("album_id = ?", albumID).
		Find(&contributors).
		Error
	return contributors, err
}

func (r *AlbumContributorRepository) AddContributor(contributor *models.AlbumContributor) error {
	return r.DB.Create(contributor).Error
}

func (r *AlbumContributorRepository) RemoveContributor(albumID, artistID uuid.UUID, contributionType string) error {
	return r.DB.
		Where("album_id = ? AND artist_id = ? AND contribution_type = ?", albumID, artistID, contributionType).
		Delete(&models.AlbumContributor{}).
		Error
}
