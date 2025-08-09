package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	"github.com/go-playground/validator/v10"

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

	res, err := h.productClient.GetCategoryTree(ctx, &productpb.GetManyRequest{})
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

func (h *ProductHandler) GetCategoriesNoChild(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetCategoriesNoChild(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh mục sản phẩm không có con thành công", gin.H{
		"categories": res.Categories,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
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

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		common.JSON(c, http.StatusBadRequest, "Không thể parse form", nil)
		return
	}

	var req request.CreateProductForm
	form := c.Request.MultipartForm

	req.Title = strings.TrimSpace(c.PostForm("title"))
	req.Description = strings.TrimSpace(c.PostForm("description"))

	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 32); err == nil {
			req.Price = float32(price)
		}
	}

	if isSaleStr := c.PostForm("is_sale"); isSaleStr != "" {
		if isSale, err := strconv.ParseBool(isSaleStr); err == nil {
			req.IsSale = &isSale
		}
	}

	if salePriceStr := c.PostForm("sale_price"); salePriceStr != "" {
		if salePrice, err := strconv.ParseFloat(salePriceStr, 32); err == nil {
			salePriceFloat := float32(salePrice)
			req.SalePrice = &salePriceFloat
		}
	}

	if startSaleStr := c.PostForm("start_sale"); startSaleStr != "" {
		if startSale, err := time.Parse("2006-01-02", startSaleStr); err == nil {
			req.StartSale = &startSale
		}
	}

	if endSaleStr := c.PostForm("end_sale"); endSaleStr != "" {
		if endSale, err := time.Parse("2006-01-02", endSaleStr); err == nil {
			req.EndSale = &endSale
		}
	}

	req.CategoryIDs = form.Value["category_ids"]
	req.TagIDs = form.Value["tag_ids"]

	req.Variants = []request.CreateVariantForm{}
	i := 0
	for {
		skuKey := fmt.Sprintf("variants[%d][sku]", i)
		colorKey := fmt.Sprintf("variants[%d][color_id]", i)
		sizeKey := fmt.Sprintf("variants[%d][size_id]", i)
		quantityKey := fmt.Sprintf("variants[%d][quantity]", i)

		sku := strings.TrimSpace(c.PostForm(skuKey))

		if sku == "" {
			break
		}

		colorID := strings.TrimSpace(c.PostForm(colorKey))
		sizeID := strings.TrimSpace(c.PostForm(sizeKey))
		quantityStr := strings.TrimSpace(c.PostForm(quantityKey))

		quantity := 0
		if quantityStr != "" {
			quantity, _ = strconv.Atoi(quantityStr)
		}

		variant := request.CreateVariantForm{
			SKU:      sku,
			ColorID:  colorID,
			SizeID:   sizeID,
			Quantity: quantity,
		}

		req.Variants = append(req.Variants, variant)
		i++
	}

	req.Images = []request.CreateImageForm{}
	j := 0
	for {
		colorIDKey := fmt.Sprintf("images[%d][color_id]", j)
		isThumbnailKey := fmt.Sprintf("images[%d][is_thumbnail]", j)
		sortOrderKey := fmt.Sprintf("images[%d][sort_order]", j)
		fileKey := fmt.Sprintf("images[%d][file]", j)

		colorID := strings.TrimSpace(c.PostForm(colorIDKey))
		if colorID == "" {
			break
		}

		isThumbnailStr := strings.TrimSpace(c.PostForm(isThumbnailKey))
		sortOrderStr := strings.TrimSpace(c.PostForm(sortOrderKey))

		isThumbnail := false
		if isThumbnailStr != "" {
			isThumbnail, _ = strconv.ParseBool(isThumbnailStr)
		}

		sortOrder := 0
		if sortOrderStr != "" {
			sortOrder, _ = strconv.Atoi(sortOrderStr)
		}

		file, err := c.FormFile(fileKey)
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Không tìm thấy file cho image %d: %s", j, err.Error()), nil)
			return
		}

		image := request.CreateImageForm{
			ColorID:     colorID,
			IsThumbnail: &isThumbnail,
			SortOrder:   sortOrder,
			File:        file,
		}

		req.Images = append(req.Images, image)
		j++
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
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

	variants := make([]*productpb.CreateVariantRequest, 0, len(req.Variants))
	for _, v := range req.Variants {
		variants = append(variants, &productpb.CreateVariantRequest{
			Sku:      v.SKU,
			ColorId:  v.ColorID,
			SizeId:   v.SizeID,
			Quantity: int64(v.Quantity),
		})
	}

	images := make([]*productpb.CreateImageRequest, 0, len(req.Images))
	for _, img := range req.Images {
		openedFile, err := img.File.Open()
		if err != nil {
			common.JSON(c, http.StatusBadRequest, "Không mở được file", nil)
			return
		}
		defer openedFile.Close()

		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			common.JSON(c, http.StatusInternalServerError, "Đọc file thất bại", nil)
			return
		}

		base64Data := base64.StdEncoding.EncodeToString(fileBytes)

		images = append(images, &productpb.CreateImageRequest{
			ColorId:     img.ColorID,
			Base64Data:  base64Data,
			FileName:    img.File.Filename,
			IsThumbnail: *img.IsThumbnail,
			SortOrder:   int32(img.SortOrder),
		})
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
		TagIds:      req.TagIDs,
		Variants:    variants,
		Images:      images,
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

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
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

	productID := c.Param("id")

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		common.JSON(c, http.StatusBadRequest, "Không thể parse form", nil)
		return
	}

	var req request.UpdateProductForm
	form := c.Request.MultipartForm

	if title := strings.TrimSpace(c.PostForm("title")); title != "" {
		req.Title = &title
	}

	if description := strings.TrimSpace(c.PostForm("description")); description != "" {
		req.Description = &description
	}

	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 32); err == nil {
			priceFloat := float32(price)
			req.Price = &priceFloat
		}
	}

	if isSaleStr := c.PostForm("is_sale"); isSaleStr != "" {
		if isSale, err := strconv.ParseBool(isSaleStr); err == nil {
			req.IsSale = &isSale
		}
	}

	if salePriceStr := c.PostForm("sale_price"); salePriceStr != "" {
		if salePrice, err := strconv.ParseFloat(salePriceStr, 32); err == nil {
			salePriceFloat := float32(salePrice)
			req.SalePrice = &salePriceFloat
		}
	}

	if startSaleStr := c.PostForm("start_sale"); startSaleStr != "" {
		if startSale, err := time.Parse("2006-01-02", startSaleStr); err == nil {
			req.StartSale = &startSale
		}
	}

	if endSaleStr := c.PostForm("end_sale"); endSaleStr != "" {
		if endSale, err := time.Parse("2006-01-02", endSaleStr); err == nil {
			req.EndSale = &endSale
		}
	}

	req.CategoryIDs = form.Value["category_ids"]
	req.TagIDs = form.Value["tag_ids"]
	req.DeleteImageIDs = form.Value["delete_image_ids"]
	req.DeleteVariantIDs = form.Value["delete_variant_ids"]

	req.UpdateImages = []request.UpdateImageForm{}
	i := 0
	for {
		idKey := fmt.Sprintf("update_images[%d][id]", i)
		isThumbnailKey := fmt.Sprintf("update_images[%d][is_thumbnail]", i)
		sortOrderKey := fmt.Sprintf("update_images[%d][sort_order]", i)

		id := strings.TrimSpace(c.PostForm(idKey))
		if id == "" {
			break
		}

		updateImage := request.UpdateImageForm{
			ID: id,
		}

		if isThumbnailStr := c.PostForm(isThumbnailKey); isThumbnailStr != "" {
			if isThumbnail, err := strconv.ParseBool(isThumbnailStr); err == nil {
				updateImage.IsThumbnail = &isThumbnail
			}
		}

		if sortOrderStr := c.PostForm(sortOrderKey); sortOrderStr != "" {
			if sortOrder, err := strconv.Atoi(sortOrderStr); err == nil {
				updateImage.SortOrder = &sortOrder
			}
		}

		req.UpdateImages = append(req.UpdateImages, updateImage)
		i++
	}

	req.NewImages = []request.CreateImageForm{}
	j := 0
	for {
		colorIDKey := fmt.Sprintf("new_images[%d][color_id]", j)
		isThumbnailKey := fmt.Sprintf("new_images[%d][is_thumbnail]", j)
		sortOrderKey := fmt.Sprintf("new_images[%d][sort_order]", j)
		fileKey := fmt.Sprintf("new_images[%d][file]", j)

		colorID := strings.TrimSpace(c.PostForm(colorIDKey))
		if colorID == "" {
			break
		}

		isThumbnailStr := strings.TrimSpace(c.PostForm(isThumbnailKey))
		sortOrderStr := strings.TrimSpace(c.PostForm(sortOrderKey))

		isThumbnail := false
		if isThumbnailStr != "" {
			isThumbnail, _ = strconv.ParseBool(isThumbnailStr)
		}

		sortOrder := 0
		if sortOrderStr != "" {
			sortOrder, _ = strconv.Atoi(sortOrderStr)
		}

		file, err := c.FormFile(fileKey)
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Không tìm thấy file cho new image %d: %s", j, err.Error()), nil)
			return
		}

		newImage := request.CreateImageForm{
			ColorID:     colorID,
			IsThumbnail: &isThumbnail,
			SortOrder:   sortOrder,
			File:        file,
		}

		req.NewImages = append(req.NewImages, newImage)
		j++
	}

	req.UpdateVariants = []request.UpdateVariantForm{}
	k := 0
	for {
		idKey := fmt.Sprintf("update_variants[%d][id]", k)
		skuKey := fmt.Sprintf("update_variants[%d][sku]", k)
		colorKey := fmt.Sprintf("update_variants[%d][color_id]", k)
		sizeKey := fmt.Sprintf("update_variants[%d][size_id]", k)
		quantityKey := fmt.Sprintf("update_variants[%d][quantity]", k)

		id := strings.TrimSpace(c.PostForm(idKey))
		if id == "" {
			break
		}

		updateVariant := request.UpdateVariantForm{
			ID: id,
		}

		if sku := strings.TrimSpace(c.PostForm(skuKey)); sku != "" {
			updateVariant.SKU = sku
		}

		if colorID := strings.TrimSpace(c.PostForm(colorKey)); colorID != "" {
			updateVariant.ColorID = colorID
		}

		if sizeID := strings.TrimSpace(c.PostForm(sizeKey)); sizeID != "" {
			updateVariant.SizeID = sizeID
		}

		if quantityStr := strings.TrimSpace(c.PostForm(quantityKey)); quantityStr != "" {
			if quantity, err := strconv.Atoi(quantityStr); err == nil {
				updateVariant.Quantity = &quantity
			}
		}

		req.UpdateVariants = append(req.UpdateVariants, updateVariant)
		k++
	}

	req.NewVariants = []request.CreateVariantForm{}
	l := 0
	for {
		skuKey := fmt.Sprintf("new_variants[%d][sku]", l)
		colorKey := fmt.Sprintf("new_variants[%d][color_id]", l)
		sizeKey := fmt.Sprintf("new_variants[%d][size_id]", l)
		quantityKey := fmt.Sprintf("new_variants[%d][quantity]", l)

		sku := strings.TrimSpace(c.PostForm(skuKey))
		if sku == "" {
			break
		}

		colorID := strings.TrimSpace(c.PostForm(colorKey))
		sizeID := strings.TrimSpace(c.PostForm(sizeKey))
		quantityStr := strings.TrimSpace(c.PostForm(quantityKey))

		quantity := 0
		if quantityStr != "" {
			quantity, _ = strconv.Atoi(quantityStr)
		}

		newVariant := request.CreateVariantForm{
			SKU:      sku,
			ColorID:  colorID,
			SizeID:   sizeID,
			Quantity: quantity,
		}

		req.NewVariants = append(req.NewVariants, newVariant)
		l++
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var isSale *bool
	var salePrice *float32
	var startSale, endSale *string

	if req.IsSale != nil {
		isSale = req.IsSale
		if *req.IsSale {
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
	} else {
		if req.SalePrice != nil {
			salePrice = req.SalePrice
		}
		if req.StartSale != nil {
			formattedStartSale := req.StartSale.Format("2006-01-02")
			startSale = &formattedStartSale
		}
		if req.EndSale != nil {
			formattedEndSale := req.EndSale.Format("2006-01-02")
			endSale = &formattedEndSale
		}
	}

	updateVariants := make([]*productpb.UpdateVariantRequest, 0, len(req.UpdateVariants))
	for _, v := range req.UpdateVariants {
		updateVar := &productpb.UpdateVariantRequest{
			Id: v.ID,
		}
		if v.SKU != "" {
			updateVar.Sku = &v.SKU
		}
		if v.ColorID != "" {
			updateVar.ColorId = &v.ColorID
		}
		if v.SizeID != "" {
			updateVar.SizeId = &v.SizeID
		}
		if v.Quantity != nil {
			quantity := int64(*v.Quantity)
			updateVar.Quantity = &quantity
		}
		updateVariants = append(updateVariants, updateVar)
	}

	newVariants := make([]*productpb.CreateVariantRequest, 0, len(req.NewVariants))
	for _, v := range req.NewVariants {
		newVariants = append(newVariants, &productpb.CreateVariantRequest{
			Sku:      v.SKU,
			ColorId:  v.ColorID,
			SizeId:   v.SizeID,
			Quantity: int64(v.Quantity),
		})
	}

	updateImages := make([]*productpb.UpdateImageRequest, 0, len(req.UpdateImages))
	for _, img := range req.UpdateImages {
		updateImg := &productpb.UpdateImageRequest{
			Id: img.ID,
		}
		if img.IsThumbnail != nil {
			updateImg.IsThumbnail = img.IsThumbnail
		}
		if img.SortOrder != nil {
			sortOrder := int32(*img.SortOrder)
			updateImg.SortOrder = &sortOrder
		}
		updateImages = append(updateImages, updateImg)
	}

	newImages := make([]*productpb.CreateImageRequest, 0, len(req.NewImages))
	for _, img := range req.NewImages {
		openedFile, err := img.File.Open()
		if err != nil {
			common.JSON(c, http.StatusBadRequest, "Không mở được file", nil)
			return
		}
		defer openedFile.Close()

		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			common.JSON(c, http.StatusInternalServerError, "Đọc file thất bại", nil)
			return
		}

		base64Data := base64.StdEncoding.EncodeToString(fileBytes)

		newImages = append(newImages, &productpb.CreateImageRequest{
			ColorId:     img.ColorID,
			Base64Data:  base64Data,
			FileName:    img.File.Filename,
			IsThumbnail: *img.IsThumbnail,
			SortOrder:   int32(img.SortOrder),
		})
	}

	res, err := h.productClient.UpdateProduct(ctx, &productpb.UpdateProductRequest{
		Id:               productID,
		Title:            req.Title,
		Description:      req.Description,
		Price:            req.Price,
		IsSale:           isSale,
		SalePrice:        salePrice,
		StartSale:        startSale,
		EndSale:          endSale,
		CategoryIds:      req.CategoryIDs,
		TagIds:           req.TagIDs,
		DeleteImageIds:   req.DeleteImageIDs,
		DeleteVariantIds: req.DeleteVariantIDs,
		UpdateVariants:   updateVariants,
		NewVariants:      newVariants,
		UpdateImages:     updateImages,
		NewImages:        newImages,
		UserId:           user.Id,
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

	common.JSON(c, http.StatusOK, "Cập nhật sản phẩm thành công", gin.H{
		"product": res,
	})
}

func (h *ProductHandler) GetProductBySlug(c *gin.Context) {
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

func (h *ProductHandler) GetAllCategoriesAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllCategoriesAdmin(ctx, &productpb.GetManyRequest{})
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

func (h *ProductHandler) GetCategoriesNoProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetCategoriesNoProduct(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách danh mục sản phẩm thành công", gin.H{
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

func (h *ProductHandler) GetAllColorsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllColorsAdmin(ctx, &productpb.GetManyRequest{})
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

func (h *ProductHandler) GetAllColors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllColors(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả màu sắc sản phẩm thành công", gin.H{
		"colors": res.Colors,
	})
}

func (h *ProductHandler) GetAllSizesAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllSizesAdmin(ctx, &productpb.GetManyRequest{})
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

func (h *ProductHandler) GetAllSizes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllSizes(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả size sản phẩm thành công", gin.H{
		"sizes": res.Sizes,
	})
}

