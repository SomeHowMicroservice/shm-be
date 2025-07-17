package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *model.Profile) error

	UpdateByUserID(ctx context.Context, userID string, updateData map[string]interface{}) error
}