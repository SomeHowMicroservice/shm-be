package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
)

type ProductService interface {
	CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*model.Category, error)

	GetCategoryTree(ctx context.Context) ([]*model.Category, error) 

	CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*model.Product, error)

	GetProductBySlug(ctx context.Context, slug string) (*model.Product, error)

	CreateColor(ctx context.Context, req *protobuf.CreateColorRequest) (*model.Color, error)

	CreateSize(ctx context.Context, req *protobuf.CreateSizeRequest) (*model.Size, error)

	CreateVariant(ctx context.Context, req *protobuf.CreateVariantRequest) (*model.Variant, error)

	CreateImage(ctx context.Context, req *protobuf.CreateImageRequest) (*model.Image, error)

	GetProductsByCategory(ctx context.Context, categorySlug string) ([]*model.Product, error)

	CreateTag(ctx context.Context, req *protobuf.CreateTagRequest) (*model.Tag, error)

	GetAllCategories(ctx context.Context) ([]*model.Category, error)

	GetCategoryByID(ctx context.Context, id string) (*protobuf.CategoryAdminDetailsResponse, error)

	UpdateCategory(ctx context.Context, req *protobuf.UpdateCategoryRequest) (*protobuf.CategoryAdminDetailsResponse, error)

	GetAllColors(ctx context.Context) (*protobuf.ColorsAdminResponse, error)

	GetAllSizes(ctx context.Context) (*protobuf.SizesAdminResponse, error)

	GetAllTags(ctx context.Context) (*protobuf.TagsAdminResponse, error)

	UpdateTag(ctx context.Context, req *protobuf.UpdateTagRequest) error
}