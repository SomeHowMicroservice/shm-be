package repository

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepositoryImpl{db}
}

func (r *profileRepositoryImpl) Create(ctx context.Context, profile *model.Profile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return err
	}
	return nil
}

func (r *profileRepositoryImpl) UpdateByUserID(ctx context.Context, userID string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", userID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrProfileNotFound
	}
	return nil
}
