package router

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	"github.com/SomeHowMicroservice/shm-be/gateway/middleware"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/gin-gonic/gin"
)

func AuthRouter(rg *gin.RouterGroup, cfg *config.AppConfig, userClient userpb.UserServiceClient, authHandler *handler.AuthHandler) {
	accessName := cfg.Jwt.AccessName
	refreshName := cfg.Jwt.RefreshName
	secretKey := cfg.Jwt.SecretKey

	auth := rg.Group("/auth")
	{
		auth.POST("/sign-up", authHandler.SignUp)

		auth.POST("/sign-up/verify", authHandler.VerifySignUp)

		auth.POST("/sign-in", authHandler.SignIn)

		auth.POST("/sign-out", middleware.RequireAuth(accessName, secretKey, userClient), authHandler.SignOut)

		auth.GET("/me", middleware.RequireAuth(accessName, secretKey, userClient), authHandler.GetMe)

		auth.GET("/refresh", middleware.RequireRefreshToken(refreshName, secretKey, userClient), authHandler.RefreshToken)

		auth.POST("/change-password", middleware.RequireAuth(accessName, secretKey, userClient), authHandler.ChangePassword)

		auth.POST("/forgot-password", authHandler.ForgotPassword)

		auth.POST("/forgot-password/verify", authHandler.VerifyForgotPassword)

		auth.POST("/reset-password", authHandler.ResetPassword)
	}
}