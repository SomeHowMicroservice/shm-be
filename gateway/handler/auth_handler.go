package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/request"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	client authpb.AuthServiceClient
	cfg    *config.AppConfig
}

func NewAuthHandler(client authpb.AuthServiceClient, cfg *config.AppConfig) *AuthHandler {
	return &AuthHandler{
		client,
		cfg,
	}
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
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.JSON(c, http.StatusOK, "Đăng ký thành công, vui lòng kiểm tra email để lấy mã xác thực", gin.H{
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
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)
	common.JSON(c, http.StatusOK, "Đăng ký thành công", gin.H{
		"user": res.User,
	})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		message := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, message, nil)
		return
	}
	res, err := h.client.SignIn(ctx, &authpb.SignInRequest{
		Username: req.Username,
		Password: req.Password,
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
				common.JSON(c, http.StatusInternalServerError, st.Message(), nil)
			}
			return
		}
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	c.SetCookie(h.cfg.Jwt.AccessName, res.AccessToken, int(res.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, res.RefreshToken, int(res.RefreshExpiresIn), "/api/v1/auth/refresh", "", false, true)
	common.JSON(c, http.StatusOK, "Đăng nhập thành công", gin.H{
		"user": res.User,
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	c.SetCookie(h.cfg.Jwt.AccessName, "", 0, "/", "", false, true)
	c.SetCookie(h.cfg.Jwt.RefreshName, "", 0, "/api/v1/auth/refresh", "", false, true)
	common.JSON(c, http.StatusOK, "Đăng xuất thành công", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
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
	common.JSON(c, http.StatusOK, "Lấy thông tin người dùng thành công", gin.H{
		"user": toAuthResponse(user),
	})
}

func toAuthResponse(userRes *userpb.UserPublicResponse) *authpb.AuthResponse {
	return &authpb.AuthResponse{
		Id: userRes.Id,
		Username: userRes.Username,
		Email: userRes.Email,
		CreatedAt: userRes.CreatedAt,
	}
}
