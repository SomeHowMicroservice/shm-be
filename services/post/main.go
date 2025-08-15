package main

import (
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/services/post/config"
	"github.com/SomeHowMicroservice/shm-be/services/post/initialization"
	"google.golang.org/grpc"
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

	log.Println(db)

	grpcServer := grpc.NewServer()
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
