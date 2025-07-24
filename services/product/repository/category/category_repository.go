package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	FindByID(ctx context.Context, id bson.ObjectID) (*model.Category, error)

	ExistsByID(ctx context.Context, id bson.ObjectID) (bool, error)

	ExistsBySlug(ctx context.Context, slug string) (bool, error)
}