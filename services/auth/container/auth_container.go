package container

import (
	"github.com/SomeHowMicroservice/shm-be/common/smtp"
	"github.com/SomeHowMicroservice/shm-be/services/auth/config"
	"github.com/SomeHowMicroservice/shm-be/services/auth/handler"
	"github.com/SomeHowMicroservice/shm-be/services/auth/repository"
	"github.com/SomeHowMicroservice/shm-be/services/auth/service"
	userpb "github.com/SomeHowMicroservice/shm-be/services/user/protobuf"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type Container struct {
	GRPCHandler *handler.GRPCHandler
	Mailer      smtp.Mailer
}

func NewContainer(cfg *config.Config, rdb *redis.Client, mqChan *amqp091.Channel, grpcServer *grpc.Server, userClient userpb.UserServiceClient) *Container {
	mailerCfg := &smtp.MailerConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
	}
	mailer := smtp.NewMailer(mailerCfg)
	cacheRepo := repository.NewCacheRepository(rdb)
	svc := service.NewAuthService(cacheRepo, userClient, mailer, cfg, mqChan)
	hdl := handler.NewGRPCHandler(grpcServer, svc)
	return &Container{hdl,mailer}
}
