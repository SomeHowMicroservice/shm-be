package repository

import (
	"context"

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