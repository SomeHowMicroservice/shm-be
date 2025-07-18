package repository

import (
	"context"
	"errors"

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

func (r *profileRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrProfileNotFound
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
