package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"threat-intelligence/internal/config"
	"threat-intelligence/internal/database"
	"threat-intelligence/internal/handler"
	"threat-intelligence/internal/handler/middleware"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

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
		Handler:      middleware.RequestIDMiddleware(mux),
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
	}

	log.Printf("Server starting on port %s", cfg.HTTPPort)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
		return
	case <-ctx.Done():
		log.Printf("Shutdown signal received")
		shutDownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutDownCtx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)
			if err := server.Close(); err != nil {
				log.Printf("Force close failed: %v", err)
			}
		}

		err := <-serverErrors
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server error during shutdown: %v", err)
		}

		log.Println("Server stopped")
	}
}