func (h *ProductHandler) GetAllTagsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllTagsAdmin(ctx, &productpb.GetManyRequest{})
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
		"tags": res.Tags,
	})
}

func (h *ProductHandler) GetAllTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllTags(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả tag sản phẩm thành công", gin.H{
		"tags": res.Tags,
	})
}

func (h *ProductHandler) UpdateTag(c *gin.Context) {
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

	tagID := c.Param("id")

	var req request.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.productClient.UpdateTag(ctx, &productpb.UpdateTagRequest{
		Id:     tagID,
		Name:   req.Name,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Cập nhật tag sản phẩm thành công", nil)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productID := c.Param("id")

	res, err := h.productClient.GetProductById(ctx, &productpb.GetProductByIdRequest{
		Id: productID,
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

	common.JSON(c, http.StatusOK, "Lấy chi tiết sản phẩm thành công", gin.H{
		"product": res,
	})
}

func (h *ProductHandler) GetAllProductsAdmin(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetAllProductsAdmin(ctx, &productpb.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả sản phẩm thành công", gin.H{
		"products": res.Products,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
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

	productID := c.Param("id")

	if _, err := h.productClient.DeleteProduct(ctx, &productpb.DeleteOneRequest{
		Id:     productID,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển sản phẩm vào thùng rác thành công", nil)
}

func (h *ProductHandler) DeleteProducts(c *gin.Context) {
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

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.productClient.DeleteProducts(ctx, &productpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển danh sách sản phẩm vào thùng rác thành công", nil)
}

func (h *ProductHandler) PermanentlyDeleteCategory(c *gin.Context) {
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

	if _, err := h.productClient.PermanentlyDeleteCategory(ctx, &productpb.DeleteOneRequest{
		Id:     categoryID,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Xóa danh mục sản phẩm thành công", nil)
}

func (h *ProductHandler) PermanentlyDeleteCategories(c *gin.Context) {
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

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.productClient.PermanentlyDeleteCategories(ctx, &productpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Xóa danh sách danh mục sản phẩm thành công", nil)
}

func (h *ProductHandler) UpdateColor(c *gin.Context) {
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

	var req request.UpdateColorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	colorID := c.Param("id")

	if _, err := h.productClient.UpdateColor(ctx, &productpb.UpdateColorRequest{
		Id:     colorID,
		Name:   req.Name,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Cập nhật màu sắc sản phẩm thành công", nil)
}

func (h *ProductHandler) UpdateSize(c *gin.Context) {
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

	var req request.UpdateSizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	sizeID := c.Param("id")

	if _, err := h.productClient.UpdateSize(ctx, &productpb.UpdateSizeRequest{
		Id:     sizeID,
		Name:   req.Name,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Cập nhật kích cỡ sản phẩm thành công", nil)
}

func (h *ProductHandler) DeleteColor(c *gin.Context) {
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

	colorID := c.Param("id")

	if _, err := h.productClient.DeleteColor(ctx, &productpb.DeleteOneRequest{
		Id:     colorID,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển màu sắc vào thùng rác thành công", nil)
}

func (h *ProductHandler) DeleteSize(c *gin.Context) {
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

	sizeID := c.Param("id")

	if _, err := h.productClient.DeleteSize(ctx, &productpb.DeleteOneRequest{
		Id:     sizeID,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển kích cỡ vào thùng rác thành công", nil)
}

func (h *ProductHandler) DeleteColors(c *gin.Context) {
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

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.productClient.DeleteColors(ctx, &productpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển danh sách màu sắc vào thùng rác thành công", nil)
}

func (h *ProductHandler) DeleteSizes(c *gin.Context) {
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

	var req request.DeleteManyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	if _, err := h.productClient.DeleteSizes(ctx, &productpb.DeleteManyRequest{
		Ids:    req.IDs,
		UserId: user.Id,
	}); err != nil {
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

	common.JSON(c, http.StatusOK, "Chuyển danh sách màu sắc vào thùng rác thành công", nil)
}

func (h *ProductHandler) GetDeletedProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetDeletedProducts(ctx, &protobuf.GetManyRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy tất cả sản phẩm đã xóa thành công", gin.H{
		"products": res.Products,
	})
}

func (h *ProductHandler) GetDeletedProductByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productID := c.Param("id")

	res, err := h.productClient.GetDeletedProductById(ctx, &productpb.GetProductByIdRequest{
		Id: productID,
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

	common.JSON(c, http.StatusOK, "Lấy chi tiết sản phẩm thành công", gin.H{
		"product": res,
	})
}

func (h *ProductHandler) GetDeletedColors(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetDeletedColors(ctx, &protobuf.GetManyRequest{})
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

	common.JSON(c, http.StatusOK, "Lấy tất cả màu sắc đã xóa thành công", gin.H{
		"colors": res.Colors,
	})
}

func (h *ProductHandler) GetDeletedSizes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	res, err := h.productClient.GetDeletedSizes(ctx, &protobuf.GetManyRequest{})
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

	common.JSON(c, http.StatusOK, "Lấy tất cả kích cỡ đã xóa thành công", gin.H{
		"sizes": res.Sizes,
	})
}
