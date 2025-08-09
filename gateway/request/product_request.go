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

type UpdateColorRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type UpdateSizeRequest struct {
	Name string `json:"name" binding:"required,max=20"`
}

type CreateImageForm struct {
	ColorID     string                `form:"color_id" validate:"required,uuid4"`
	IsThumbnail *bool                 `form:"is_thumbnail" validate:"required"`
	SortOrder   int                   `form:"sort_order" validate:"required,gt=0"`
	File        *multipart.FileHeader `form:"file" validate:"required"`
}

type CreateProductForm struct {
	Title       string              `form:"title" validate:"required,min=2"`
	Description string              `form:"description" validate:"required,min=1"`
	Price       float32             `form:"price" validate:"required,gt=0"`
	IsSale      *bool               `form:"is_sale" validate:"required"`
	SalePrice   *float32            `form:"sale_price" validate:"omitempty,gt=0"`
	StartSale   *time.Time          `form:"start_sale" validate:"omitempty"`
	EndSale     *time.Time          `form:"end_sale" validate:"omitempty"`
	CategoryIDs []string            `form:"category_ids" validate:"required,dive,uuid4"`
	TagIDs      []string            `form:"tag_ids" validate:"required,dive,uuid4"`
	Variants    []CreateVariantForm `form:"variants" validate:"required,dive"`
	Images      []CreateImageForm   `form:"images" validate:"required,dive"`
}

type CreateVariantForm struct {
	SKU      string `form:"sku" validate:"required,max=50"`
	ColorID  string `form:"color_id" validate:"required,uuid4"`
	SizeID   string `form:"size_id" validate:"required,uuid4"`
	Quantity int    `form:"quantity" validate:"required,min=1"`
}

type UpdateProductForm struct {
	Title            *string             `form:"title" validate:"omitempty,min=2"`
	Description      *string             `form:"description" validate:"omitempty,min=1"`
	Price            *float32            `form:"price" validate:"omitempty,gt=0"`
	IsSale           *bool               `form:"is_sale" validate:"omitempty"`
	SalePrice        *float32            `form:"sale_price" validate:"omitempty,gt=0"`
	StartSale        *time.Time          `form:"start_sale" validate:"omitempty"`
	EndSale          *time.Time          `form:"end_sale" validate:"omitempty"`
	CategoryIDs      []string            `form:"category_ids" validate:"omitempty,dive,uuid4"`
	TagIDs           []string            `form:"tag_ids" validate:"omitempty,dive,uuid4"`
	DeleteImageIDs   []string            `form:"delete_image_ids" validate:"omitempty,dive,uuid4"`
	UpdateImages     []UpdateImageForm   `form:"update_images" validate:"omitempty,dive"`
	NewImages        []CreateImageForm   `form:"new_images" validate:"omitempty,dive"`
	DeleteVariantIDs []string            `form:"delete_variant_ids" validate:"omitempty,dive,uuid4"`
	UpdateVariants   []UpdateVariantForm `form:"update_variants" validate:"omitempty,dive"`
	NewVariants      []CreateVariantForm `form:"new_variants" validate:"omitempty,dive"`
}

type UpdateVariantForm struct {
	ID       string `form:"id" validate:"required,uuid4"`
	SKU      string `form:"sku" validate:"omitempty,max=50"`
	ColorID  string `form:"color_id" validate:"omitempty,uuid4"`
	SizeID   string `form:"size_id" validate:"omitempty,uuid4"`
	Quantity *int   `form:"quantity" validate:"omitempty,min=1"`
}

type UpdateImageForm struct {
	ID          string `form:"id" validate:"required,uuid4"`
	IsThumbnail *bool  `form:"is_thumbnail" validate:"omitempty"`
	SortOrder   *int   `form:"sort_order" validate:"omitempty,min=1"`
}

type DeleteManyRequest struct {
	IDs []string `json:"ids" binding:"required,dive,uuid4"`
}
