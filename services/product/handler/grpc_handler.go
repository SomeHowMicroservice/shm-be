package handler

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/product/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	protobuf.UnimplementedProductServiceServer
	svc service.ProductService
}

func NewGRPCHandler(grpcServer *grpc.Server, svc service.ProductService) *GRPCHandler {
	return &GRPCHandler{svc: svc}
}

func (h *GRPCHandler) CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*protobuf.CategoryResponse, error) {
	category, err := h.svc.CreateCategory(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSlugAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrCategoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	// Convert []bson.ObjectID to []string for ParentIds and ChildrenIds
	parentIds := make([]string, len(category.ParentIDs))
	for i, id := range category.ParentIDs {
		parentIds[i] = id.Hex()
	}
	childrenIds := make([]string, len(category.ChildrenIDs))
	for i, id := range category.ChildrenIDs {
		childrenIds[i] = id.Hex()
	}

	return &protobuf.CategoryResponse{
		Id: category.ID.Hex(),
		Name: category.Name,
		Slug: category.Slug,
		ParentIds: parentIds,
		ChildrenIds: childrenIds,
	}, nil
}