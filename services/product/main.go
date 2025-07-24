package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"github.com/SomeHowMicroservice/shm-be/services/product/handler"
	"github.com/SomeHowMicroservice/shm-be/services/product/initialization"
	"github.com/SomeHowMicroservice/shm-be/services/product/protobuf"
	repository "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	"github.com/SomeHowMicroservice/shm-be/services/product/service"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình User Service thất bại: %v", err)
	}

	client, err := initialization.InitDB(cfg)
	if err != nil {
		log.Fatalf("Lỗi kết nối DB ở Product Services: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Lỗi ngắt kết nối DB ở Product Services: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()
	categoryRepo := repository.NewCategoryRepository(client.Database("product_db"))
	svc := service.NewProductService(categoryRepo)
	productHandler := handler.NewGRPCHandler(grpcServer, svc)
	protobuf.RegisterProductServiceServer(grpcServer, productHandler)
	

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
