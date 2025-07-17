package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserService interface {
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	CheckUsernameExists(ctx context.Context, username string) (bool, error)

	CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*model.User, error)

	UpdateUserPassword(ctx context.Context, req *protobuf.UpdateUserPasswordRequest) error

	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	GetUserById(ctx context.Context, id string) (*model.User, error)

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	UpdateUserProfile(ctx context.Context, req *protobuf.UpdateUserProfileRequest) (*model.User, error)
}
