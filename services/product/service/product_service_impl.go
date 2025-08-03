package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/imagekit"
	"github.com/SomeHowMicroservice/shm-be/services/product/model"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	colorRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/color"
	imageRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/image"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	sizeRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/size"
	tagRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/tag"
	variantRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/variant"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productServiceImpl struct {
	userClient   userpb.UserServiceClient
	imageKitSvc  imagekit.ImageKitService
	categoryRepo categoryRepo.CategoryRepository
	productRepo  productRepo.ProductRepository
	tagRepo      tagRepo.TagRepository
	colorRepo    colorRepo.ColorRepository
	sizeRepo     sizeRepo.SizeRepository
	variantRepo  variantRepo.VariantRepository
	imageRepo    imageRepo.ImageRepository
}

func NewProductService(userClient userpb.UserServiceClient, imageKitSvc imagekit.ImageKitService, categoryRepo categoryRepo.CategoryRepository, productRepo productRepo.ProductRepository, tagRepo tagRepo.TagRepository, colorRepo colorRepo.ColorRepository, sizeRepo sizeRepo.SizeRepository, variantRepo variantRepo.VariantRepository, imageRepo imageRepo.ImageRepository) ProductService {
	return &productServiceImpl{
		userClient,
		imageKitSvc,
		categoryRepo,
		productRepo,
		tagRepo,
		colorRepo,
		sizeRepo,
		variantRepo,
		imageRepo,
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
		if len(parents) != len(req.ParentIds) {
			return nil, customErr.ErrHasCategoryNotFound
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
	categories, err := s.categoryRepo.FindAllWithParentsAndChildren(ctx)
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
		categories, err = s.categoryRepo.FindAllByIDInWithChildren(ctx, req.CategoryIds)
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
		ID:          uuid.NewString(),
		SKU:         req.Sku,
		ProductID:   req.ProductId,
		ColorID:     req.ColorId,
		SizeID:      req.SizeId,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
		Inventory: &model.Inventory{
			ID:           uuid.NewString(),
			Quantity:     int(req.Quantity),
			SoldQuantity: 0,
			UpdatedByID:  req.UserId,
		},
	}
	variant.Inventory.SetStock()
	if err = s.variantRepo.Create(ctx, variant); err != nil {
		return nil, fmt.Errorf("tạo biến thể sản phẩm thất bại: %w", err)
	}

	return variant, nil
}

func (s *productServiceImpl) CreateImage(ctx context.Context, req *protobuf.CreateImageRequest) (*model.Image, error) {
	product, err := s.productRepo.FindByID(ctx, req.ProductId)
	if err != nil {
		return nil, fmt.Errorf("tìm thông tin sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	exists, err := s.colorRepo.ExistsByID(ctx, req.ColorId)
	if err != nil {
		return nil, fmt.Errorf("tìm thông tin màu sắc thất bại: %w", err)
	}
	if !exists {
		return nil, customErr.ErrColorNotFound
	}

	ext := strings.ToLower(filepath.Ext(req.FileName))
	if ext == "" {
		ext = ".jpg"
	}

	fileName := fmt.Sprintf("%s-%d%s", product.Slug, req.SortOrder, ext)

	uploadFileRequest := &common.UploadFileRequest{
		File:     bytes.NewReader(req.File),
		FileName: fileName,
		Folder:   "somehow_microservice/product",
	}
	uploadedRes, err := s.imageKitSvc.UploadFile(ctx, uploadFileRequest)
	if err != nil {
		return nil, err
	}

	image := &model.Image{
		ID:          uuid.NewString(),
		ProductID:   req.ProductId,
		ColorID:     req.ColorId,
		Url:         uploadedRes.URL,
		IsThumbnail: req.IsThumbnail,
		SortOrder:   int(req.SortOrder),
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err = s.imageRepo.Create(ctx, image); err != nil {
		return nil, fmt.Errorf("tạo ảnh sản phẩm thất bại: %w", err)
	}

	return image, nil
}

func (s *productServiceImpl) GetProductsByCategory(ctx context.Context, categorySlug string) ([]*model.Product, error) {
	exists, err := s.categoryRepo.ExistsBySlug(ctx, categorySlug)
	if err != nil {
		return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm thất bại: %w", err)
	}
	if !exists {
		return nil, customErr.ErrCategoryNotFound
	}

	products, err := s.productRepo.FindAllByCategorySlug(ctx, categorySlug)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách sản phẩm theo danh mục thất bại: %w", err)
	}

	return products, nil
}

func (s *productServiceImpl) CreateTag(ctx context.Context, req *protobuf.CreateTagRequest) (*model.Tag, error) {
	slug := common.GenerateSlug(req.Name)
	exists, err := s.tagRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("kiểm tra tag sản phẩm tồn tại thất bại: %w", err)
	}
	if exists {
		return nil, customErr.ErrTagAlreadyExists
	}

	tag := &model.Tag{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Slug:        slug,
		CreatedByID: req.UserId,
		UpdatedByID: req.UserId,
	}
	if err = s.tagRepo.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("tạo nhãn sản phẩm thất bại: %w", err)
	}

	return tag, nil
}

func (s *productServiceImpl) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	categories, err := s.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy tất cả danh mục sản phẩm thất bại: %w", err)
	}

	return categories, nil
}

