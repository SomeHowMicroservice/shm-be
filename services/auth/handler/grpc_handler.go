package handler

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	protobuf.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.AuthService) *grpcHandler {
	return &grpcHandler{
		svc: svc,
	}
}

func (h *grpcHandler) SignUp(ctx context.Context, req *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	token, err := h.svc.SignUp(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameAlreadyExists, customErr.ErrEmailAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Lỗi auth service")
		}
	}

	return &protobuf.SignUpResponse{
		Message:           "Đăng ký thành công, vui lòng kiểm tra email để lấy mã OTP",
		RegistrationToken: token,
	}, nil
}
