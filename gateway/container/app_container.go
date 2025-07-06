package container

import authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"

type Container struct {
	Auth *AuthContainer
}

func NewContainer(authClient authpb.AuthServiceClient) *Container {
	return &Container{
		Auth: NewAuthContainer(authClient),
	}
}
