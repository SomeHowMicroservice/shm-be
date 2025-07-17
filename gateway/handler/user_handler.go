package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
	cfg        *config.AppConfig
}

func NewUserHandler(userClient userpb.UserServiceClient, cfg *config.AppConfig) *UserHandler {
	return &UserHandler{
		userClient,
		cfg,
	}
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}

	var firstName, lastName, gender, dob string
	if req.FirstName != nil {
		firstName = *req.FirstName
	}
	if req.LastName != nil {
		lastName = *req.LastName
	}
	if req.Gender != nil {
		gender = *req.Gender
	}
	if req.DOB != nil {
		dob = req.DOB.Format("2006-01-02")
	}

	res, err := h.userClient.UpdateUserProfile(ctx, &userpb.UpdateUserProfileRequest{
		UserId:    user.Id,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    gender,
		Dob:       dob,
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

	common.JSON(c, http.StatusOK, "Cập nhật hồ sơ thành công", gin.H{
		"user": ToAuthResponse(res),
	})
}
