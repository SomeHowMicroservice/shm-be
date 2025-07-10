package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/gin-gonic/gin"
)

func AuthRouter(rg *gin.RouterGroup, authClient authpb.AuthServiceClient,  authHandler handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)
		auth.POST("/sign-up/verify", authHandler.VerifySignUp)
		auth.POST("/sign-in", authHandler.SignIn)
	}
}