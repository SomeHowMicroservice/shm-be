package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserService interface {
	CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*model.User, error)
}
