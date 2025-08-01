package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type TagRepository interface {
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	Create(ctx context.Context, tag *model.Tag) error
}