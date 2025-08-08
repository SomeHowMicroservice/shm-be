package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) error

	CreateAll(ctx context.Context, images []*model.Image) error

	FindAllByID(ctx context.Context, ids []string) ([]*model.Image, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	DeleteAllByID(ctx context.Context, ids []string) error

	UpdateFileID(ctx context.Context, id string, fileID string) error
}
