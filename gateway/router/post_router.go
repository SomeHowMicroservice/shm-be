package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
)

func PostRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, postHandler *handler.PostHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{model.RoleContributor}))
	{
		admin.POST("/topics", postHandler.CreateTopic)
	}
}
