package request

type CreateCategoryRequest struct {
	Name      string   `json:"name" binding:"required,min=1,max=50"`
	Slug      *string  `json:"slug" binding:"omitempty,min=1"`
	ParentIDs []string `json:"parent_ids"`
}
