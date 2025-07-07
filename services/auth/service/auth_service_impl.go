package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/google/uuid"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

type authServiceImpl struct {
	userClient userpb.UserServiceClient
}

func NewAuthService(userClient userpb.UserServiceClient) AuthService {
	return &authServiceImpl{
		userClient: userClient,
	}
}

func (s *authServiceImpl) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (string, error) {
	_, err := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return "", err
	}

	registrationToken := uuid.NewString()

	return registrationToken, nil
}
