package service

import (
	"context"
	"errors"
	"testing"

	"threat-intelligence/internal/domain"
)

// Mock Repository

type mockThreatRepository struct {
	createErr error
	getErr    error
	updateErr error

	threat *domain.Threat

	createCalls int
	getCalls    int
	updateCalls int
}

func (r *mockThreatRepository) Create(
	ctx context.Context,
	threat *domain.Threat,
) error {
	r.createCalls++
	return r.createErr
}

func (r *mockThreatRepository) GetByMD5(
	ctx context.Context,
	targetMD5 string,
) (*domain.Threat, error) {
	r.getCalls++
	return r.threat, r.getErr
}

func (r *mockThreatRepository) Update(
	ctx context.Context,
	threat *domain.Threat,
) error {
	r.updateCalls++
	return r.updateErr
}

// GetByMD5 tests

func TestGetByMD5_InvalidMD5(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	_, err := svc.GetByMD5(context.Background(), "12345")

	if !errors.Is(err, domain.ErrBadMD5) {
		t.Errorf("expected error %v, got %v", domain.ErrBadMD5, err)
	}

	if repo.getCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.getCalls)
	}
}

func TestGetByMD5_Success(t *testing.T) {
	expected := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	repo := &mockThreatRepository{
		threat: expected,
	}

	svc := NewThreatService(repo)

	got, err := svc.GetByMD5(context.Background(), expected.MD5)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != expected {
		t.Errorf("expected threat %v, got %v", expected, got)
	}

	if repo.getCalls != 1 {
		t.Errorf("expected 1 repository call, got %d", repo.getCalls)
	}
}

// Create tests

func TestCreate_Success(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Create(context.Background(), threat)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.createCalls != 1 {
		t.Errorf("expected 1 repository call, got %d", repo.createCalls)
	}
}

func TestCreate_InvalidMD5(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "12345",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrBadMD5) {
		t.Errorf("expected error %v, got %v", domain.ErrBadMD5, err)
	}

	if repo.createCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.createCalls)
	}
}

func TestCreate_EmptyName(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "   ",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected error %v, got %v", domain.ErrEmptyName, err)
	}

	if repo.createCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.createCalls)
	}
}

func TestCreate_InvalidType(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     "unknown",
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrBadType) {
		t.Errorf("expected error %v, got %v", domain.ErrBadType, err)
	}

	if repo.createCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.createCalls)
	}
}

func TestCreate_InvalidSeverity(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: "extreme",
	}

	err := svc.Create(context.Background(), threat)

	if !errors.Is(err, domain.ErrBadSeverity) {
		t.Errorf("expected error %v, got %v", domain.ErrBadSeverity, err)
	}

	if repo.createCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.createCalls)
	}
}

// Update tests

func TestUpdate_Success(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Updated Threat",
		Type:     domain.ThreatTypeRansomware,
		Severity: domain.ThreatSeverityCritical,
	}

	err := svc.Update(context.Background(), threat)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.updateCalls != 1 {
		t.Errorf("expected 1 repository call, got %d", repo.updateCalls)
	}
}

func TestUpdate_InvalidMD5(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "12345",
		Name:     "Updated Threat",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Update(context.Background(), threat)

	if !errors.Is(err, domain.ErrBadMD5) {
		t.Errorf("expected error %v, got %v", domain.ErrBadMD5, err)
	}

	if repo.updateCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.updateCalls)
	}
}

func TestUpdate_EmptyName(t *testing.T) {
	repo := &mockThreatRepository{}

	svc := NewThreatService(repo)

	threat := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "",
		Type:     domain.ThreatTypeTrojan,
		Severity: domain.ThreatSeverityHigh,
	}

	err := svc.Update(context.Background(), threat)

	if !errors.Is(err, domain.ErrEmptyName) {
		t.Errorf("expected error %v, got %v", domain.ErrEmptyName, err)
	}

	if repo.updateCalls != 0 {
		t.Errorf("expected 0 repository calls, got %d", repo.updateCalls)
	}
}
