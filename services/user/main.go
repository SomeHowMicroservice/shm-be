package main

import (
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/services/user/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình thất bại: %v", err)
	}

	fmt.Println("Listening on:", cfg.App.GRPCPort)
	fmt.Println("DB host:", cfg.Database.DBHost)
}