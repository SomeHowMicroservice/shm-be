package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error 

	ExistsByUsername(ctx context.Context, username string) (bool, error)

	ExistsByEmail(ctx context.Context, email string) (bool, error)

	FindByUsername(ctx context.Context, username string) (*model.User, error)

	FindById(ctx context.Context, id string) (*model.User, error)

	FindByEmail(ctx context.Context, email string) (*model.User, error)

	UpdatePassword(ctx context.Context, id, password string) error
}