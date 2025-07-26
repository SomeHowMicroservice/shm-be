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

func (h *GRPCHandler) CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*protobuf.CategoryAdminResponse, error) {
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

	return toCategoryAdminResponse(category), nil
}

func (h *GRPCHandler) GetCategoryTree(ctx context.Context, req *protobuf.GetCategoryTreeRequest) (*protobuf.CategoryTreeResponse, error) {
	categoryTree, err := h.svc.GetCategoryTree(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toCategoryTreeResponse(categoryTree), nil
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*protobuf.ProductAdminResponse, error) {
	product, err := h.svc.CreateProduct(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSlugAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrCategoryNotFound, customErr.ErrHasCategoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toProductAdminResponse(product), nil
}

func toProductAdminResponse(product *model.Product) *protobuf.ProductAdminResponse {
	var startSalePtr, endSalePtr *string
	if product.StartSale != nil {
		formatted := product.StartSale.Format("2006-01-02")
		startSalePtr = &formatted
	}
	if product.EndSale != nil {
		formatted := product.EndSale.Format("2006-01-02")
		endSalePtr = &formatted
	}

	categories := make([]*protobuf.BaseCategoryResponse, len(product.Categories))
	for i, category := range product.Categories {
		categories[i] = toBaseCategoryResponse(category)
	}

	return &protobuf.ProductAdminResponse{
		Id:          product.ID,
		Title:       product.Title,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		IsSale:      &product.IsSale,
		SalePrice:   product.SalePrice,
		StartSale:   startSalePtr,
		EndSale:     endSalePtr,
		Categories:  categories,
	}
}

func toCategoryPublicResponse(category *model.Category) *protobuf.CategoryPublicResponse {
	children := make([]*protobuf.CategoryPublicResponse, 0, len(category.Children))
	for _, child := range category.Children {
		children = append(children, toCategoryPublicResponse(child))
	}

	return &protobuf.CategoryPublicResponse{
		Id:       category.ID,
		Name:     category.Name,
		Slug:     category.Slug,
		Children: children,
	}
}

func toCategoryTreeResponse(categories []*model.Category) *protobuf.CategoryTreeResponse {
	result := make([]*protobuf.CategoryPublicResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryPublicResponse(c))
	}

	return &protobuf.CategoryTreeResponse{
		Categories: result,
	}
}

func toCategoryAdminResponse(category *model.Category) *protobuf.CategoryAdminResponse {
	parents := make([]*protobuf.BaseCategoryResponse, len(category.Parents))
	for i, parent := range category.Parents {
		parents[i] = toBaseCategoryResponse(parent)
	}
	return &protobuf.CategoryAdminResponse{
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
