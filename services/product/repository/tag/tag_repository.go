package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type TagRepository interface {
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	Create(ctx context.Context, tag *model.Tag) error

	FindAll(ctx context.Context) ([]*model.Tag, error)

	FindByID(ctx context.Context, id string) (*model.Tag, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	UpdateAllByID(ctx context.Context, ids []string, updateData map[string]interface{}) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.Tag, error)

	FindAllDeleted(ctx context.Context) ([]*model.Tag, error)

	FindDeletedByID(ctx context.Context, id string) (*model.Tag, error)

	FindAllDeletedByID(ctx context.Context, ids []string) ([]*model.Tag, error)
}
