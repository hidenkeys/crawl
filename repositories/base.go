package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"reflect"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func (r *BaseRepository[T]) GetAll(offset int, limit int, where ...interface{}) ([]T, error) {
	var models []T
	db := r.DB

	// Apply WHERE clause if conditions are provided
	if len(where) > 0 {
		// First element is the query, rest are args
		query := where[0].(string)
		args := make([]interface{}, 0)
		if len(where) > 1 {
			args = where[1:]
		}
		db = db.Where(query, args...)
	}

	// Apply pagination
	if limit > 0 {
		db = db.Offset(offset).Limit(limit)
	}

	err := db.Find(&models).Error
	return models, err
}

func (r *BaseRepository[T]) GetByID(id uuid.UUID) (*T, error) {
	var model T
	err := r.DB.First(&model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}
	return &model, err
}

func (r *BaseRepository[T]) Create(model *T) (*T, error) {
	err := r.DB.Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *BaseRepository[T]) Update(model *T) (*T, error) {
	result := r.DB.Model(model).Updates(model)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, ErrEditConflict
	}

	// For some databases/ORM configurations, Updates might not refresh all fields
	// So we fetch the updated record to ensure we have the latest data
	var updatedModel T
	err := r.DB.First(&updatedModel, getID(model)).Error
	if err != nil {
		return nil, err
	}

	return &updatedModel, nil
}

// Helper function to get the ID of a model
// Assumes your models have an ID field of type uuid.UUID
func getID[T any](model *T) uuid.UUID {
	// Use reflection to get the ID field
	val := reflect.ValueOf(model).Elem()
	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		// Try BaseModel.ID if using embedded BaseModel
		baseModel := val.FieldByName("BaseModel")
		if baseModel.IsValid() {
			idField = baseModel.FieldByName("ID")
		}
	}

	if idField.IsValid() {
		return idField.Interface().(uuid.UUID)
	}

	return uuid.Nil
}

func (r *BaseRepository[T]) Delete(id uuid.UUID) error {
	result := r.DB.Delete(new(T), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *BaseRepository[T]) Exists(id uuid.UUID) (bool, error) {
	var count int64
	err := r.DB.Model(new(T)).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *BaseRepository[T]) Count(where ...interface{}) (int64, error) {
	var count int64
	db := r.DB.Model(new(T))

	if len(where) > 0 {
		query := where[0].(string)
		args := make([]interface{}, 0)
		if len(where) > 1 {
			args = where[1:]
		}
		db = db.Where(query, args...)
	}

	err := db.Count(&count).Error
	return count, err
}
