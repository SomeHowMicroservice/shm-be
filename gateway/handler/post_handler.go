package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	postpb "github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostHandler struct {
	postClient postpb.PostServiceClient
}

func NewPostHandler(postClient postpb.PostServiceClient) *PostHandler {
	return &PostHandler{postClient}
}

func (h *PostHandler) CreateTopic(c *gin.Context) {
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

	var req request.CreateTopicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var slug *string
	if req.Slug != nil {
		slug = req.Slug
	}

	res, err := h.postClient.CreateTopic(ctx, &postpb.CreateTopicRequest{
		Name:   req.Name,
		Slug:   slug,
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

	common.JSON(c, http.StatusCreated, "Tạo chủ đề bài viết thành công", gin.H{
		"topic_id": res.Id,
	})
}
