package main

import (
	"context"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err != nil {
		logger.Warn("No .env file found", slog.Any("err", err))
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", slog.Any("err", err))
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
		logger.Error("Failed to create database pool", slog.Any("err", err))
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

	var handler http.Handler = mux
	limiter := middleware.NewRateLimiter(cfg.RateLimit, cfg.RateWindow)

	handler = middleware.RequestIDMiddleware(handler)
	handler = middleware.LoggingMiddleware(logger, handler)
	handler = middleware.APIKeyMiddleware(cfg.APIKey, handler)
	handler = middleware.RecoveryMiddleware(logger, handler)
	handler = middleware.RateLimitMiddleware(limiter, handler)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
	}

	logger.Info("Server starting", slog.String("port", cfg.HTTPPort))

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.Any("err", err))
		}
		return
	case <-ctx.Done():
		logger.Info("Shutdown signal received")
		shutDownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutDownCtx); err != nil {
			logger.Error("Graceful shutdown failed", slog.Any("err", err))
			if err := server.Close(); err != nil {
				logger.Error("Force close failed", slog.Any("err", err))
			}
		}

		err := <-serverErrors
		if err != nil && err != http.ErrServerClosed {
			logger.Error("Server error during shutdown", slog.Any("err", err))
		}

		logger.Info("Server stopped")
	}
}