func (s *productServiceImpl) GetCategoryByID(ctx context.Context, id string) (*protobuf.CategoryAdminDetailsResponse, error) {
	category, err := s.categoryRepo.FindByIDWithParentsAndProducts(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	cRes, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{
		Id: category.CreatedByID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	uRes, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{
		Id: category.CreatedByID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	productResponses := toBaseProductResponse(category)

	return toCategoryAdminDetailsResponse(category, productResponses, cRes, uRes), nil
}

func (s *productServiceImpl) UpdateCategory(ctx context.Context, req *protobuf.UpdateCategoryRequest) (*protobuf.CategoryAdminDetailsResponse, error) {
	category, err := s.categoryRepo.FindByIDWithParents(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	updateData := map[string]interface{}{}
	if category.Name != req.Name {
		updateData["name"] = req.Name
	}
	if category.Slug != req.Slug {
		updateData["slug"] = req.Slug
	}
	if category.UpdatedByID != req.UserId {
		updateData["updated_by_id"] = req.UserId
	}

	if len(updateData) > 0 {
		if err = s.categoryRepo.Update(ctx, category.ID, updateData); err != nil {
			if errors.Is(err, customErr.ErrCategoryNotFound) {
				return nil, err
			}
			return nil, fmt.Errorf("cập nhật danh mục sản phẩm thất bại: %w", err)
		}
	}

	parentIDs := getParentIDsFromParents(category.Parents)

	if !slices.Equal(parentIDs, req.ParentIds) {
		parents, err := s.categoryRepo.FindAllByIDIn(ctx, req.ParentIds)
		if err != nil {
			return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm cha thất bại: %w", err)
		}
		if len(parents) != len(req.ParentIds) {
			return nil, customErr.ErrHasCategoryNotFound
		}

		if err = s.categoryRepo.UpdateParents(ctx, category, parents); err != nil {
			return nil, fmt.Errorf("cập nhật danh mục cha thất bại: %w", err)
		}
	}

	category, err = s.categoryRepo.FindByIDWithParentsAndProducts(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("tìm kiếm danh mục sản phẩm thất bại: %w", err)
	}
	if category == nil {
		return nil, customErr.ErrCategoryNotFound
	}

	cRes, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{
		Id: category.CreatedByID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	uRes, err := s.userClient.GetUserById(ctx, &userpb.GetUserByIdRequest{
		Id: category.CreatedByID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	productResponses := toBaseProductResponse(category)

	return toCategoryAdminDetailsResponse(category, productResponses, cRes, uRes), nil
}

func (s *productServiceImpl) GetAllColors(ctx context.Context) (*protobuf.ColorsAdminResponse, error) {
	colors, err := s.colorRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách màu sắc sản phẩm thất bại: %w", err)
	}

	userIDMap := map[string]struct{}{}
	for _, color := range colors {
		userIDMap[color.CreatedByID] = struct{}{}
		userIDMap[color.UpdatedByID] = struct{}{}
	}

	var userIDs []string
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	userRes, err := s.userClient.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{
		Ids: userIDs,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	userMap := make(map[string]*userpb.UserPublicResponse)
	for _, user := range userRes.Users {
		userMap[user.Id] = user
	}

	var colorResponses []*protobuf.ColorAdminResponse
	for _, color := range colors {
		colorResponses = append(colorResponses, &protobuf.ColorAdminResponse{
			Id:   color.ID,
			Name: color.Name,
			CreatedAt: color.CreatedAt.Format(time.RFC3339),
			CreatedBy: &protobuf.BaseUserResponse{
				Id:       color.CreatedByID,
				Username: userMap[color.CreatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[color.CreatedByID].Profile.Id,
					FirstName: userMap[color.CreatedByID].Profile.FirstName,
					LastName:  userMap[color.CreatedByID].Profile.LastName,
				},
			},
			UpdatedAt: color.UpdatedAt.Format(time.RFC3339),
			UpdatedBy: &protobuf.BaseUserResponse{
				Id:       color.UpdatedByID,
				Username: userMap[color.UpdatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[color.UpdatedByID].Profile.Id,
					FirstName: userMap[color.UpdatedByID].Profile.FirstName,
					LastName:  userMap[color.UpdatedByID].Profile.LastName,
				},
			},
		})
	}

	return &protobuf.ColorsAdminResponse{
		Colors: colorResponses,
	}, nil
}

func (s *productServiceImpl) GetAllSizes(ctx context.Context) (*protobuf.SizesAdminResponse, error) {
	sizes, err := s.sizeRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("lấy danh sách size sản phẩm thất bại: %w", err)
	}

	userIDMap := map[string]struct{}{}
	for _, size := range sizes {
		userIDMap[size.CreatedByID] = struct{}{}
		userIDMap[size.UpdatedByID] = struct{}{}
	}

	var userIDs []string
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	userRes, err := s.userClient.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{
		Ids: userIDs,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	userMap := make(map[string]*userpb.UserPublicResponse)
	for _, user := range userRes.Users {
		userMap[user.Id] = user
	}

	var sizeResponses []*protobuf.SizeAdminResponse
	for _, size := range sizes {
		sizeResponses = append(sizeResponses, &protobuf.SizeAdminResponse{
			Id:   size.ID,
			Name: size.Name,
			CreatedAt: size.CreatedAt.Format(time.RFC3339),
			CreatedBy: &protobuf.BaseUserResponse{
				Id:       size.CreatedByID,
				Username: userMap[size.CreatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[size.CreatedByID].Profile.Id,
					FirstName: userMap[size.CreatedByID].Profile.FirstName,
					LastName:  userMap[size.CreatedByID].Profile.LastName,
				},
			},
			UpdatedAt: size.UpdatedAt.Format(time.RFC3339),
			UpdatedBy: &protobuf.BaseUserResponse{
				Id:       size.UpdatedByID,
				Username: userMap[size.UpdatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[size.UpdatedByID].Profile.Id,
					FirstName: userMap[size.UpdatedByID].Profile.FirstName,
					LastName:  userMap[size.UpdatedByID].Profile.LastName,
				},
			},
		})
	}

	return &protobuf.SizesAdminResponse{
		Sizes: sizeResponses,
	}, nil
}

func (s *productServiceImpl) GetAllTags(ctx context.Context) (*protobuf.TagsAdminResponse, error) {
	tags, err := s.tagRepo.FindAll(ctx)
	userIDMap := map[string]struct{}{}
	for _, tag := range tags {
		userIDMap[tag.CreatedByID] = struct{}{}
		userIDMap[tag.UpdatedByID] = struct{}{}
	}

	var userIDs []string
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	userRes, err := s.userClient.GetUsersByIds(ctx, &userpb.GetUsersByIdsRequest{
		Ids: userIDs,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrHasUserNotFound
			default:
				return nil, fmt.Errorf("lỗi từ user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi không xác định: %w", err)
	}

	userMap := make(map[string]*userpb.UserPublicResponse)
	for _, user := range userRes.Users {
		userMap[user.Id] = user
	}

	var tagResponses []*protobuf.TagAdminResponse
	for _, tag := range tags {
		tagResponses = append(tagResponses, &protobuf.TagAdminResponse{
			Id:   tag.ID,
			Name: tag.Name,
			CreatedAt: tag.CreatedAt.Format(time.RFC3339),
			CreatedBy: &protobuf.BaseUserResponse{
				Id:       tag.CreatedByID,
				Username: userMap[tag.CreatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[tag.CreatedByID].Profile.Id,
					FirstName: userMap[tag.CreatedByID].Profile.FirstName,
					LastName:  userMap[tag.CreatedByID].Profile.LastName,
				},
			},
			UpdatedAt: tag.UpdatedAt.Format(time.RFC3339),
			UpdatedBy: &protobuf.BaseUserResponse{
				Id:       tag.UpdatedByID,
				Username: userMap[tag.UpdatedByID].Username,
				Profile: &protobuf.BaseProfileResponse{
					Id:        userMap[tag.UpdatedByID].Profile.Id,
					FirstName: userMap[tag.UpdatedByID].Profile.FirstName,
					LastName:  userMap[tag.UpdatedByID].Profile.LastName,
				},
			},
		})
	}

	return &protobuf.TagsAdminResponse{
		Tags: tagResponses,
	}, nil
}

func (s *productServiceImpl) UpdateTag(ctx context.Context, req *protobuf.UpdateTagRequest) error {
	tag, err := s.tagRepo.FindByID(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("tìm kiếm tag sản phẩm thất bại: %w", err)
	}
	if tag == nil {
		return customErr.ErrTagNotFound
	}

	updateData := map[string]interface{}{}
	if tag.Name != req.Name {
		slug := common.GenerateSlug(req.Name)
		exists, err := s.tagRepo.ExistsBySlug(ctx, slug)
		if err != nil {
			return fmt.Errorf("kiểm tra tồn tại slug thất bại: %w", err)
		}
		if exists {
			return customErr.ErrTagAlreadyExists
		}
		updateData["name"] = req.Name
		updateData["slug"] = slug
	}
	if tag.UpdatedByID != req.UserId {
		updateData["updated_by_id"] = req.UserId
	}

	if len(updateData) > 0 {
		if err = s.tagRepo.Update(ctx, tag.ID, updateData); err != nil {
			if errors.Is(err, customErr.ErrTagNotFound) {
				return err
			}
			return fmt.Errorf("cập nhật tag sản phẩm thất bại: %w", err)
		}
	}

	return nil
}

func getParentIDsFromParents(categories []*model.Category) []string {
	var parentIDs []string
	for _, cat := range categories {
		parentIDs = append(parentIDs, cat.ID)
	}

	return parentIDs
}

func toCategoryAdminDetailsResponse(category *model.Category, productResponses []*protobuf.BaseProductResponse, cRes *userpb.UserResponse, uRes *userpb.UserResponse) *protobuf.CategoryAdminDetailsResponse {
	return &protobuf.CategoryAdminDetailsResponse{
		Id:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		Parents:   toBaseCategoriesResponse(category.Parents),
		CreatedAt: category.CreatedAt.Format(time.RFC3339),
		UpdatedAt: category.UpdatedAt.Format(time.RFC3339),
		CreatedBy: &protobuf.BaseUserResponse{
			Id:       cRes.Id,
			Username: cRes.Username,
			Profile: &protobuf.BaseProfileResponse{
				Id:        cRes.Profile.Id,
				FirstName: cRes.Profile.FirstName,
				LastName:  cRes.Profile.LastName,
			},
		},
		UpdatedBy: &protobuf.BaseUserResponse{
			Id:       uRes.Id,
			Username: uRes.Username,
			Profile: &protobuf.BaseProfileResponse{
				Id:        uRes.Profile.Id,
				FirstName: uRes.Profile.FirstName,
				LastName:  uRes.Profile.LastName,
			},
		},
		Products: productResponses,
	}
}

func toBaseProductResponse(category *model.Category) []*protobuf.BaseProductResponse {
	var productResponses []*protobuf.BaseProductResponse
	for _, p := range category.Products {
		var thumb *protobuf.BaseImageResponse
		for _, img := range p.Images {
			if img.IsThumbnail {
				thumb = &protobuf.BaseImageResponse{
					Id:          img.ID,
					Url:         img.Url,
					IsThumbnail: img.IsThumbnail,
				}
				break
			}
		}

		productResponses = append(productResponses, &protobuf.BaseProductResponse{
			Id:    p.ID,
			Title: p.Title,
			Slug:  p.Slug,
			Image: thumb,
		})
	}

	return productResponses
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
