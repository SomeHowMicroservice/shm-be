package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/gin-gonic/gin"
)

func AuthRouter(rg *gin.RouterGroup, authHandler handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-up/verify", authHandler.VerifySignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}
}