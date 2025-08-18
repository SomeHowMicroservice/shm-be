package initialization

import (
	"context"
	"log"
	"time"

	"github.com/SomeHowMicroservice/shm-be/gateway/common"
	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	postpb "github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient    authpb.AuthServiceClient
	UserClient    userpb.UserServiceClient
	ProductClient productpb.ProductServiceClient
	PostClient    postpb.PostServiceClient
}

func InitClients(ca *common.ClientAddresses) *GRPCClients {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authConn, err := grpc.DialContext(ctx, ca.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Auth Service thất bại: %v", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	userConn, err := grpc.DialContext(ctx, ca.UserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới User Service thất bại: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	productConn, err := grpc.DialContext(ctx, ca.ProductAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Product Service thất bại: %v", err)
	}
	productClient := productpb.NewProductServiceClient(productConn)

	postConn, err := grpc.DialContext(ctx, ca.PostAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Post Service thất bại: %v", err)
	}
	postClient := postpb.NewPostServiceClient(postConn)

	return &GRPCClients{
		authClient,
		userClient,
		productClient,
		postClient,
	}
}
