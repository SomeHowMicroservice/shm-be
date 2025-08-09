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
	"google.golang.org/protobuf/proto"
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

func (h *GRPCHandler) GetCategoryTree(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.CategoryTreeResponse, error) {
	categoryTree, err := h.svc.GetCategoryTree(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toCategoryTreeResponse(categoryTree), nil
}

func (h *GRPCHandler) GetCategoriesNoProduct(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.BaseCategoriesResponse, error) {
	categories, err := h.svc.GetCategoriesNoProduct(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &protobuf.BaseCategoriesResponse{
		Categories: toBaseCategoriesResponse(categories),
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

func (h *GRPCHandler) GetAllCategoriesAdmin(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.BaseCategoriesResponse, error) {
	categories, err := h.svc.GetAllCategories(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.BaseCategoriesResponse{
		Categories: toBaseCategoriesResponse(categories),
	}, nil
}

func (h *GRPCHandler) GetCategoryById(ctx context.Context, req *protobuf.GetCategoryByIdRequest) (*protobuf.CategoryAdminDetailsResponse, error) {
	convertedCategory, err := h.svc.GetCategoryByID(ctx, req.Id)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedCategory, nil
}

func (h *GRPCHandler) UpdateCategory(ctx context.Context, req *protobuf.UpdateCategoryRequest) (*protobuf.CategoryAdminDetailsResponse, error) {
	convertedCategory, err := h.svc.UpdateCategory(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound, customErr.ErrHasCategoryNotFound, customErr.ErrCategoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedCategory, nil
}

func (h *GRPCHandler) GetAllColorsAdmin(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.ColorsAdminResponse, error) {
	convertedColors, err := h.svc.GetAllColorsAdmin(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedColors, nil
}

func (h *GRPCHandler) GetAllColors(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.ColorsPublicResponse, error) {
	colors, err := h.svc.GetAllColors(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toColorsPublicResponse(colors), nil
}

func (h *GRPCHandler) GetAllSizesAdmin(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.SizesAdminResponse, error) {
	convertedSizes, err := h.svc.GetAllSizesAdmin(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedSizes, nil
}

func (h *GRPCHandler) GetAllSizes(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.SizesPublicResponse, error) {
	sizes, err := h.svc.GetAllSizes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toSizesPublicResponse(sizes), nil
}

func (h *GRPCHandler) GetAllTagsAdmin(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.TagsAdminResponse, error) {
	convertedTags, err := h.svc.GetAllTagsAdmin(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedTags, nil
}

func (h *GRPCHandler) GetAllTags(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.TagsPublicResponse, error) {
	tags, err := h.svc.GetAllTags(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var baseTags []*protobuf.BaseTagResponse
	for _, tag := range tags {
		baseTags = append(baseTags, &protobuf.BaseTagResponse{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}

	return &protobuf.TagsPublicResponse{
		Tags: baseTags,
	}, nil
}

func (h *GRPCHandler) UpdateTag(ctx context.Context, req *protobuf.UpdateTagRequest) (*protobuf.UpdatedResponse, error) {
	if err := h.svc.UpdateTag(ctx, req); err != nil {
		switch err {
		case customErr.ErrTagNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrTagAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.UpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*protobuf.CreatedResponse, error) {
	product, err := h.svc.CreateProduct(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSlugAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrHasCategoryNotFound, customErr.ErrHasTagNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.CreatedResponse{
		Id: product.ID,
	}, nil
}

func (h *GRPCHandler) GetCategoriesNoChild(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.BaseCategoriesResponse, error) {
	categories, err := h.svc.GetCategoriesNoChild(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.BaseCategoriesResponse{
		Categories: toBaseCategoriesResponse(categories),
	}, nil
}

func (h *GRPCHandler) GetProductById(ctx context.Context, req *protobuf.GetProductByIdRequest) (*protobuf.ProductAdminDetailsResponse, error) {
	convertedProduct, err := h.svc.GetProductByID(ctx, req.Id)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedProduct, nil
}

func (h *GRPCHandler) GetAllProductsAdmin(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.ProductsAdminResponse, error) {
	products, err := h.svc.GetAllProductsAdmin(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProductsAdminResponse(products), nil
}

func (h *GRPCHandler) UpdateProduct(ctx context.Context, req *protobuf.UpdateProductRequest) (*protobuf.ProductAdminDetailsResponse, error) {
	convertedProduct, err := h.svc.UpdateProduct(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrSlugAlreadyExists, customErr.ErrSKUAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case customErr.ErrHasCategoryNotFound, customErr.ErrHasTagNotFound, customErr.ErrHasImageNotFound, customErr.ErrHasVariantNotFound, customErr.ErrProductNotFound, customErr.ErrVariantNotFound, customErr.ErrImageNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedProduct, nil
}

func (h *GRPCHandler) DeleteProduct(ctx context.Context, req *protobuf.DeleteOneRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteProduct(ctx, req); err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) DeleteProducts(ctx context.Context, req *protobuf.DeleteManyRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteProducts(ctx, req); err != nil {
		switch err {
		case customErr.ErrHasProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) PermanentlyDeleteCategory(ctx context.Context, req *protobuf.DeleteOneRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.PermanentlyDeleteCategory(ctx, req); err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) PermanentlyDeleteCategories(ctx context.Context, req *protobuf.DeleteManyRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.PermanentlyDeleteCategories(ctx, req); err != nil {
		switch err {
		case customErr.ErrHasCategoryNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) UpdateColor(ctx context.Context, req *protobuf.UpdateColorRequest) (*protobuf.UpdatedResponse, error) {
	if err := h.svc.UpdateColor(ctx, req); err != nil {
		switch err {
		case customErr.ErrColorNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrColorAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.UpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) UpdateSize(ctx context.Context, req *protobuf.UpdateSizeRequest) (*protobuf.UpdatedResponse, error) {
	if err := h.svc.UpdateSize(ctx, req); err != nil {
		switch err {
		case customErr.ErrSizeNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case customErr.ErrSizeAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.UpdatedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) DeleteColor(ctx context.Context, req *protobuf.DeleteOneRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteColor(ctx, req); err != nil {
		switch err {
		case customErr.ErrColorNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) DeleteSize(ctx context.Context, req *protobuf.DeleteOneRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteSize(ctx, req); err != nil {
		switch err {
		case customErr.ErrSizeNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) DeleteColors(ctx context.Context, req *protobuf.DeleteManyRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteColors(ctx, req); err != nil {
		switch err {
		case customErr.ErrHasColorNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) DeleteSizes(ctx context.Context, req *protobuf.DeleteManyRequest) (*protobuf.DeletedResponse, error) {
	if err := h.svc.DeleteSizes(ctx, req); err != nil {
		switch err {
		case customErr.ErrHasSizeNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &protobuf.DeletedResponse{
		Success: true,
	}, nil
}

func (h *GRPCHandler) GetDeletedProducts(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.ProductsAdminResponse, error) {
	products, err := h.svc.GetDeletedProducts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProductsAdminResponse(products), nil
}

func (h *GRPCHandler) GetDeletedProductById(ctx context.Context, req *protobuf.GetProductByIdRequest) (*protobuf.ProductAdminDetailsResponse, error) {
	convertedProduct, err := h.svc.GetDeletedProductByID(ctx, req.Id)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedProduct, nil
}

func (h *GRPCHandler) GetDeletedColors(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.ColorsAdminResponse, error) {
	convertedColors, err := h.svc.GetDeletedColors(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedColors, nil
}

func (h *GRPCHandler) GetDeletedSizes(ctx context.Context, req *protobuf.GetManyRequest) (*protobuf.SizesAdminResponse, error) {
	convertedSizes, err := h.svc.GetDeletedSizes(ctx)
	if err != nil {
		switch err {
		case customErr.ErrHasUserNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return convertedSizes, nil
}

func toProductsAdminResponse(products []*model.Product) *protobuf.ProductsAdminResponse {
	var productResponses []*protobuf.ProductAdminResponse
	for _, pro := range products {
		productResponses = append(productResponses, toProductAdminResponse(pro))
	}

	return &protobuf.ProductsAdminResponse{
		Products: productResponses,
	}
}

func toProductAdminResponse(product *model.Product) *protobuf.ProductAdminResponse {
	categories := toBaseCategoriesResponse(product.Categories)
	thumbnail := &protobuf.SimpleImageResponse{
		Id:  product.Images[0].ID,
		Url: product.Images[0].Url,
	}
	return &protobuf.ProductAdminResponse{
		Id:         product.ID,
		Title:      product.Title,
		Price:      product.Price,
		Categories: categories,
		Thumbnail:  thumbnail,
	}
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

func toBaseCategoriesResponse(categories []*model.Category) []*protobuf.BaseCategoryResponse {
	var baseCategories []*protobuf.BaseCategoryResponse
	for _, category := range categories {
		baseCategories = append(baseCategories, toBaseCategoryResponse(category))
	}
	return baseCategories
}

func toBaseCategoryResponse(category *model.Category) *protobuf.BaseCategoryResponse {
	return &protobuf.BaseCategoryResponse{
		Id:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func toColorsPublicResponse(colors []*model.Color) *protobuf.ColorsPublicResponse {
	var baseColors []*protobuf.BaseColorResponse
	for _, color := range colors {
		baseColors = append(baseColors, toBaseColorResponse(color))
	}
	return &protobuf.ColorsPublicResponse{
		Colors: baseColors,
	}
}

func toBaseColorResponse(color *model.Color) *protobuf.BaseColorResponse {
	return &protobuf.BaseColorResponse{
		Id:   color.ID,
		Name: color.Name,
	}
}

func toSizesPublicResponse(sizes []*model.Size) *protobuf.SizesPublicResponse {
	var baseSizes []*protobuf.BaseSizeResponse
	for _, size := range sizes {
		baseSizes = append(baseSizes, toBaseSizeResponse(size))
	}
	return &protobuf.SizesPublicResponse{
		Sizes: baseSizes,
	}
}

func toBaseSizeResponse(size *model.Size) *protobuf.BaseSizeResponse {
	return &protobuf.BaseSizeResponse{
		Id:   size.ID,
		Name: size.Name,
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
			SoldQuantity: proto.Int64(int64(variant.Inventory.SoldQuantity)),
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
		IsThumbnail: &image.IsThumbnail,
	}
}
