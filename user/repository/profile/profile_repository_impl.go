package repository

import (
	"context"
	"errors"

	"github.com/SomeHowMicroservice/shm-be/user/common"
	"github.com/SomeHowMicroservice/shm-be/user/model"
	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepositoryImpl{db}
}

func (r *profileRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return common.ErrProfileNotFound
	}

	return nil
}

func (r *profileRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Profile, error) {
	var profile model.Profile
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}
