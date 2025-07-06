package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/handler"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
)

type AuthContainer struct {
	Handler handler.AuthHandler
}

func NewAuthContainer(client authpb.AuthServiceClient) *AuthContainer {
	handler := handler.NewAuthHandler(client)
	return &AuthContainer{
		Handler: *handler,
	}
}
