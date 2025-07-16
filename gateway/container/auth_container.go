package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthContainer struct {
	Handler *handler.AuthHandler
}

func NewAuthContainer(client authpb.AuthServiceClient, cfg *config.AppConfig) *AuthContainer {
	handler := handler.NewAuthHandler(client, cfg)
	return &AuthContainer{
		Handler: handler,
	}
}
