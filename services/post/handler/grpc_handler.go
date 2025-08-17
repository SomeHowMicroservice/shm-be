package handler

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/post/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	protobuf.UnimplementedPostServiceServer
	svc service.PostService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.PostService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) CreateTopic(ctx context.Context, req *protobuf.CreateTopicRequest) (*protobuf.CreatedResponse, error) {
	topicID, err := h.svc.CreateTopic(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrTopicAlreadyExists:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: topicID,
	}, nil
}
