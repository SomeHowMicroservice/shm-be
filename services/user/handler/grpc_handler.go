package handler

import (
	"github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/user/service"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	protobuf.UnimplementedUserServiceServer
	svc service.UserService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.UserService) *grpcHandler {
	return &grpcHandler{
		svc: svc,
		UnimplementedUserServiceServer: protobuf.UnimplementedUserServiceServer{},
	}
}