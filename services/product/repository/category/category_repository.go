package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	FindByID(ctx context.Context, id string) (*model.Category, error)

	FindAllByID(ctx context.Context, ids []string) ([]*model.Category, error)

	FindAllByIDWithChildren(ctx context.Context, ids []string) ([]*model.Category, error)

	ExistsByID(ctx context.Context, id string) (bool, error)

	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	FindAllWithParentsAndChildren(ctx context.Context) ([]*model.Category, error)

	FindAll(ctx context.Context) ([]*model.Category, error)

	FindByIDWithParentsAndProducts(ctx context.Context, id string) (*model.Category, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindByIDWithParents(ctx context.Context, id string) (*model.Category, error)

	UpdateParents(ctx context.Context, category *model.Category, parents []*model.Category) error

	FindAllWithChildren(ctx context.Context) ([]*model.Category, error)

	DeleteAllByID(ctx context.Context, ids []string) error

	Delete(ctx context.Context, id string) error
}