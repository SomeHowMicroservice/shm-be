package repository

import (
	"context"
	"errors"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"gorm.io/gorm"
)

type categoryRepositoryImpl struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepositoryImpl{db}
}

func (r *categoryRepositoryImpl) Create(ctx context.Context, category *model.Category) error {
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Category, error) {
	return r.findAllByIDInBase(ctx, ids)
}

func (r *categoryRepositoryImpl) FindAllByIDInWithChildren(ctx context.Context, ids []string) ([]*model.Category, error) {
	return r.findAllByIDInBase(ctx, ids, "Children")
}

func (r *categoryRepositoryImpl) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Category, error) {
	return r.findByIDBase(ctx, id)
}

func (r *categoryRepositoryImpl) FindAllWithParentsAndChildren(ctx context.Context) ([]*model.Category, error) {
	return r.findAllBase(ctx, "Parents", "Children")
}

func (r *categoryRepositoryImpl) FindAll(ctx context.Context) ([]*model.Category, error) {
	return r.findAllBase(ctx)
}

func (r *categoryRepositoryImpl) FindByIDWithParentsAndProducts(ctx context.Context, id string) (*model.Category, error) {
	return r.findByIDBase(ctx, id, common.Preload{Relation: "Parents"}, common.Preload{Relation: "Products"}, common.Preload{Relation: "Products.Images", Scope: func(db *gorm.DB) *gorm.DB { return db.Where("is_thumbnail = ?", true) }})
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepositoryImpl) FindByIDWithParents(ctx context.Context, id string) (*model.Category, error) {
	return r.findByIDBase(ctx, id, common.Preload{Relation: "Parents"})
}

func (r *categoryRepositoryImpl) UpdateParents(ctx context.Context, category *model.Category, parents []*model.Category) error {
	if err := r.db.WithContext(ctx).Model(category).Association("Parents").Replace(parents); err != nil {
		return err
	}

	return nil
}

func (r *categoryRepositoryImpl) FindAllWithChildren(ctx context.Context) ([]*model.Category, error) {
	return r.findAllBase(ctx, "Children")
}

func (r *categoryRepositoryImpl) findByIDBase(ctx context.Context, id string, preloads ...common.Preload) (*model.Category, error) {
	var category model.Category
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		if preload.Scope != nil {
			query = query.Preload(preload.Relation, preload.Scope)
		} else {
			query = query.Preload(preload.Relation)
		}
	}

	if err := query.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepositoryImpl) findAllBase(ctx context.Context, preloads ...string) ([]*model.Category, error) {
	var categories []*model.Category
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepositoryImpl) findAllByIDInBase(ctx context.Context, ids []string, preloads ...string) ([]*model.Category, error) {
	var categories []*model.Category
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Where("id IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
