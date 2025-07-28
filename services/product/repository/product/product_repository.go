package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	FindBySlug(ctx context.Context, slug string) (*model.Product, error)

	ExistsByID(ctx context.Context, id string) (bool, error)
}