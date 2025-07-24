package service

import (
	"context"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type productServiceImpl struct {
	categoryRepo categoryRepo.CategoryRepository
}

func NewProductService(categoryRepo categoryRepo.CategoryRepository) ProductService {
	return &productServiceImpl{categoryRepo}
}

func (s *productServiceImpl) CreateCategory(ctx context.Context, req *protobuf.CreateCategoryRequest) (*model.Category, error) {
	if req.Slug == nil {
		slug := common.GenerateSlug(req.Name)
		req.Slug = &slug
	}

	exists, err := s.categoryRepo.ExistsBySlug(ctx, *req.Slug)
	if err != nil {
		return nil, fmt.Errorf("kiếm tra tồn tại slug thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrSlugAlreadyExists
	}

	var parentIDs []bson.ObjectID
	if len(req.ParentIds) > 0 {
		parentIDs = make([]bson.ObjectID, 0, len(req.ParentIds))
		for _, parentIDStr := range req.ParentIds {
			parentID, err := bson.ObjectIDFromHex(parentIDStr)
			if err != nil {
				return nil, fmt.Errorf("định dạng parent_id không chính xác: %w", err)
			}

			exists, err := s.categoryRepo.ExistsByID(ctx, parentID)
			if err != nil {
				return nil, fmt.Errorf("kiểm tra danh mục sản phẩm tồn tại thất bại:%w", err)
			}
			if !exists {
				return nil, customErr.ErrCategoryNotFound
			}

			parentIDs = append(parentIDs, parentID)
		}
	}

	category := &model.Category{
		ID: bson.NewObjectID(),
		Name: req.Name,
		Slug: *req.Slug,
		ParentIDs: parentIDs,		
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("tạo danh mục sản phẩm thất bại: %w", err)
	}

	return category, nil
}
