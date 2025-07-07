package handler

import (
	"context"
	"errors"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	protobuf.UnimplementedUserServiceServer
	svc service.UserService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.UserService) *grpcHandler {
	return &grpcHandler{
		svc: svc,
	}
}

func (h *grpcHandler) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserResponse, error) {
	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, customErr.ErrUsernameAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, customErr.ErrUsernameAlreadyExists.Error())
		case errors.Is(err, customErr.ErrEmailAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, customErr.ErrEmailAlreadyExists.Error())
		default:
			return nil, status.Error(codes.Internal, "Lỗi user service")
		}
	}

	return &protobuf.UserResponse{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
