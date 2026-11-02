package main

import (
	"context"
	"log"
	"net/http"

	"threat-intelligence/internal/config"
	"threat-intelligence/internal/database"
	"threat-intelligence/internal/handler"
	"threat-intelligence/internal/repository/postgres"
	"threat-intelligence/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Printf("Cant read config file: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Printf("Cant make database pool: %v", err)
		return
	}

	defer pool.Close()

	threatRepo := postgres.NewThreatRepository(pool)

	threatService := service.NewThreatService(threatRepo)

	threatHandler := handler.NewThreatHandler(threatService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /info/{md5}", threatHandler.GetByMD5)
	mux.HandleFunc("POST /info", threatHandler.Create)
	mux.HandleFunc("PUT /info/{md5}", threatHandler.Update)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
	}

	log.Printf("Server starting on port %s", cfg.HTTPPort)

	if err := server.ListenAndServe(); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}
