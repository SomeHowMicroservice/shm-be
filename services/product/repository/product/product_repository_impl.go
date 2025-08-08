package repository

import (
	"context"
	"errors"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
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

func (r *productRepositoryImpl) FindBySlugWithDetails(ctx context.Context, slug string) (*model.Product, error) {
	return r.findBySlugBase(ctx, slug, "Categories", "Tags", "Variants", "Variants.Color", "Variants.Size", "Variants.Inventory", "Images", "Images.Color")
}

func (r *productRepositoryImpl) FindByIDWithDetails(ctx context.Context, id string) (*model.Product, error) {
	return r.findByIDBase(ctx, id, "Categories", "Tags", "Variants", "Variants.Color", "Variants.Size", "Variants.Inventory", "Images", "Images.Color")
}

func (r *productRepositoryImpl) ExistsByID(ctx context.Context, id string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Product, error) {
	return r.findByIDBase(ctx, id)
}

func (r *productRepositoryImpl) FindAllByCategorySlug(ctx context.Context, categorySlug string) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Joins("JOIN product_categories pc ON pc.product_id = products.id").Joins("JOIN categories c ON c.id = pc.category_id").Where("c.slug = ? AND products.is_deleted = false", categorySlug).Preload("Categories").Preload("Variants").Preload("Variants.Color").Preload("Variants.Size").Preload("Variants.Inventory").Preload("Images").Preload("Images.Color").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) FindAllWithCategoriesAndThumbnail(ctx context.Context) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Preload("Categories").Preload("Images", "is_thumbnail = ?", true).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) UpdateCategories(ctx context.Context, product *model.Product, categories []*model.Category) error {
	if err := r.db.WithContext(ctx).Model(product).Association("Categories").Replace(categories); err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindByIDWithCategoriesAndTags(ctx context.Context, id string) (*model.Product, error) {
	return r.findByIDBase(ctx, id, "Categories", "Tags")
}

func (r *productRepositoryImpl) UpdateTags(ctx context.Context, product *model.Product, tags []*model.Tag) error {
	if err := r.db.WithContext(ctx).Model(product).Association("Tags").Replace(tags); err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return customErr.ErrProductNotFound
	}

	return nil
}

func (r *productRepositoryImpl) FindAllByID(ctx context.Context, ids []string) ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.WithContext(ctx).Scopes(notDeleted).Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepositoryImpl) UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error {
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Where("id IN ?", ids).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) findByIDBase(ctx context.Context, id string, preloads ...string) (*model.Product, error) {
	var product model.Product
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Scopes(notDeleted).Where("id = ?", id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepositoryImpl) findBySlugBase(ctx context.Context, slug string, preloads ...string) (*model.Product, error) {
	var product model.Product
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if err := query.Scopes(notDeleted).Where("slug = ?", slug).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func notDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = false")
}
