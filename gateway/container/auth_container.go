package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthContainer struct {
	Handler *handler.AuthHandler
}

func NewAuthContainer(authClient authpb.AuthServiceClient, cfg *config.AppConfig) *AuthContainer {
	handler := handler.NewAuthHandler(authClient, cfg)
	return &AuthContainer{handler}
}
