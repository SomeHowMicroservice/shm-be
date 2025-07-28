package repository

import (
	"context"
	"errors"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepositoryImpl{db}
}

func (r *productRepositoryImpl) Create(ctx context.Context, product *model.Product) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return err 
	}

	return nil
}

func (r *productRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64 
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *productRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).Preload("Categories").Where("slug = ?", slug).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} 
		return nil, err
	}

	return &product, nil
}

func (r *productRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}