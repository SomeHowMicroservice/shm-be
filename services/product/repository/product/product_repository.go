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

	FindAllByCategorySlug(ctx context.Context, categorySlug string) ([]*model.Product, error)

	FindAllWithCategoriesAndThumbnail(ctx context.Context) ([]*model.Product, error)
}