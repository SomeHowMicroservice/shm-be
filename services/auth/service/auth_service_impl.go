package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authServiceImpl struct {
	protobuf.UnimplementedAuthServiceServer
	userClient userpb.UserServiceClient
}

func NewAuthService(userClient userpb.UserServiceClient) AuthService {
	return &authServiceImpl{
		userClient: userClient,
		UnimplementedAuthServiceServer: protobuf.UnimplementedAuthServiceServer{},
	}
}

func (s *authServiceImpl) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	_, err := s.userClient.CreateUser(ctx, &userpb.CreateUserRequest{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Tạo user thất bại: %v", err)
	}

	registrationToken := uuid.NewString()

	return &protobuf.SignUpResponse{
		Message: "Đăng ký thành công",
		RegistrationToken: registrationToken,
	}, nil
}