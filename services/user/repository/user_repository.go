package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error 

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)
}