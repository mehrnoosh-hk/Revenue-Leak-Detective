package main

import (
	"log"
	"net/http"
	"rld/services/api/config"
	"rld/services/api/handlers"
)

type APIService struct {
	Config   *config.Config
	Servemux *http.ServeMux
}

func main() {
	// Initialize configuration
	cfg := &config.Config{
		Port: "8080",
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/health", handlers.HealthCheckHandler)

	// Create API service instance
	apiService := &APIService{
		Config:   cfg,
		Servemux: mux,
	}

	// Start the HTTP server
	err := http.ListenAndServe(":"+apiService.Config.Port, apiService.Servemux)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
