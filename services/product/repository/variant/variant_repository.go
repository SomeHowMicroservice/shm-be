package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
)

type VariantRepository interface {
	Create(ctx context.Context, variant *model.Variant) error

	CreateAll(ctx context.Context, variants []*model.Variant) error

	ExistsBySKU(ctx context.Context, sku string) (bool, error)

	FindAllByIDIn(ctx context.Context, ids []string) ([]*model.Variant, error)

	UpdateIsDeletedByIDIn(ctx context.Context, ids []string) error

	Update(ctx context.Context, id string, updateData map[string]interface{}) error
}