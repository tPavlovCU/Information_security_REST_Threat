package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPPort         string
	DatabaseURL      string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
}

func getEnv(key, fallback string) string {
	data := os.Getenv(key)
	if data == "" {
		data = fallback
	}
	return data
}

var ErrNoDatabaseURL = errors.New("No DATABASE_URL in env")

func Load() (*Config, error) {
	var cfg Config

	cfg.HTTPPort = getEnv("HTTP_PORT", "8080")

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		return nil, ErrNoDatabaseURL
	}
	cfg.DatabaseURL = dbUrl

	readTime, err := time.ParseDuration(getEnv("HTTP_READ_TIMEOUT", "5s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_READ_TIMEOUT: %w", err)
	}

	cfg.HTTPReadTimeout = readTime

	writeTime, err := time.ParseDuration(getEnv("HTTP_WRITE_TIMEOUT", "10s"))
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_WRITE_TIMEOUT: %w", err)
	}

	cfg.HTTPWriteTimeout = writeTime

	return &cfg, nil
}
