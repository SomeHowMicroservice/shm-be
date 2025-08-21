package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/user/config"
	"github.com/SomeHowMicroservice/shm-be/user/container"
	"github.com/SomeHowMicroservice/shm-be/user/initialization"
	userpb "github.com/SomeHowMicroservice/shm-be/user/protobuf/user"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình User Service thất bại: %v", err)
	}

	db, err := initialization.InitDB(cfg)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB ở User Service: %v", err)
	}

	grpcServer := grpc.NewServer()
	userContainer := container.NewContainer(db, grpcServer)
	userpb.RegisterUserServiceServer(grpcServer, userContainer.GRPCHandler)

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
