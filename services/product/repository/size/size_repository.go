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

	FindByID(ctx context.Context, id string) (*model.Size, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.Size, error)

	UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error

	FindAllDeleted(ctx context.Context) ([]*model.Size, error)
}