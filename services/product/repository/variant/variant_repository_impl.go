package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type variantRepositoryImpl struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) VariantRepository {
	return &variantRepositoryImpl{db}
}

func (r *variantRepositoryImpl) Create(ctx context.Context, variant *model.Variant) error {
	if err := r.db.WithContext(ctx).Create(variant).Error; err != nil {
		return err
	}

	return nil
}

func (r *variantRepositoryImpl) ExistsBySKU(ctx context.Context, sku string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Variant{}).Where("sku = ?", sku).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}