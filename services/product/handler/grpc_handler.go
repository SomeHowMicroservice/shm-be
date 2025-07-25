package handler

import (
	"context"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
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

func (h *GRPCHandler) CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*protobuf.NewCategoryResponse, error) {
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

	return toNewCategoryResponse(category), nil
}

func (h *GRPCHandler) GetCategoryTree(ctx context.Context, req *protobuf.GetCategoryTreeRequest) (*protobuf.CategoryTreeResponse, error) {
	categoryTree, err := h.svc.GetCategoryTree(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toCategoryTreeResponse(categoryTree), nil
}

func toCategoryResponse(category *model.Category) *protobuf.CategoryResponse {
	children := make([]*protobuf.CategoryResponse, 0, len(category.Children))
	for _, child := range category.Children {
		children = append(children, toCategoryResponse(child))
	}

	return &protobuf.CategoryResponse{
		Id:       category.ID,
		Name:     category.Name,
		Slug:     category.Slug,
		Children: children,
	}
}

func toCategoryTreeResponse(categories []*model.Category) *protobuf.CategoryTreeResponse {
	result := make([]*protobuf.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(c))
	}

	return &protobuf.CategoryTreeResponse{
		Categories: result,
	}
}

func toNewCategoryResponse(category *model.Category) *protobuf.NewCategoryResponse {
	parents := make([]*protobuf.BaseCategoryResponse, len(category.Parents))
	for i, child := range category.Children {
		parents[i] = toBaseCategoryResponse(child)
	}
	return &protobuf.NewCategoryResponse{
		Id:      category.ID,
		Name:    category.Name,
		Slug:    category.Slug,
		Parents: parents,
	}
}

func toBaseCategoryResponse(category *model.Category) *protobuf.BaseCategoryResponse {
	return &protobuf.BaseCategoryResponse{
		Id:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}
