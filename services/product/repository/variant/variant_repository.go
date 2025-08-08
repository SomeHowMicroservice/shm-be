package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type VariantRepository interface {
	Create(ctx context.Context, variant *model.Variant) error

	CreateAll(ctx context.Context, variants []*model.Variant) error

	ExistsBySKU(ctx context.Context, sku string) (bool, error)

	FindAllByID(ctx context.Context, ids []string) ([]*model.Variant, error)

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	DeleteAllByID(ctx context.Context, ids []string) error
}