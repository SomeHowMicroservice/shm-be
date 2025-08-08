package repository

import (
	"context"
	"errors"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type colorRepositoryImpl struct {
	db *gorm.DB
}

func NewColorRepository(db *gorm.DB) ColorRepository {
	return &colorRepositoryImpl{db}
}

func (r *colorRepositoryImpl) Create(ctx context.Context, color *model.Color) error {
	if err := r.db.WithContext(ctx).Create(color).Error; err != nil {
		return err
	}

	return nil
}

func (r *colorRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Color{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *colorRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Color{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *colorRepositoryImpl) FindAll(ctx context.Context) ([]*model.Color, error) {
	var colors []*model.Color
	if err := r.db.WithContext(ctx).Find(&colors).Error; err != nil {
		return nil, err
	}

	return colors, nil
}

func (r *colorRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Color, error) {
	var color model.Color
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&color).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &color, nil
}

func (r *colorRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Color{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrColorNotFound
	}

	return nil
}
