package repositories

import (
	"crawl/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MonthlyRoyaltyRepository struct {
	BaseRepository[models.MonthlyRoyalty]
}

func NewMonthlyRoyaltyRepository(db *gorm.DB) IMonthlyRoyaltyRepository {
	return &MonthlyRoyaltyRepository{
		BaseRepository: BaseRepository[models.MonthlyRoyalty]{DB: db},
	}
}

func (r *MonthlyRoyaltyRepository) FindByArtistAndPeriod(artistID uuid.UUID, year int, month int) (*models.MonthlyRoyalty, error) {
	var royalty models.MonthlyRoyalty
	err := r.DB.
		Where("artist_id = ? AND year = ? AND month = ?", artistID, year, month).
		First(&royalty).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &royalty, err
}

func (r *MonthlyRoyaltyRepository) MarkAsPaid(artistID uuid.UUID, year int, month int) error {
	return r.DB.Model(&models.MonthlyRoyalty{}).
		Where("artist_id = ? AND year = ? AND month = ?", artistID, year, month).
		Update("paid_status", true).
		Error
}

func (r *MonthlyRoyaltyRepository) GetArtistRoyalties(artistID uuid.UUID) ([]models.MonthlyRoyalty, error) {
	var royalties []models.MonthlyRoyalty
	err := r.DB.
		Where("artist_id = ?", artistID).
		Order("year DESC, month DESC").
		Find(&royalties).
		Error
	return royalties, err
}

func (r *MonthlyRoyaltyRepository) CalculatePendingRoyalties() (float64, error) {
	var total float64
	err := r.DB.Model(&models.MonthlyRoyalty{}).
		Where("paid_status = ?", false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).
		Error
	return total, err
}
