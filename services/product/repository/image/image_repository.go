package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) error

	CreateAll(ctx context.Context, images []*model.Image) error

	UpdateIsDeletedByIDIn(ctx context.Context, ids []string) error

	FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Image, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error
}
