package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type Container struct {
	Auth *AuthContainer
	User *UserContainer
}

func NewContainer(authClient authpb.AuthServiceClient, userClient userpb.UserServiceClient, cfg *config.AppConfig) *Container {
	auth := NewAuthContainer(authClient, cfg)
	user := NewUserContainer(userClient, cfg)
	return &Container{
		auth,
		user,
	}
}
