package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ColorRepository interface {
	Create(ctx context.Context, color *model.Color) error

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	ExistsByID(ctx context.Context, id string) (bool, error)

	FindAll(ctx context.Context) ([]*model.Color, error)

	FindAllDeleted(ctx context.Context) ([]*model.Color, error)

	FindByID(ctx context.Context, id string) (*model.Color, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.Color, error)

	UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error

	FindAllDeletedByID(ctx context.Context, ids []string) ([]*model.Color, error)

	FindDeletedByID(ctx context.Context, id string) (*model.Color, error)
}