package middleware

import (
	"net/http"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	"github.com/gin-gonic/gin"
)

func RequireAuth(accessName string, secretKey string) gin.HandlerFunc {
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
		c.Set("userID", userID)
		c.Set("userRoles", userRoles)
		c.Next()
	}
}
