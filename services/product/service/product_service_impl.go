package service

import (
	"context"
	"fmt"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	colorRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/color"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	sizeRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/size"
	variantRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/variant"
	"github.com/google/uuid"
)

type productServiceImpl struct {
	categoryRepo categoryRepo.CategoryRepository
	productRepo  productRepo.ProductRepository
	colorRepo    colorRepo.ColorRepository
	sizeRepo     sizeRepo.SizeRepository
	variantRepo  variantRepo.VariantRepository
}

func NewProductService(categoryRepo categoryRepo.CategoryRepository, productRepo productRepo.ProductRepository, colorRepo colorRepo.ColorRepository, sizeRepo sizeRepo.SizeRepository, variantRepo variantRepo.VariantRepository) ProductService {
	return &productServiceImpl{
		categoryRepo,
		productRepo,
		colorRepo,
		sizeRepo,
		variantRepo,
	}
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

	catMap := make(map[string]*model.Category)
	for _, c := range categories {
		c.Children = nil
		catMap[c.ID] = c
	}

	var roots []*model.Category
	for _, c := range categories {
		if len(c.Parents) == 0 {
			roots = append(roots, c)
		} else {
			for _, p := range c.Parents {
				if parent, ok := catMap[p.ID]; ok {
					parent.Children = append(parent.Children, c)
				}
			}
		}
	}

	return roots, nil
}

func (s *productServiceImpl) CreateProduct(ctx context.Context, req *protobuf.CreateProductRequest) (*model.Product, error) {
	slug := common.GenerateSlug(req.Title)
	exists, err := s.productRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra tồn tại slug thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrSlugAlreadyExists
	}

	var categories []*model.Category
	if len(req.CategoryIds) > 0 {
		categories, err = s.categoryRepo.FindAllByIDIn(ctx, req.CategoryIds)
		if err != nil {
			return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm thất bại: %w", err)
		}
		if len(categories) != len(req.CategoryIds) {
			return nil, customErr.ErrHasCategoryNotFound
		}

		for _, c := range categories {
			if len(c.Children) > 0 {
				return nil, fmt.Errorf("danh mục %s có danh mục con, không thể được gán cho sản phẩm", c.Name)
			}
		}
	}

	var startSale, endSale *time.Time
	if req.StartSale != nil && req.EndSale != nil {
		parsedStartSale, err := common.ParseDate(*req.StartSale)
		if err != nil {
			return nil, fmt.Errorf("chuyển đổi kiểu dữ liệu thời gian bắt đầu giảm giá thất bại: %w", err)
		}
		startSale = &parsedStartSale

		parsedEndSale, err := common.ParseDate(*req.EndSale)
		if err != nil {
			return nil, fmt.Errorf("chuyển đổi kiểu dữ liệu thời gian kết thúc giảm giá thất bại: %w", err)
		}
		endSale = &parsedEndSale
	}

	product := &model.Product{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Price:       req.Price,
		IsSale:      req.IsSale,
		SalePrice:   req.SalePrice,
		StartSale:   startSale,
		EndSale:     endSale,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
		Categories:  categories,
	}
	if err = s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("tạo sản phẩm thất bại: %w", err)
	}

	return product, nil
}

func (s *productServiceImpl) GetProductBySlug(ctx context.Context, slug string) (*model.Product, error) {
	product, err := s.productRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("lấy sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	return product, nil
}

func (s *productServiceImpl) CreateColor(ctx context.Context, req *protobuf.CreateColorRequest) (*model.Color, error) {
	slug := common.GenerateSlug(req.Name)
	exists, err := s.colorRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra màu tồn tại thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrColorAlreadyExists
	}

	color := &model.Color{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Slug:        slug,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err = s.colorRepo.Create(ctx, color); err != nil {
		return nil, fmt.Errorf("tạo màu sắc thất bại: %w", err)
	}

	return color, nil
}

func (s *productServiceImpl) CreateSize(ctx context.Context, req *protobuf.CreateSizeRequest) (*model.Size, error) {
	slug := common.GenerateSlug(req.Name)
	exists, err := s.sizeRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra size tồn tại thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrSizeAlreadyExists
	}

	size := &model.Size{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Slug:        slug,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err = s.sizeRepo.Create(ctx, size); err != nil {
		return nil, fmt.Errorf("tạo size thất bại: %w", err)
	}

	return size, nil
}

func (s *productServiceImpl) CreateVariant(ctx context.Context, req *protobuf.CreateVariantRequest) (*model.Variant, error) {
	exists, err := s.variantRepo.ExistsBySKU(ctx, req.Sku)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra tồn tại mã SKU thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrSKUAlreadyExists
	}

	exists, err = s.productRepo.ExistsByID(ctx, req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("tìm thông tin sản phẩm thất bại: %w", err)
	}
	if !exists {
		return nil, customErr.ErrProductNotFound
	}

	exists, err = s.colorRepo.ExistsByID(ctx, req.ColorId)
	if err != nil {
		return nil, fmt.Errorf("tìm thông tin màu sắc thất bại: %w", err)
	}
	if !exists {
		return nil, customErr.ErrColorNotFound
	}

	exists, err = s.sizeRepo.ExistsByID(ctx, req.SizeId)
	if err != nil {
		return nil, fmt.Errorf("tìm thông tin kích cỡ thất bại: %w", err)
	}
	if !exists {
		return nil, customErr.ErrSizeNotFound
	}

	variant := &model.Variant{
		ID: uuid.NewString(),
		SKU: req.Sku,
		ProductID: req.ProductId,
		ColorID: req.ColorId,
		SizeID: req.SizeId,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
		Inventory: &model.Inventory{
			ID: uuid.NewString(),
			Quantity: int(req.Quantity),
			SoldQuantity: 0,
			UpdatedByID: req.UserId,
		},
	}
	variant.Inventory.SetStock()
	if err = s.variantRepo.Create(ctx, variant); err != nil {
		return nil, fmt.Errorf("tạo biến thể sản phẩm thất bại: %w", err)
	}

	return variant, nil
}
