package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
)

func UserRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, userHandler *handler.UserHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	user := rg.Group("/users")
	{
		user.PATCH("/profiles/:id", middleware.RequireAuth(accessName, secretKey, userClient), userHandler.UpdateProfile)

		user.GET("/me/measurements", middleware.RequireAuth(accessName, secretKey, userClient), userHandler.MyMeasurements)

		user.PATCH("/measurements/:id", middleware.RequireAuth(accessName, secretKey, userClient), userHandler.UpdateMeasurement)
	}
}
