package service

import (
	"context"
	"time"

	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthService interface {
	SignUp(ctx context.Context, req *protobuf.SignUpRequest) (string, error)

	VerifySignUp(ctx context.Context, req *protobuf.VerifySignUpRequest) (*protobuf.AuthResponse, string, time.Duration, string, time.Duration, error)

	SignIn(ctx context.Context, req *protobuf.SignInRequest) (*protobuf.AuthResponse, string, time.Duration, string, time.Duration, error)

	RefreshToken(ctx context.Context, req *protobuf.RefreshTokenRequest) (string, time.Duration, string, time.Duration, error)
}
