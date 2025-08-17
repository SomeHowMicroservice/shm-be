package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/services/post/config"
	"github.com/SomeHowMicroservice/shm-be/services/post/container"
	"github.com/SomeHowMicroservice/shm-be/services/post/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/post/protobuf"
	"google.golang.org/grpc"
)

var (
	userAddr = "localhost:8082"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Post Service thất bại: %v", err)
	}

	db, err := initialization.InitDB(cfg)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB ở Post Service: %v", err)
	}

	userAddr = cfg.App.ServerHost + fmt.Sprintf(":%d", cfg.Services.UserPort)
	clients := initialization.InitClients(userAddr)

	grpcServer := grpc.NewServer()
	postContainer := container.NewContainer(cfg, db, grpcServer, clients.UserClient)
	protobuf.RegisterPostServiceServer(grpcServer, postContainer.GRPCHandler)

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
