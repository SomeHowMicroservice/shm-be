package request

type CreateTopicRequest struct {
	Name string  `json:"name" binding:"required,max=50"`
	Slug *string `json:"slug" binding:"omitempty,max=50"`
}

type UpdateTopic struct {
	Name string  `json:"name" binding:"required,max=50"`
	Slug string `json:"slug" binding:"required,max=50"`
}
