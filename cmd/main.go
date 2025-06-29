package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alirezaghasemi/user-manager/internal/config"
)

func main() {
	// make and run server
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading env variable: %v", err)
	}

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: nil,
	}

	fmt.Printf("Server started on %d", cfg.Server.Port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
