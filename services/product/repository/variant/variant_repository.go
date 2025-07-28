package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type VariantRepository interface {
	Create(ctx context.Context, variant *model.Variant) error

	ExistsBySKU(ctx context.Context, sku string) (bool, error)
}