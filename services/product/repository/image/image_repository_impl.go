package repository

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type imageRepositoryImpl struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepositoryImpl{db}
}

func (r *imageRepositoryImpl) Create(ctx context.Context, image *model.Image) error {
	if err := r.db.WithContext(ctx).Create(image).Error; err != nil {
		return err
	}

	return nil
}

func (r *imageRepositoryImpl) CreateAll(ctx context.Context, images []*model.Image) error {
	if err := r.db.WithContext(ctx).Create(&images).Error; err != nil {
		return err
	}

	return nil
}

func (r *imageRepositoryImpl) UpdateIsDeletedByIDIn(ctx context.Context, ids []string) error {
	result := r.db.WithContext(ctx).Model(&model.Image{}).Where("id IN ?", ids).Update("is_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrHasImageNotFound
	}

	return nil
}

func (r *imageRepositoryImpl) FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Image, error) {
	var images []*model.Image
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&images).Error; err != nil {
		return nil, err
	}

	return images, nil
}

func (r *imageRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Image{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrImageNotFound
	}

	return nil
}
