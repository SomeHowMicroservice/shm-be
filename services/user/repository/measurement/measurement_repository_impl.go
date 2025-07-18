package repository

import (
	"context"
	"errors"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"gorm.io/gorm"
)

type measurementRepositoryImpl struct {
	db *gorm.DB
}

func NewMeasurementRepository(db *gorm.DB) MeasurementRepository {
	return &measurementRepositoryImpl{db}
}

func (r *measurementRepositoryImpl) Create(ctx context.Context, measurement *model.Measurement) error {
	if err := r.db.WithContext(ctx).Create(measurement).Error; err != nil {
		return err
	}

	return nil
}

func (r *measurementRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Measurement{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrMeasurementNotFound
	}

	return nil
}

func (r *measurementRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Measurement, error) {
	var measurement model.Measurement
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&measurement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &measurement, nil
}

func (r *measurementRepositoryImpl) FindByUserID(ctx context.Context, userID string) (*model.Measurement, error) {
	var measurement model.Measurement
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&measurement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &measurement, nil
}
