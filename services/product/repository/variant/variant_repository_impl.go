package repository

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
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

func (r *variantRepositoryImpl) CreateAll(ctx context.Context, variants []*model.Variant) error {
	if err := r.db.WithContext(ctx).Create(&variants).Error; err != nil {
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

func (r *variantRepositoryImpl) FindAllByID(ctx context.Context, ids []string) ([]*model.Variant, error) {
	var variants []*model.Variant
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&variants).Error; err != nil {
		return nil, err
	}

	return variants, nil
}

func (r *variantRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Variant{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrVariantNotFound
	}

	return nil
}

func (r *variantRepositoryImpl) DeleteAllByID(ctx context.Context, ids []string) error {
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.Variant{}).Error; err != nil {
		return err
	}

	return nil
}
