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
}
