package repository

import (
	"context"
	"errors"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
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

func (r *sizeRepositoryImpl) FindAll(ctx context.Context) ([]*model.Size, error) {
	var sizes []*model.Size
	if err := r.db.WithContext(ctx).Where("is_deleted = false").Find(&sizes).Error; err != nil {
		return nil, err
	}

	return sizes, nil
}

func (r *sizeRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Size, error) {
	var size model.Size
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = false", id).First(&size).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &size, nil
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

func (r *sizeRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Size{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrSizeNotFound
	}

	return nil
}

func (r *sizeRepositoryImpl) FindAllByID(ctx context.Context, ids []string) ([]*model.Size, error) {
	var size []*model.Size
	if err := r.db.WithContext(ctx).Where("id IN ? AND is_deleted = false", ids).Find(&size).Error; err != nil {
		return nil, err
	}

	return size, nil
}

func (r *sizeRepositoryImpl) UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(&model.Size{}).Where("id IN ?", ids).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *sizeRepositoryImpl) FindAllDeleted(ctx context.Context) ([]*model.Size, error) {
	var sizes []*model.Size
	if err := r.db.WithContext(ctx).Where("is_deleted = true").Find(&sizes).Error; err != nil {
		return nil, err
	}

	return sizes, nil
}
