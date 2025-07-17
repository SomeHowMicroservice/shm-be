package main

import (
	"log"
	"fmt"

	"github.com/SomeHowMicroservice/shm-be/gateway/config"
	"github.com/SomeHowMicroservice/shm-be/gateway/container"
	"github.com/SomeHowMicroservice/shm-be/gateway/initialization"
	"github.com/SomeHowMicroservice/shm-be/gateway/router"
	"github.com/gin-gonic/gin"
)

var (
	authAddr = "localhost:8081"
	userAddr = "localhost:8082"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Gateway thất bại: %v", err)
	}

	authAddr = cfg.App.ServerHost + fmt.Sprintf(":%d", cfg.Services.AuthPort)
	userAddr = cfg.App.ServerHost + fmt.Sprintf(":%d", cfg.Services.UserPort)
	clients := initialization.InitClients(authAddr, userAddr)

	appContainer := container.NewContainer(clients.AuthClient, clients.UserClient, cfg)

	r := gin.Default()
	config.CORSConfig(r)
	api := r.Group("/api/v1")
	router.AuthRouter(api, cfg, clients.UserClient, appContainer.Auth.Handler)

	r.Run(fmt.Sprintf(":%d", cfg.App.GRPCPort))
}
