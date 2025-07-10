package handler

import (
	"context"
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
	return &grpcHandler{svc: svc}
}

func (h *grpcHandler) CheckEmailExists(ctx context.Context, req *protobuf.CheckEmailExistsRequest) (*protobuf.CheckEmailExistsResponse, error) {
	exists, err := h.svc.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &protobuf.CheckEmailExistsResponse{
		Exists: exists,
	}, nil
}

func (h *grpcHandler) CheckUsernameExists(ctx context.Context, req *protobuf.CheckUsernameExistsRequest) (*protobuf.CheckUsernameExistsResponse, error) {
	exists, err := h.svc.CheckUsernameExists(ctx, req.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &protobuf.CheckUsernameExistsResponse{
		Exists: exists,
	}, nil
}

func (h *grpcHandler) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserResponse, error) {
	user, err := h.svc.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameAlreadyExists, customErr.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (h *grpcHandler) GetUserByUsername(ctx context.Context, req *protobuf.GetUserByUsernameRequest) (*protobuf.UserResponse, error) {
	user, err := h.svc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &protobuf.UserResponse{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}, nil
}
