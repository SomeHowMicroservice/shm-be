package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type AuthContainer struct {
	Handler *handler.AuthHandler
}

func NewAuthContainer(authClient authpb.AuthServiceClient, userClient userpb.UserServiceClient, cfg *config.AppConfig) *AuthContainer {
	handler := handler.NewAuthHandler(authClient, userClient, cfg)
	return &AuthContainer{
		Handler: handler,
	}
}
