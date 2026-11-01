package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	//"threat-intelligence/internal/service"

	"threat-intelligence/internal/domain"
)

type mockThreatService struct {
	threat       *domain.Threat
	getByMD5Err  error
	requestedMD5 string
}

func (s *mockThreatService) Create(ctx context.Context, threat *domain.Threat) error {
	return nil
}

func (s *mockThreatService) GetByMD5(ctx context.Context, targetMD5 string) (*domain.Threat, error) {
	s.requestedMD5 = targetMD5
	return s.threat, s.getByMD5Err
}

func (s *mockThreatService) Update(ctx context.Context, threat *domain.Threat) error {
	return nil
}

func TestGetByMD5_InvalidMD5(t *testing.T) {
	mock := &mockThreatService{nil, domain.ErrBadMD5, ""}
	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /info/{md5}", handler.GetByMD5)

	req := httptest.NewRequest(http.MethodGet, "/info/12345", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetByMD5_NotFound(t *testing.T) {
	mock := &mockThreatService{nil, domain.ErrThreatNotFound, ""}
	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /info/{md5}", handler.GetByMD5)

	req := httptest.NewRequest(http.MethodGet, "/info/12345", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetByMD5_Success(t *testing.T) {
	mock := &mockThreatService{&domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     "trojan",
		Severity: "high",
	}, nil, ""}

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /info/{md5}", handler.GetByMD5)
	req := httptest.NewRequest(http.MethodGet, "/info/d41d8cd98f00b204e9800998ecf8427e", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	var got domain.Threat
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.MD5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("unexpected MD5: %s", got.MD5)
	}

	if got.Name != "Test Threat" {
		t.Errorf("unexpected name: %s", got.Name)
	}

	if got.Type != "trojan" {
		t.Errorf("unexpected Type: %s", got.Type)
	}

	if got.Severity != "high" {
		t.Errorf("unexpected severity: %s", got.Severity)
	}

	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("unexpected Content-Type: %s", got)
	}

	if mock.requestedMD5 != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("unexpected rrequested MD5: %s", mock.requestedMD5)
	}
}
