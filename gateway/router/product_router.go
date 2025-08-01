package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	"github.com/SomeHowMicroservice/shm-be/services/user/model"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
)

func ProductRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, productHandler *handler.ProductHandler) {
	accessName := cfg.Jwt.AccessName
	secretKey := cfg.Jwt.SecretKey

	product := rg.Group("/products")
	{
		product.GET("/:slug", productHandler.ProductDetails)
	}

	category := rg.Group("/categories")
	{
		category.GET("/tree", productHandler.GetCategoryTree)
		category.GET("/:slug/products", productHandler.GetProductsByCategory)
	}

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{model.RoleAdmin}))
	{
		admin.POST("/categories", productHandler.CreateCategory)
		admin.GET("/categories", productHandler.GetAllCategories)
		admin.GET("/categories/:id", productHandler.CategoryAdminDetails)
		admin.POST("/products", productHandler.CreateProduct)
		admin.POST("/products/colors", productHandler.CreateColor)
		admin.POST("/products/sizes", productHandler.CreateSize)
		admin.POST("/products/variants", productHandler.CreateVariant)
		admin.POST("/products/images", productHandler.CreateImage)
		admin.POST("/products/tags", productHandler.CreateTag)
	}
}
