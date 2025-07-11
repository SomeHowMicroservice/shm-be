package middleware

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireAuth(accessName string, secretKey string, userClient userpb.UserServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "Không tìm thấy token",
			})
			return
		}
		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}
		userID, userRoles, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: err.Error(),
			})
			return
		}
		ctx := c.Request.Context()
		userRes, err := fetchUserFromUserService(ctx, userID, userClient)
		if err != nil {
			switch err {
			case customErr.ErrUserNotFound:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
					Message: err.Error(),
				})
				return
			}
		}
		if !slices.Equal(userRes.Roles, userRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse{
				Message: "người dùng không hợp lệ",
			})
			return
		}
		c.Set("user", userRes)
		c.Next()
	}
}

func fetchUserFromUserService(ctx context.Context, userID string, userClient userpb.UserServiceClient) (*userpb.UserPublicResponse, error) {
	userRes, err := userClient.GetUserPublicById(ctx, &userpb.GetUserByIdRequest{Id: userID})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, customErr.ErrUserNotFound
			default:
				return nil, fmt.Errorf("lỗi user service: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("lỗi hệ thống: %w", err)
	}
	return userRes, nil
}
