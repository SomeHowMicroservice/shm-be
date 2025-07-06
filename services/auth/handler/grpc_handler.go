package handler

import (
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/service"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	protobuf.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.AuthService) *grpcHandler {
	return &grpcHandler{
		svc: svc,
		UnimplementedAuthServiceServer: protobuf.UnimplementedAuthServiceServer{},
	}
}