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
		category.GET("/no-child", productHandler.GetCategoriesNoChild)
	}

	color := rg.Group("/colors")
	{
		color.GET("", productHandler.GetAllColors)
	}

	size := rg.Group("/sizes")
	{
		size.GET("", productHandler.GetAllSizes)
	}

	tag := rg.Group("/tags")
	{
		tag.GET("", productHandler.GetAllTags)
	}

	admin := rg.Group("/admin", middleware.RequireAuth(accessName, secretKey, userClient), middleware.RequireMultiRoles([]string{model.RoleAdmin}))
	{
		admin.POST("/categories", productHandler.CreateCategory)
		admin.GET("/categories", productHandler.GetAllCategoriesAdmin)
		admin.GET("/categories/:id", productHandler.CategoryAdminDetails)
		admin.PUT("/categories/:id", productHandler.UpdateCategory)
		admin.POST("/products", productHandler.CreateProduct)
		admin.POST("/colors", productHandler.CreateColor)
		admin.GET("/colors", productHandler.GetAllColorsAdmin)
		admin.POST("/sizes", productHandler.CreateSize)
		admin.GET("/sizes", productHandler.GetAllSizesAdmin)
		admin.POST("/tags", productHandler.CreateTag)
		admin.GET("/tags", productHandler.GetAllTagsAdmin)
		admin.PUT("/tags/:id", productHandler.UpdateTag)
	}
}
