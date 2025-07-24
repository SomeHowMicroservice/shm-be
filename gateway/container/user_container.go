package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type UserContainer struct {
	Handler *handler.UserHandler
}

func NewUserContainer(userClient userpb.UserServiceClient) *UserContainer {
	handler := handler.NewUserHandler(userClient)
	return &UserContainer{handler}
}