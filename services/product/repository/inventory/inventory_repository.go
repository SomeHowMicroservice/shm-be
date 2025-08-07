package repository

import "context"

type InventoryRepository interface {
	UpdateByVariantID(ctx context.Context, variantID string, updateData map[string]interface{}) error
}