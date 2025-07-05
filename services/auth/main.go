package main

import (
	"fmt"
	"log"

	"github.com/SomeHowMicroservice/shm-be/services/auth/config"
	"github.com/SomeHowMicroservice/shm-be/services/auth/initialization"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Tải cấu hình Auth Service thất bại: %v", err)
	}

	rdb, err := initialization.InitCache(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Listening on:", cfg.App.GRPCPort)
	fmt.Println("Cache host:", cfg.Cache.CHost)
	fmt.Println("Cache", rdb)
}
