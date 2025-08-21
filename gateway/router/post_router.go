package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/gateway/protobuf/user"
	"github.com/gin-gonic/gin"
)

func PostRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, postHandler *handler.PostHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{"contributor"}))
	{
		admin.POST("/topics", postHandler.CreateTopic)
		admin.GET("/topics", postHandler.GetAllTopicsAdmin)
		admin.PUT("/topics/:id", postHandler.UpdateTopic)
	}
}
