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

func (h *GRPCHandler) CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*protobuf.CreatedResponse, error) {
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

	return &protobuf.CreatedResponse{
		Id: category.ID,
	}, nil
}

func (h *GRPCHandler) GetCategoryTree(ctx context.Context, req *protobuf.GetCategoryTreeRequest) (*protobuf.CategoryTreeResponse, error) {
	categoryTree, err := h.svc.GetCategoryTree(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toCategoryTreeResponse(categoryTree), nil
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*protobuf.CreatedResponse, error) {
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

	return &protobuf.CreatedResponse{
		Id: product.ID,
	}, nil
}

func (h *GRPCHandler) GetProductBySlug(ctx context.Context, req *protobuf.GetProductBySlugRequest) (*protobuf.ProductPublicResponse, error) {
	product, err := h.svc.GetProductBySlug(ctx, req.Slug)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toProductPublicResponse(product), nil
}

func (h *GRPCHandler) CreateColor(ctx context.Context, req *protobuf.CreateColorRequest) (*protobuf.CreatedResponse, error) {
	color, err := h.svc.CreateColor(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrColorAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: color.ID,
	}, nil
}

func (h *GRPCHandler) CreateSize(ctx context.Context, req *protobuf.CreateSizeRequest) (*protobuf.CreatedResponse, error) {
	size, err := h.svc.CreateSize(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSizeAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: size.ID,
	}, nil
}

func (h *GRPCHandler) CreateVariant(ctx context.Context, req *protobuf.CreateVariantRequest) (*protobuf.CreatedResponse, error) {
	variant, err := h.svc.CreateVariant(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSKUAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrProductNotFound, customErr.ErrColorNotFound, customErr.ErrSizeNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: variant.ID,
	}, nil
}

func (h *GRPCHandler) CreateImage(ctx context.Context, req *protobuf.CreateImageRequest) (*protobuf.CreatedResponse, error) {
	image, err := h.svc.CreateImage(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound, customErr.ErrColorNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: image.ID,
	}, nil
}

func (h *GRPCHandler) GetProductsByCategory(ctx context.Context, req *protobuf.GetProductsByCategoryRequest) (*protobuf.ProductsPublicResponse, error) {
	products, err := h.svc.GetProductsByCategory(ctx, req.Slug)
	if err != nil {
		switch err {
		case customErr.ErrCategoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return toProductsPublicResponse(products), nil
}

func (h *GRPCHandler) CreateTag(ctx context.Context, req *protobuf.CreateTagRequest) (*protobuf.CreatedResponse, error) {
	tag, err := h.svc.CreateTag(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrTagAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: tag.ID,
	}, nil
}

func (h *GRPCHandler) GetAllCategories(ctx context.Context, req *protobuf.GetAllCategoriesRequest) (*protobuf.CategoriesAdminResponse, error) {
	categories, err := h.svc.GetAllCategories(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return categories, nil
}

func toProductsPublicResponse(products []*model.Product) *protobuf.ProductsPublicResponse {
	var productResponses []*protobuf.ProductPublicResponse
	for _, pro := range products {
		productResponses = append(productResponses, toProductPublicResponse(pro))
	}

	return &protobuf.ProductsPublicResponse{
		Products: productResponses,
	}
}

func toProductPublicResponse(product *model.Product) *protobuf.ProductPublicResponse {
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

	variants := make([]*protobuf.BaseVariantResponse, len(product.Variants))
	for i, variant := range product.Variants {
		variants[i] = toBaseVariantResponse(variant)
	}

	images := make([]*protobuf.BaseImageResponse, len(product.Images))
	for i, image := range product.Images {
		images[i] = toBaseImageResponse(image)
	}

	return &protobuf.ProductPublicResponse{
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
		Variants:    variants,
		Images:      images,
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

func toBaseCategoryResponse(category *model.Category) *protobuf.BaseCategoryResponse {
	return &protobuf.BaseCategoryResponse{
		Id:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func toBaseVariantResponse(variant *model.Variant) *protobuf.BaseVariantResponse {
	return &protobuf.BaseVariantResponse{
		Id: variant.ID,
		Color: &protobuf.BaseColorResponse{
			Id:   variant.ColorID,
			Name: variant.Color.Name,
		},
		Size: &protobuf.BaseSizeResponse{
			Id:   variant.SizeID,
			Name: variant.Size.Name,
		},
		Inventory: &protobuf.BaseInventoryResponse{
			Id:           variant.Inventory.ID,
			SoldQuantity: int64(variant.Inventory.SoldQuantity),
			Stock:        int64(variant.Inventory.Stock),
			IsStock:      &variant.Inventory.IsStock,
		},
	}
}

func toBaseImageResponse(image *model.Image) *protobuf.BaseImageResponse {
	return &protobuf.BaseImageResponse{
		Id: image.ID,
		Color: &protobuf.BaseColorResponse{
			Id:   image.ColorID,
			Name: image.Color.Name,
		},
		Url:         image.Url,
		SortOrder:   int32(image.SortOrder),
		IsThumbnail: image.IsThumbnail,
	}
}
