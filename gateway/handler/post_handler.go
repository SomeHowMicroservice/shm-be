package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	postpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/post"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
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

// func (h *PostHandler) GetAllTopicsAdmin(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
// 	defer cancel()

// 	res, err := h.postClient.GetAllTopicsAdmin(ctx, &postpb.GetManyRequest{})
// 	if err != nil {
// 		if st, ok := status.FromError(err); ok {
// 			common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
// 			return
// 		}
// 		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
// 		return
// 	}

// 	common.JSON(c, http.StatusOK, "Lấy danh sách chủ đề bài viết thành công", gin.H{
// 		"topics": res.Topics,
// 	})
// }

// func (h *PostHandler) UpdateTopic(c *gin.Context) {
// 	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
// 	defer cancel()

// 	userAny, exists := c.Get("user")
// 	if !exists {
// 		common.JSON(c, http.StatusUnauthorized, "không có thông tin người dùng", nil)
// 		return
// 	}

// 	user, ok := userAny.(*userpb.UserPublicResponse)
// 	if !ok {
// 		common.JSON(c, http.StatusUnauthorized, "không thể chuyển đổi thông tin người dùng", nil)
// 		return
// 	}

// 	var req request.UpdateTopic
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		message := common.HandleValidationError(err)
// 		common.JSON(c, http.StatusBadRequest, message, nil)
// 		return
// 	}

// 	topicID := c.Param("id")

// 	if _, err := h.postClient.UpdateTopic(ctx, &postpb.UpdateTopicRequest{
// 		Id: topicID,
// 		Name: req.Name,
// 		Slug: req.Slug,
// 		UserId: user.Id,
// 	}); err != nil {
// 		if st, ok := status.FromError(err); ok {
// 			switch st.Code() {
// 			case codes.NotFound:
// 				common.JSON(c, http.StatusNotFound, st.Message(), nil)
// 			case codes.AlreadyExists:
// 				common.JSON(c, http.StatusConflict, st.Message(), nil)
// 			default:
// 				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
// 			}
// 			return
// 		}
// 	}

// 	common.JSON(c, http.StatusOK, "Cập nhật chủ đề bài viết thành công", nil)
// }
