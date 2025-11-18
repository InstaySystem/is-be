package main

import (
	"log"

	_ "github.com/InstaySystem/is-be/docs"
	"github.com/InstaySystem/is-be/internal/config"
	"github.com/InstaySystem/is-be/internal/server"
)

// @title Instay API
// @version 1.0
// @description Instay backend API documentation
// @host localhost:8080
// @BasePath /api/v1

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	sv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Server initialization failed: %v", err)
	}

	ch := make(chan error, 1)
	go func() {
		if err := sv.Start(); err != nil {
			ch <- err
		}
	}()

	log.Printf("Server is running at: http://localhost:%d", cfg.Server.Port)

	sv.GracefulShutdown(ch)
}
