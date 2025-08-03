package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type SizeRepository interface {
	Create(ctx context.Context, size *model.Size) error

	FindAll(ctx context.Context) ([]*model.Size, error)

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	ExistsByID(ctx context.Context, id string) (bool, error)
}