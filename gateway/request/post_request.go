package request

type CreateTopicRequest struct {
	Name      string   `json:"name" binding:"required,max=50"`
	Slug      *string  `json:"slug" binding:"omitempty"`
}