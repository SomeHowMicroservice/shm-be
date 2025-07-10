package repository

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
)

type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*model.Role, error)

	CreateUserRoles(ctx context.Context, userID string, roleID string) error
}