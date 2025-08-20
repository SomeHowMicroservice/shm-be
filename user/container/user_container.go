package container

import (
	"github.com/SomeHowMicroservice/shm-be/user/handler"
	addressRepo "github.com/SomeHowMicroservice/shm-be/user/repository/address"
	measurementRepo "github.com/SomeHowMicroservice/shm-be/user/repository/measurement"
	profileRepo "github.com/SomeHowMicroservice/shm-be/user/repository/profile"
	roleRepo "github.com/SomeHowMicroservice/shm-be/user/repository/role"
	userRepo "github.com/SomeHowMicroservice/shm-be/user/repository/user"
	"github.com/SomeHowMicroservice/shm-be/user/service"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
}

func NewContainer(db *gorm.DB, grpcServer *grpc.Server) *Container {
	userRepo := userRepo.NewUserRepository(db)
	profileRepo := profileRepo.NewProfileRepository(db)
	measurementRepo := measurementRepo.NewMeasurementRepository(db)
	roleRepo := roleRepo.NewRoleRepository(db)
	addressRepo := addressRepo.NewAddressRepository(db)
	svc := service.NewUserService(userRepo, roleRepo, profileRepo, measurementRepo, addressRepo)
	hdl := handler.NewGRPCHandler(grpcServer, svc)
	return &Container{hdl}
}
