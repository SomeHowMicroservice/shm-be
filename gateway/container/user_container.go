package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserContainer struct {
	Handler *handler.UserHandler
}

func NewUserContainer(userClient userpb.UserServiceClient, cfg *config.AppConfig) *UserContainer {
	handler := handler.NewUserHandler(userClient, cfg)
	return &UserContainer{handler}
}