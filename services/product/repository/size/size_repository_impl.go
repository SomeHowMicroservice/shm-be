package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type sizeRepositoryImpl struct {
	db *gorm.DB
}

func NewSizeRepository(db *gorm.DB) SizeRepository {
	return &sizeRepositoryImpl{db}
}

func (r *sizeRepositoryImpl) Create(ctx context.Context, size *model.Size) error {
	if err := r.db.WithContext(ctx).Create(size).Error; err != nil {
		return err
	}

	return nil
}

func (r *sizeRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Size{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *sizeRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Size{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
