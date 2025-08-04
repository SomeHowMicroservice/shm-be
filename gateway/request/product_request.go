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
	Slug      string   `json:"slug" binding:"omitempty"`
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

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

type CreateVariantRequest struct {
	SKU       string `json:"sku" binding:"required,max=50"`
	ProductID string `json:"product_id" binding:"required,uuid4"`
	ColorID   string `json:"color_id" binding:"required,uuid4"`
	SizeID    string `json:"size_id" binding:"required,uuid4"`
	Quantity  int    `json:"quantity" binding:"required"`
}

type CreateImageForm struct {
	ColorID     string                `form:"color_id" binding:"required,uuid4"`
	IsThumbnail bool                 `form:"is_thumbnail" binding:"required"`
	SortOrder   int                   `form:"sort_order" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

type CreateImageRequest struct {
	ProductID   string                `form:"product_id" binding:"required,uuid4"`
	ColorID     string                `form:"color_id" binding:"required,uuid4"`
	IsThumbnail *bool                 `form:"is_thumbnail" binding:"required"`
	SortOrder   int                   `form:"sort_order" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

type CreateProductForm struct {
	Title       string              `form:"title" binding:"required,min=2"`
	Description string              `form:"description" binding:"required"`
	Price       float32             `form:"price"`
	IsSale      *bool               `form:"is_sale" binding:"required"`
	SalePrice   *float32            `form:"sale_price" binding:"omitempty"`
	StartSale   *time.Time          `form:"start_sale" binding:"omitempty"`
	EndSale     *time.Time          `form:"end_sale" binding:"omitempty"`
	CategoryIDs []string            `form:"category_ids" binding:"required,dive,uuid4"`
	TagIDs      []string            `form:"tag_ids" validate:"required,dive,uuid4"`
	Variants    []CreateVariantForm `form:"variants" validate:"required,dive"`
	Images      []CreateImageForm   `form:"images" validate:"required,dive"`
}

type CreateVariantForm struct {
	SKU       string `form:"sku" binding:"required,max=50"`
	ColorID   string `form:"color_id" binding:"required,uuid4"`
	SizeID    string `form:"size_id" binding:"required,uuid4"`
	Quantity  int    `form:"quantity" binding:"required"`
}
