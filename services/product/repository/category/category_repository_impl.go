package repository

import (
	"context"
	"errors"

	// customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Category, error) {
	var categories []*model.Category

	if err := r.db.WithContext(ctx).Preload("Children").Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*model.Category, error) {
	var category model.Category

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category

	if err := r.db.WithContext(ctx).Preload("Children").Preload("Parents").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
