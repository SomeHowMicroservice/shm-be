package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthService interface {
	SignUp(ctx context.Context, req *protobuf.SignUpRequest) (string, error)

	VerifySignUp(ctx context.Context, req *protobuf.VerifySignUpRequest) (*protobuf.AuthResponse, string, int64, string, int64, error)
}
