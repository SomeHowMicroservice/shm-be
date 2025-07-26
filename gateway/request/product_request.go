package request

import "time"

type CreateCategoryRequest struct {
	Name      string   `json:"name" binding:"required,min=1,max=50"`
	Slug      *string  `json:"slug" binding:"omitempty,min=1"`
	ParentIDs []string `json:"parent_ids" binding:"required,dive,uuid4"`
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
