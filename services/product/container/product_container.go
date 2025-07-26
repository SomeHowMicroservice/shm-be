package container

import (
	"github.com/SomeHowMicroservice/shm-be/services/product/handler"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	"github.com/SomeHowMicroservice/shm-be/services/product/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
}

func NewContainer(db *gorm.DB, grpcServcer *grpc.Server) *Container {
	categoryRepo := categoryRepo.NewCategoryRepository(db)
	productRepo := productRepo.NewProductRepository(db)
	svc := service.NewProductService(categoryRepo, productRepo)
	hdl := handler.NewGRPCHandler(grpcServcer, svc)
	return &Container{hdl}
}