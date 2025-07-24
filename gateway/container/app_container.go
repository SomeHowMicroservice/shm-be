package container

import (
	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
)

type Container struct {
	Auth *AuthContainer
	User *UserContainer
	Product *ProductContainer
}

func NewContainer(authClient authpb.AuthServiceClient, userClient userpb.UserServiceClient, productClient productpb.ProductServiceClient, cfg *config.AppConfig) *Container {
	auth := NewAuthContainer(authClient, cfg)
	user := NewUserContainer(userClient)
	product := NewProductHandler(productClient)
	return &Container{
		auth,
		user,
		product,
	}
}
