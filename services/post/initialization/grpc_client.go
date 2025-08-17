package initialization

import (
	"context"
	"log"
	"time"

	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	UserClient userpb.UserServiceClient
}

func InitClients(userAddr string) *GRPCClients {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	userConn, err := grpc.DialContext(ctx, userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới User Service thất bại: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	return &GRPCClients{userClient}
}