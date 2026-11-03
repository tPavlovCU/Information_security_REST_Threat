package config

import (
	"errors"
	"testing"
	"time"
)

func TestLoad_Success(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/test_db")
	t.Setenv("HTTP_READ_TIMEOUT", "3s")
	t.Setenv("HTTP_WRITE_TIMEOUT", "7s")
	t.Setenv("API_KEY", "smt")
	t.Setenv("RATE_LIMIT", "50")
	t.Setenv("RATE_WINDOW", "2m")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.HTTPPort != "9090" {
		t.Errorf("expected HTTP_PORT 9090, got %s", cfg.HTTPPort)
	}

	if cfg.HTTPReadTimeout != 3*time.Second {
		t.Errorf("expected read timeout 3s, got %v", cfg.HTTPReadTimeout)
	}

	if cfg.HTTPWriteTimeout != 7*time.Second {
		t.Errorf("expected write timeout 7s, got %v", cfg.HTTPWriteTimeout)
	}

	if cfg.APIKey != "smt" {
		t.Errorf("expected API_KEY smt, got %v", cfg.APIKey)
	}

	if cfg.RateLimit != 50 {
		t.Errorf("expected rate limit 50, got %d", cfg.RateLimit)
	}

	if cfg.RateWindow != 2*time.Minute {
		t.Errorf("expected rate window 2m, got %v", cfg.RateWindow)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/test_db")
	t.Setenv("HTTP_READ_TIMEOUT", "")
	t.Setenv("HTTP_WRITE_TIMEOUT", "")
	t.Setenv("API_KEY", "smt")
	t.Setenv("RATE_LIMIT", "")
	t.Setenv("RATE_WINDOW", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.HTTPPort != "8080" {
		t.Errorf("expected default port 8080, got %s", cfg.HTTPPort)
	}

	if cfg.HTTPReadTimeout != 5*time.Second {
		t.Errorf("expected default read timeout 5s, got %v", cfg.HTTPReadTimeout)
	}

	if cfg.HTTPWriteTimeout != 10*time.Second {
		t.Errorf("expected default write timeout 10s, got %v", cfg.HTTPWriteTimeout)
	}

	if cfg.APIKey != "smt" {
		t.Errorf("expected API_KEY smt, got %v", cfg.APIKey)
	}

	if cfg.RateLimit != 100 {
		t.Errorf("expected default rate limit 100, got %d", cfg.RateLimit)
	}

	if cfg.RateWindow != time.Minute {
		t.Errorf("expected default rate window 1m, got %v", cfg.RateWindow)
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("API_KEY", "a")

	_, err := Load()
	if !errors.Is(err, ErrNoDatabaseURL) {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func TestLoad_MissingAPIKey(t *testing.T) {
	t.Setenv("DATABASE_URL", "a")
	t.Setenv("API_KEY", "")

	_, err := Load()
	if !errors.Is(err, ErrNoAPIKey) {
		t.Fatal("expected error when API_KEY is missing")
	}
}

func TestLoad_InvalidReadTimeout(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/test_db")
	t.Setenv("HTTP_READ_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid HTTP_READ_TIMEOUT")
	}
}

func TestLoad_InvalidWriteTimeout(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/test_db")
	t.Setenv("HTTP_WRITE_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid HTTP_WRITE_TIMEOUT")
	}
}
