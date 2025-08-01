package container

import (
	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"github.com/SomeHowMicroservice/shm-be/services/product/handler"
	"github.com/SomeHowMicroservice/shm-be/services/product/imagekit"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	colorRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/color"
	imageRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/image"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	sizeRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/size"
	tagRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/tag"
	variantRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/variant"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/SomeHowMicroservice/shm-be/services/product/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB, grpcServer *grpc.Server, userClient userpb.UserServiceClient) *Container {
	categoryRepo := categoryRepo.NewCategoryRepository(db)
	productRepo := productRepo.NewProductRepository(db)
	tagRepo := tagRepo.NewTagRepository(db)
	colorRepo := colorRepo.NewColorRepository(db)
	sizeRepo := sizeRepo.NewSizeRepository(db)
	variantRepo := variantRepo.NewVariantRepository(db)
	imageRepo := imageRepo.NewImageRepository(db)
	imageKitSvc := imagekit.NewImageKitService(cfg)
	svc := service.NewProductService(userClient, imageKitSvc, categoryRepo, productRepo, tagRepo, colorRepo, sizeRepo, variantRepo, imageRepo)
	hdl := handler.NewGRPCHandler(grpcServer, svc)
	return &Container{hdl}
}
