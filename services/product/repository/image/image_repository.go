package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) error
}