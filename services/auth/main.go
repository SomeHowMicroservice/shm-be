package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/SomeHowMicroservice/shm-be/common/smtp"
	"github.com/SomeHowMicroservice/shm-be/services/auth/config"
	"github.com/SomeHowMicroservice/shm-be/services/auth/handler"
	"github.com/SomeHowMicroservice/shm-be/services/auth/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/auth/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/auth/repository"
	"github.com/SomeHowMicroservice/shm-be/services/auth/service"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userAddr = "localhost:8082"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Auth Service thất bại: %v", err)
	}

	rdb, err := initialization.InitCache(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userAddr = cfg.App.ServerHost + fmt.Sprintf(":%d", cfg.Services.UserPort)
	userConn, err := grpc.DialContext(ctx, userAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Kết nối tới UserService thất bại: %v", err)
	}
	userClient := userpb.NewUserServiceClient(userConn)

	grpcServer := grpc.NewServer()
	cacheRepo := repository.NewCacheRepository(rdb)

	mailerCfg := &smtp.MailerConfig{
		Host: cfg.SMTP.Host,
		Port: cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
	}

	mailer := smtp.NewMailer(mailerCfg)
	svc := service.NewAuthService(cacheRepo, userClient, mailer, cfg)
	authHandler := handler.NewGRPCHandler(grpcServer, svc)

	protobuf.RegisterAuthServiceServer(grpcServer, authHandler)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.App.GRPCPort))
	if err != nil {
		log.Fatalf("Không thể lắng nghe: %v", err)
	}
	defer lis.Close()

	log.Println("Khởi chạy service thành công")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Kết nối tới phục vụ thất bại: %v", err)
	}
}
