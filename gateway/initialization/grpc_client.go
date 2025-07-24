package initialization

import (
	"context"
	"log"
	"time"

	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	productpb "github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient    authpb.AuthServiceClient
	UserClient    userpb.UserServiceClient
	ProductClient productpb.ProductServiceClient
}

func InitClients(authAddr string, userAddr string, productAddr string) *GRPCClients {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Kết nối gRPC tới Auth Service
	authConn, err := grpc.DialContext(ctx, authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Auth Service thất bại: %v", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	// Kết nối gRPC tới User Service
	userConn, err := grpc.DialContext(ctx, userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới User Service thất bại: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	// Kết nối gRPC tới Product Service
	productConn, err := grpc.DialContext(ctx, productAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Product Service thất bại: %v", err)
	}
	productClient := productpb.NewProductServiceClient(productConn)

	return &GRPCClients{
		authClient,
		userClient,
		productClient,
	}
}
