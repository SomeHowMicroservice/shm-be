package container

import (
	"github.com/SomeHowMicroservice/shm-be/services/product/handler"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	colorRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/color"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	sizeRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/size"
	variantRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/variant"
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
	colorRepo := colorRepo.NewColorRepository(db)
	sizeRepo := sizeRepo.NewSizeRepository(db)
	variantRepo := variantRepo.NewVariantRepository(db)
	svc := service.NewProductService(categoryRepo, productRepo, colorRepo, sizeRepo, variantRepo)
	hdl := handler.NewGRPCHandler(grpcServcer, svc)
	return &Container{hdl}
}
