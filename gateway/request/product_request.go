package request

import (
	"mime/multipart"
	"time"
)

type CreateCategoryRequest struct {
	Name      string   `json:"name" binding:"required,max=50"`
	Slug      *string  `json:"slug" binding:"omitempty"`
	ParentIDs []string `json:"parent_ids" binding:"omitempty,dive,uuid4"`
}

type UpdateCategoryRequest struct {
	Name      string   `json:"name" binding:"required,max=50"`
	Slug      string  `json:"slug" binding:"omitempty"`
	ParentIDs []string `json:"parent_ids" binding:"omitempty,dive,uuid4"`
}

type CreateProductRequest struct {
	Title       string     `json:"title" binding:"required,min=2"`
	Description string     `json:"description" binding:"required"`
	Price       float32    `json:"price"`
	IsSale      *bool      `json:"is_sale" binding:"required"`
	SalePrice   *float32   `json:"sale_price" binding:"omitempty"`
	StartSale   *time.Time `json:"start_sale" binding:"omitempty"`
	EndSale     *time.Time `json:"end_sale" binding:"omitempty"`
	CategoryIDs []string   `json:"category_ids" binding:"required,dive,uuid4"`
}

type CreateColorRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type CreateSizeRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

type CreateVariantRequest struct {
	SKU       string `json:"sku" binding:"required,max=50"`
	ProductID string `json:"product_id" binding:"required,uuid4"`
	ColorID   string `json:"color_id" binding:"required,uuid4"`
	SizeID    string `json:"size_id" binding:"required,uuid4"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateImageRequest struct {
	ProductID   string                `form:"product_id" binding:"required,uuid4"`
	ColorID     string                `form:"color_id" binding:"required,uuid4"`
	IsThumbnail *bool                 `form:"is_thumbnail" binding:"required"`
	SortOrder   int                   `form:"sort_order" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}
