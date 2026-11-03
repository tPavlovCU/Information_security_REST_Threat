package integration

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"threat-intelligence/internal/domain"
	"threat-intelligence/internal/repository/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	testDBURL := os.Getenv("DATABASE_URL")
	if testDBURL == "" {
		t.Fatal("DATABASE_URL environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, testDBURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("failed to ping test database: %v", err)
	}

	_, err = pool.Exec(ctx, "TRUNCATE TABLE threats RESTART IDENTITY CASCADE")
	if err != nil {
		pool.Close()
		t.Fatalf("failed to clean test database: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestCreate_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "0123456789abcdef0123456789abcdef",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
		Tags:     []string{"test", "malware"},
	}

	err := repo.Create(context.Background(), threat)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	var count int
	err = pool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM threats WHERE md5 = $1",
		threat.MD5,
	).Scan(&count)

	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 record, got %d", count)
	}
}

func TestCreate_AlreadyExists(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "0123456789abcdef0123456789abcdef",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
		Tags:     []string{},
	}

	if err := repo.Create(context.Background(), threat); err != nil {
		t.Fatalf("failed to create initial threat: %v", err)
	}

	err := repo.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrThreatAlreadyExists) {
		t.Fatalf("expected ErrThreatAlreadyExists, got: %v", err)
	}
}

func TestGetByMD5_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "0123456789abcdef0123456789abcdef",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
		Tags:     []string{"test"},
	}

	if err := repo.Create(context.Background(), threat); err != nil {
		t.Fatalf("failed to create threat: %v", err)
	}

	got, err := repo.GetByMD5(context.Background(), threat.MD5)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if got.MD5 != threat.MD5 {
		t.Errorf("expected MD5 %s, got %s", threat.MD5, got.MD5)
	}

	if got.Name != threat.Name {
		t.Errorf("expected name %s, got %s", threat.Name, got.Name)
	}

	if got.ID == 0 {
		t.Error("expected ID to be assigned by database")
	}
}

func TestGetByMD5_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	_, err := repo.GetByMD5(
		context.Background(),
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	)

	if !errors.Is(err, domain.ErrThreatNotFound) {
		t.Fatalf("expected ErrThreatNotFound, got: %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "0123456789abcdef0123456789abcdef",
		Name:     "Old Name",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityLow,
		Tags:     []string{"old"},
	}

	if err := repo.Create(context.Background(), threat); err != nil {
		t.Fatalf("failed to create threat: %v", err)
	}

	threat.Name = "Updated Name"
	threat.Severity = domain.ThreatSeverityCritical
	threat.Tags = []string{"updated"}

	err := repo.Update(context.Background(), threat)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	got, err := repo.GetByMD5(context.Background(), threat.MD5)
	if err != nil {
		t.Fatalf("failed to get updated threat: %v", err)
	}

	if got.Name != "Updated Name" {
		t.Errorf("expected updated name, got %s", got.Name)
	}

	if got.Severity != domain.ThreatSeverityCritical {
		t.Errorf("expected severity %s, got %s",
			domain.ThreatSeverityCritical, got.Severity)
	}

	if len(got.Tags) != 1 || got.Tags[0] != "updated" {
		t.Errorf("expected updated tags, got %v", got.Tags)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Name:     "Missing Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := repo.Update(context.Background(), threat)

	if !errors.Is(err, domain.ErrThreatNotFound) {
		t.Fatalf("expected ErrThreatNotFound, got: %v", err)
	}
}

func TestCreate_NoTags(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewThreatRepository(pool)

	threat := &domain.Threat{
		MD5:      "0123456789abcdef0123456789abcdef",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	if err := repo.Create(context.Background(), threat); err != nil {
		t.Fatalf("failed to create initial threat: %v", err)
	}

	err := repo.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrThreatAlreadyExists) {
		t.Fatalf("expected ErrThreatAlreadyExists, got: %v", err)
	}
}
