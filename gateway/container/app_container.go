package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type Container struct {
	Auth *AuthContainer
}

func NewContainer(authClient authpb.AuthServiceClient, cfg *config.AppConfig) *Container {
	return &Container{
		Auth: NewAuthContainer(authClient, cfg),
	}
}
