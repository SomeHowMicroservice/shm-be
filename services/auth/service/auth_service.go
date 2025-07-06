package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthService interface {
	SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error)
}