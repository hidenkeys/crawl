package repositories

import (
	"crawl/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SongContributorRepository struct {
	DB *gorm.DB
}

func NewSongContributorRepository(db *gorm.DB) ISongContributorRepository {
	return &SongContributorRepository{DB: db}
}

func (r *SongContributorRepository) FindBySongID(songID uuid.UUID) ([]models.SongContributor, error) {
	var contributors []models.SongContributor
	err := r.DB.
		Where("song_id = ?", songID).
		Find(&contributors).
		Error
	return contributors, err
}

func (r *SongContributorRepository) FindByArtistID(artistID uuid.UUID) ([]models.SongContributor, error) {
	var contributors []models.SongContributor
	err := r.DB.
		Where("artist_id = ?", artistID).
		Find(&contributors).
		Error
	return contributors, err
}

func (r *SongContributorRepository) AddContributor(contributor *models.SongContributor) error {
	return r.DB.Create(contributor).Error
}

func (r *SongContributorRepository) RemoveContributor(songID, artistID uuid.UUID, contributionType string) error {
	return r.DB.
		Where("song_id = ? AND artist_id = ? AND contribution_type = ?", songID, artistID, contributionType).
		Delete(&models.SongContributor{}).
		Error
}

func (r *SongContributorRepository) UpdateRoyalty(songID, artistID uuid.UUID, contributionType string, royalty int) error {
	return r.DB.Model(&models.SongContributor{}).
		Where("song_id = ? AND artist_id = ? AND contribution_type = ?", songID, artistID, contributionType).
		Update("royalty_percentage", royalty).
		Error
}
