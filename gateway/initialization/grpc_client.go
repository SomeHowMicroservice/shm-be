package initialization

import (
	"context"
	"log"
	"time"

	authpb "github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	AuthClient authpb.AuthServiceClient
}

func InitClients(authAddr string) *GRPCClients {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	authConn, err := grpc.DialContext(ctx, authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới Auth Service thất bại: %v", err)
	}
	authClient := authpb.NewAuthServiceClient(authConn)

	return &GRPCClients{
		AuthClient: authClient,
	}
}
