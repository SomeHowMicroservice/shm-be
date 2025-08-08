package container

import (
	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"github.com/SomeHowMicroservice/shm-be/services/product/handler"
	categoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/category"
	colorRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/color"
	imageRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/image"
	inventoryRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/inventory"
	productRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/product"
	sizeRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/size"
	tagRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/tag"
	variantRepo "github.com/SomeHowMicroservice/shm-be/services/product/repository/variant"
	"github.com/SomeHowMicroservice/shm-be/services/product/service"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
	ImageRepo imageRepo.ImageRepository
}

func NewContainer(cfg *config.Config, db *gorm.DB, mqChannel *amqp091.Channel, grpcServer *grpc.Server, userClient userpb.UserServiceClient) *Container {
	categoryRepo := categoryRepo.NewCategoryRepository(db)
	productRepo := productRepo.NewProductRepository(db)
	tagRepo := tagRepo.NewTagRepository(db)
	colorRepo := colorRepo.NewColorRepository(db)
	sizeRepo := sizeRepo.NewSizeRepository(db)
	variantRepo := variantRepo.NewVariantRepository(db)
	inventoryRepo := inventoryRepo.NewInventoryRepository(db)
	imageRepo := imageRepo.NewImageRepository(db)
	svc := service.NewProductService(cfg, userClient, mqChannel, categoryRepo, productRepo, tagRepo, colorRepo, sizeRepo, variantRepo, inventoryRepo, imageRepo)
	hdl := handler.NewGRPCHandler(grpcServer, svc)
	return &Container{hdl, imageRepo}
}
