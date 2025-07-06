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
	//userAddr = "localhost:8082"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Gateway thất bại: %v", err)
	}

	authAddr = cfg.App.ServerHost + fmt.Sprintf(":%d", cfg.Services.AuthPort)
	clients := initialization.InitClients(authAddr)

	appContainer := container.NewContainer(clients.AuthClient)

	r := gin.Default()
	api := r.Group("/api/v1")
	router.AuthRouter(api, clients.AuthClient, appContainer.Auth.Handler)

	r.Run(fmt.Sprintf(":%d", cfg.App.GRPCPort))
}
