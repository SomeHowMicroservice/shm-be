package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	FindBySlugWithDetails(ctx context.Context, slug string) (*model.Product, error)

	FindByIDWithDetails(ctx context.Context, id string) (*model.Product, error)

	ExistsByID(ctx context.Context, id string) (bool, error)

	FindByID(ctx context.Context, id string) (*model.Product, error)

	FindByIDWithCategoriesAndTags(ctx context.Context, id string) (*model.Product, error)

	FindAllByCategorySlug(ctx context.Context, categorySlug string) ([]*model.Product, error)

	FindAllWithCategoriesAndThumbnail(ctx context.Context) ([]*model.Product, error)

	UpdateCategories(ctx context.Context, product *model.Product, categories []*model.Category) error

	UpdateTags(ctx context.Context, product *model.Product, tags []*model.Tag) error

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.Product, error)

	UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error
}
