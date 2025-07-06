package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserService interface {
	CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserResponse, error)
}
