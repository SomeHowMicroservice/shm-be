package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	client authpb.AuthServiceClient
}

func NewAuthHandler(client authpb.AuthServiceClient) *AuthHandler {
	return &AuthHandler{client}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}
	res, err := h.client.SignUp(ctx, &authpb.SignUpRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, "Lỗi hệ thống", nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, "Lỗi hệ thống", nil)
		return
	}

	common.JSON(c, http.StatusOK, res.Message, gin.H{
		"registration_token": res.RegistrationToken,
	})
}

func (h *AuthHandler) VerifySignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.VerifySignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}
	res, err := h.client.VerifySignUp(ctx, &authpb.VerifySignUpRequest{
		RegistrationToken: req.RegistrationToken,
		Otp:               req.Otp,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				common.JSON(c, http.StatusConflict, st.Message(), nil)
			case codes.NotFound:
				common.JSON(c, http.StatusNotFound, st.Message(), nil)
			case codes.InvalidArgument:
				common.JSON(c, http.StatusBadRequest, st.Message(), nil)
			default:
				common.JSON(c, http.StatusInternalServerError, "Lỗi hệ thống", nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, "Lỗi hệ thống", nil)
		return
	}
	c.SetCookie("access-token", res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie("refresh-token", res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)
	common.JSON(c, http.StatusOK, res.Message, gin.H{
		"user": res.User,
	})
}
