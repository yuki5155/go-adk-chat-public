package main

import (
	"fmt"
	"log"

	"github.com/yuki5155/go-google-auth/internal/infrastructure/config"
	"github.com/yuki5155/go-google-auth/internal/infrastructure/container"
	"github.com/yuki5155/go-google-auth/internal/presentation/http/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create dependency injection container
	c := container.NewContainer(cfg)

	// Setup router with container
	r := router.Setup(c)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
