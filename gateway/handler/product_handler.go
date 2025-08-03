package handler

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"

	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	productClient productpb.ProductServiceClient
}

func NewProductHandler(productClient productpb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{productClient}
}

func (h *ProductHandler) CreateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var slug *string
	if req.Slug != nil {
		slug = req.Slug
	}

	res, err := h.productClient.CreateCategory(ctx, &productpb.CreateCategoryRequest{
		Name:      req.Name,
		Slug:      slug,
		ParentIds: req.ParentIDs,
		UserId:    user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo danh mục sản phẩm thành công", gin.H{
		"category_id": res.Id,
	})
}

func (h *ProductHandler) GetCategoryTree(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetCategoryTree(ctx, &productpb.GetCategoryTreeRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh mục sản phẩm thành công", gin.H{
		"categories": res.Categories,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var isSale bool
	var salePrice *float32
	var startSale, endSale *string
	if req.IsSale != nil {
		isSale = *req.IsSale
		if isSale {
			if req.SalePrice == nil || req.StartSale == nil || req.EndSale == nil {
				common.JSON(c, http.StatusBadRequest, "Sản phẩm giảm giá phải bổ sung thêm thông tin", nil)
				return
			}
			salePrice = req.SalePrice

			formattedStartSale := req.StartSale.Format("2006-01-02")
			startSale = &formattedStartSale
			formattedEndSale := req.EndSale.Format("2006-01-02")
			endSale = &formattedEndSale
		} else {
			if req.SalePrice != nil || req.StartSale != nil || req.EndSale != nil {
				common.JSON(c, http.StatusBadRequest, "Sản phẩm không được giảm giá vui lòng không điền thông tin liên quan", nil)
				return
			}
		}
	}

	res, err := h.productClient.CreateProduct(ctx, &productpb.CreateProductRequest{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		IsSale:      isSale,
		SalePrice:   salePrice,
		StartSale:   startSale,
		EndSale:     endSale,
		CategoryIds: req.CategoryIDs,
		UserId:      user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo sản phẩm thành công", gin.H{
		"product_id": res.Id,
	})
}

func (h *ProductHandler) ProductDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productSlug := c.Param("slug")

	res, err := h.productClient.GetProductBySlug(ctx, &productpb.GetProductBySlugRequest{
		Slug: productSlug,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy sản phẩm thành công", gin.H{
		"product": res,
	})
}

func (h *ProductHandler) CreateColor(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.productClient.CreateColor(ctx, &productpb.CreateColorRequest{
		Name:   req.Name,
		UserId: user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo màu sắc thành công", gin.H{
		"color_id": res.Id,
	})
}

func (h *ProductHandler) CreateSize(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.productClient.CreateSize(ctx, &productpb.CreateSizeRequest{
		Name:   req.Name,
		UserId: user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo size thành công", gin.H{
		"size_id": res.Id,
	})
}

func (h *ProductHandler) CreateVariant(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.productClient.CreateVariant(ctx, &productpb.CreateVariantRequest{
		Sku:       req.SKU,
		ProductId: req.ProductID,
		ColorId:   req.ColorID,
		SizeId:    req.SizeID,
		Quantity:  int64(req.Quantity),
		UserId:    user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo biến thể sản phẩm thành công", gin.H{
		"variant_id": res.Id,
	})
}

func (h *ProductHandler) CreateImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateImageRequest
	if err := c.ShouldBind(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	fileHeader := req.File

	maxSize := int64(10 * 1024 * 1024)
	if fileHeader.Size > maxSize {
		common.JSON(c, http.StatusRequestEntityTooLarge, "File phải có kích thước bé hơn hoặc bằng 10MB", nil)
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()

	data, _ := io.ReadAll(file)

	var isThumbnail bool
	if req.IsThumbnail != nil {
		isThumbnail = *req.IsThumbnail
	}

	res, err := h.productClient.CreateImage(ctx, &productpb.CreateImageRequest{
		ProductId:   req.ProductID,
		ColorId:     req.ColorID,
		File:        data,
		FileName:    fileHeader.Filename,
		IsThumbnail: isThumbnail,
		SortOrder:   int32(req.SortOrder),
		UserId:      user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Thêm hình ảnh sản phẩm thành công", gin.H{
		"image_id": res.Id,
	})
}

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categorySlug := c.Param("slug")

	res, err := h.productClient.GetProductsByCategory(ctx, &productpb.GetProductsByCategoryRequest{
		Slug: categorySlug,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách sản phẩm thành công", gin.H{
		"products": res.Products,
	})
}

func (h *ProductHandler) CreateTag(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.productClient.CreateTag(ctx, &productpb.CreateTagRequest{
		Name:   req.Name,
		UserId: user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo tag sản phẩm thành công", gin.H{
		"tag_id": res.Id,
	})
}

func (h *ProductHandler) GetAllCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllCategories(ctx, &productpb.GetAllCategoriesRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả danh mục sản phẩm thành công", gin.H{
		"categories": res.Categories,
	})
}

func (h *ProductHandler) CategoryAdminDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categoryID := c.Param("id")

	res, err := h.productClient.GetCategoryById(ctx, &productpb.GetCategoryByIdRequest{
		Id: categoryID,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh mục sản phẩm thành công", gin.H{
		"category": res,
	})
}

func (h *ProductHandler) UpdateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*userpb.UserPublicResponse)
	if !ok {
		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	categoryID := c.Param("id")

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	res, err := h.productClient.UpdateCategory(ctx, &productpb.UpdateCategoryRequest{
		Id:        categoryID,
		Name:      req.Name,
		Slug:      req.Slug,
		ParentIds: req.ParentIDs,
		UserId:    user.Id,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật danh mục sản phẩm thành công", gin.H{
		"category": res,
	})
}

func (h *ProductHandler) GetAllColors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllColors(ctx, &productpb.GetAllColorsRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả màu sắc sản phẩm thành công", gin.H{
		"colors": res.Colors,
	})
}

func (h *ProductHandler) GetAllSizes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllSizes(ctx, &productpb.GetAllSizesRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả size sản phẩm thành công", gin.H{
		"sizes": res.Sizes,
	})
}
