package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/user/model"
)

type AddressRepository interface {
	Create(ctx context.Context, address *model.Address) error

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindByID(ctx context.Context, id string) (*model.Address, error)

	FindByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CountByUserID(ctx context.Context, userID string) (int64, error)

	Delete(ctx context.Context, id string) error
}
