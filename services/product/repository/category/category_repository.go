package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	FindByID(ctx context.Context, id string) (*model.Category, error)

	FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Category, error)

	ExistsByID(ctx context.Context, id string) (bool, error)

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	FindAll(ctx context.Context) ([]*model.Category, error)
}