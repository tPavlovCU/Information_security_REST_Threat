package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort         string
	DatabaseURL      string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	APIKey           string

	RateLimit  int
	RateWindow time.Duration
}

func getEnv(key, fallback string) string {
	data := os.Getenv(key)
	if data == "" {
		data = fallback
	}
	return data
}

var ErrNoDatabaseURL = errors.New("No DATABASE_URL in env")
var ErrNoAPIKey = errors.New("No API_KEY in env")

var ErrBadRateLimit = errors.New("Bad RATE_LIMIT in env")
var ErrBadRateWindow = errors.New("Bad RATE_WINDOW in env")

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

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, ErrNoAPIKey
	}
	cfg.APIKey = apiKey

	cfg.RateLimit, err = strconv.Atoi(getEnv("RATE_LIMIT", "100"))

	if err != nil || cfg.RateLimit <= 0 {
		return nil, ErrBadRateLimit
	}

	cfg.RateWindow, err = time.ParseDuration(getEnv("RATE_WINDOW", "1m"))
	if err != nil || cfg.RateWindow <= 0 {
		return nil, ErrBadRateWindow
	}

	return &cfg, nil
}
