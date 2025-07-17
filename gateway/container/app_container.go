package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type Container struct {
	Auth *AuthContainer
}

func NewContainer(authClient authpb.AuthServiceClient, userClient userpb.UserServiceClient, cfg *config.AppConfig) *Container {
	return &Container{
		Auth: NewAuthContainer(authClient, userClient, cfg),
	}
}
