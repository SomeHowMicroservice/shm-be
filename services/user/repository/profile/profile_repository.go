package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error

	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindByID(ctx context.Context, id string) (*model.Profile, error)
}