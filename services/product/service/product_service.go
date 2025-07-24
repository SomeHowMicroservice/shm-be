package service

import (
	"context"

	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
)

type ProductService interface {
	CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*model.Category, error)
}