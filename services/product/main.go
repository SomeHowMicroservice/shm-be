package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"github.com/SomeHowMicroservice/shm-be/services/product/container"
	"github.com/SomeHowMicroservice/shm-be/services/product/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
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
	productContainer := container.NewContainer(db, grpcServer)
	protobuf.RegisterProductServiceServer(grpcServer, productContainer.GRPCHandler)
	

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
