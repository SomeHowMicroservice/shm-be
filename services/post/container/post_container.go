package container

import (
	"github.com/SomeHowMicroservice/shm-be/services/post/config"
	"github.com/SomeHowMicroservice/shm-be/services/post/handler"
	topicRepo "github.com/SomeHowMicroservice/shm-be/services/post/repository/topic"
	"github.com/SomeHowMicroservice/shm-be/services/post/service"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
}

func NewContainer(cfg *config.Config, db *gorm.DB, grpcServer *grpc.Server, userClient userpb.UserServiceClient) *Container {
	topicRepo := topicRepo.NewTopicRepository(db)
	svc := service.NewPostService(topicRepo)
	hdl := handler.NewGRPCHandler(grpcServer, svc)
	return &Container{hdl}
}
