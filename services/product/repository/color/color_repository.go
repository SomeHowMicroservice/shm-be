package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ColorRepository interface {
	Create(ctx context.Context, color *model.Color) error

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	ExistsByID(ctx context.Context, id string) (bool, error)
}