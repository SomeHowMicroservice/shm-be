package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
)

type ProductService interface {
	CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*model.Category, error)

	GetCategoryTree(ctx context.Context) ([]*model.Category, error)

	GetCategoriesNoChild(ctx context.Context) ([]*model.Category, error)

	GetCategoriesNoProduct(ctx context.Context) ([]*model.Category, error)

	GetProductBySlug(ctx context.Context, productSlug string) (*model.Product, error)

	CreateColor(ctx context.Context, req *protobuf.CreateColorRequest) (*model.Color, error)

	CreateSize(ctx context.Context, req *protobuf.CreateSizeRequest) (*model.Size, error)

	GetProductsByCategory(ctx context.Context, categorySlug string) ([]*model.Product, error)

	CreateTag(ctx context.Context, req *protobuf.CreateTagRequest) (*model.Tag, error)

	GetAllCategories(ctx context.Context) ([]*model.Category, error)

	GetCategoryByID(ctx context.Context, categoryID string) (*protobuf.CategoryAdminDetailsResponse, error)

	UpdateCategory(ctx context.Context, req *protobuf.UpdateCategoryRequest) (*protobuf.CategoryAdminDetailsResponse, error)

	GetAllColorsAdmin(ctx context.Context) (*protobuf.ColorsAdminResponse, error)

	GetAllSizesAdmin(ctx context.Context) (*protobuf.SizesAdminResponse, error)

	GetAllTagsAdmin(ctx context.Context) (*protobuf.TagsAdminResponse, error)

	UpdateTag(ctx context.Context, req *protobuf.UpdateTagRequest) error

	GetAllColors(ctx context.Context) ([]*model.Color, error)

	GetAllSizes(ctx context.Context) ([]*model.Size, error)

	GetAllTags(ctx context.Context) ([]*model.Tag, error)

	CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*model.Product, error)

	GetProductByID(ctx context.Context, productID string) (*protobuf.ProductAdminDetailsResponse, error)

	GetAllProductsAdmin(ctx context.Context) ([]*model.Product, error)

	UpdateProduct(ctx context.Context, req *protobuf.UpdateProductRequest) (*protobuf.ProductAdminDetailsResponse, error)

	DeleteProduct(ctx context.Context, req *protobuf.DeleteOneRequest) error

	DeleteProducts(ctx context.Context, req *protobuf.DeleteManyRequest) error

	PermanentlyDeleteCategory(ctx context.Context, req *protobuf.DeleteOneRequest) error

	PermanentlyDeleteCategories(ctx context.Context, req *protobuf.DeleteManyRequest) error

	UpdateColor(ctx context.Context, req *protobuf.UpdateColorRequest) error

	UpdateSize(ctx context.Context, req *protobuf.UpdateSizeRequest) error

	DeleteColor(ctx context.Context, req *protobuf.DeleteOneRequest) error

	DeleteSize(ctx context.Context, req *protobuf.DeleteOneRequest) error

	DeleteColors(ctx context.Context, req *protobuf.DeleteManyRequest) error

	DeleteSizes(ctx context.Context, req *protobuf.DeleteManyRequest) error

	GetDeletedProducts(ctx context.Context) ([]*model.Product, error)

	GetDeletedProductByID(ctx context.Context, productID string) (*protobuf.ProductAdminDetailsResponse, error)

	GetDeletedColors(ctx context.Context) (*protobuf.ColorsAdminResponse, error)

	GetDeletedSizes(ctx context.Context) (*protobuf.SizesAdminResponse, error)

	GetDeletedTags(ctx context.Context) (*protobuf.TagsAdminResponse, error)

	DeleteTag(ctx context.Context, req *protobuf.DeleteOneRequest) error

	DeleteTags(ctx context.Context, req *protobuf.DeleteManyRequest) error

	RestoreProduct(ctx context.Context, req *protobuf.RestoreOneRequest) error

	RestoreProducts(ctx context.Context, req *protobuf.RestoreManyRequest) error

	RestoreColor(ctx context.Context, req *protobuf.RestoreOneRequest) error

	RestoreColors(ctx context.Context, req *protobuf.RestoreManyRequest) error

	RestoreSize(ctx context.Context, req *protobuf.RestoreOneRequest) error

	RestoreSizes(ctx context.Context, req *protobuf.RestoreManyRequest) error

	RestoreTag(ctx context.Context, req *protobuf.RestoreOneRequest) error

	RestoreTags(ctx context.Context, req *protobuf.RestoreManyRequest) error
}
