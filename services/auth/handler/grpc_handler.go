package handler

import (
	"context"

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
		if st, ok := status.FromError(err); ok {
			return nil, status.Error(st.Code(), st.Message())
		}
		return nil, status.Error(codes.Internal, "Lỗi auth service")
	}
	
	return &protobuf.SignUpResponse{
		Message:           "Đăng ký thành công",
		RegistrationToken: token,
	}, nil
}
