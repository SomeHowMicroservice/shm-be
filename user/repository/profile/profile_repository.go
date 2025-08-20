package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/user/model"
)

type ProfileRepository interface {
	Update(ctx context.Context, id string, updateData map[string]interface{}) error

	FindByID(ctx context.Context, id string) (*model.Profile, error)
}