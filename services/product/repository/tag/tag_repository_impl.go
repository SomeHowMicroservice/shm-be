package repository

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type tagRepositoryImpl struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepositoryImpl{db}
}

func (r *tagRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Tag{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *tagRepositoryImpl) Create(ctx context.Context, tag *model.Tag) error {
	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		return err
	}

	return nil
}

func (r *tagRepositoryImpl) FindAll(ctx context.Context) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.WithContext(ctx).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

func (r *tagRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &tag, nil
}

func (r *tagRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Tag{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrTagNotFound
	}

	return nil
}

func (r *tagRepositoryImpl) FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Tag, error) {
	var tags []*model.Tag
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}