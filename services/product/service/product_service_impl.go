package service

import (
	"context"
	"fmt"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	"github.com/google/uuid"
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

	var parents []*model.Category
	if len(req.ParentIds) > 0 {
		parents, err = s.categoryRepo.FindAllByIDIn(ctx, req.ParentIds)
		if err != nil {
			return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm cha thất bại: %w", err)
		}
	}

	category := &model.Category{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Slug:        *req.Slug,
		Parents:     parents,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err = s.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("tạo danh mục sản phẩm thất bại: %w", err)
	}

	return category, nil
}

func (s *productServiceImpl) GetCategoryTree(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả danh mục sản phẩm thất bại: %w", err)
	}

	categoryMap := make(map[string]*model.Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	var roots []*model.Category
	for _, cat := range categories {
		if len(cat.Parents) == 0 {
			roots = append(roots, cat)
		}
	}

	return roots, nil
}
