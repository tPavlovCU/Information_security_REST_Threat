package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//"threat-intelligence/internal/service"

	"threat-intelligence/internal/domain"
)

type mockThreatService struct {
	threat       *domain.Threat
	getByMD5Err  error
	requestedMD5 string

	createErr     error
	createdThreat *domain.Threat
	createCalls   int

	updateErr     error
	updatedThreat *domain.Threat
	updateCalls   int
}

func (s *mockThreatService) Create(ctx context.Context, threat *domain.Threat) error {
	s.createCalls++
	s.createdThreat = threat
	return s.createErr
}

func (s *mockThreatService) GetByMD5(ctx context.Context, targetMD5 string) (*domain.Threat, error) {
	s.requestedMD5 = targetMD5
	return s.threat, s.getByMD5Err
}

func (s *mockThreatService) Update(ctx context.Context, threat *domain.Threat) error {
	s.updateCalls++
	s.updatedThreat = threat
	return s.updateErr
}

func TestGetByMD5_InvalidMD5(t *testing.T) {
	mock := &mockThreatService{
		threat:       nil,
		getByMD5Err:  domain.ErrBadMD5,
		requestedMD5: ""}
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
	mock := &mockThreatService{
		threat:       nil,
		getByMD5Err:  domain.ErrThreatNotFound,
		requestedMD5: ""}
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
	mock := &mockThreatService{
		threat: &domain.Threat{
			MD5:      "d41d8cd98f00b204e9800998ecf8427e",
			Name:     "Test Threat",
			Type:     "trojan",
			Severity: "high",
		},
		getByMD5Err:  nil,
		requestedMD5: ""}

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

func TestCreate_Success(t *testing.T) {

	mock := &mockThreatService{}

	createBody := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     "trojan",
		Severity: "high",
	}

	body, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /info", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/info", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)
	var got CreateResult

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
		return
	}

	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response")
	}

	if got.MD5 != createBody.MD5 {
		t.Errorf("expected MD5 %s, got %s", createBody.MD5, got.MD5)
	}

	if got.Name != createBody.Name {
		t.Errorf("expected Name %s, got %s", createBody.Name, got.Name)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %s", contentType)
	}

	if mock.createCalls != 1 {
		t.Errorf("expected Create to be called once, got %d", mock.createCalls)
	}

	if mock.createdThreat == nil {
		t.Fatal("expected service to receive a threat, got nil")
	}

	if mock.createdThreat.MD5 != createBody.MD5 {
		t.Errorf("service received wrong MD5")
	}
}

func TestCreate_InvalidJSON(t *testing.T) {

	mock := &mockThreatService{}

	body := []byte(`{"md5": invalid json}`)

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /info", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/info", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if mock.createCalls != 0 {
		t.Errorf("expected 0 calls to Create, got %d", mock.createCalls)
	}
}

func TestCreate_ThreatAlreadyExists(t *testing.T) {
	mock := &mockThreatService{
		createErr: domain.ErrThreatAlreadyExists,
	}

	createBody := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     "trojan",
		Severity: "high",
	}

	body, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /info", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/info", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, rec.Code)
	}

	if mock.createCalls != 1 {
		t.Errorf("expected 1 call to Create, got %d", mock.createCalls)
	}
}

func TestCreate_InternalError(t *testing.T) {
	mock := &mockThreatService{
		createErr: errors.New("db connection failed"),
	}

	createBody := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "Test Threat",
		Type:     "trojan",
		Severity: "high",
	}

	body, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /info", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/info", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	if mock.createCalls != 1 {
		t.Errorf("expected 1 call to Create, got %d", mock.createCalls)
	}
}

func TestCreate_EmptyName(t *testing.T) {
	mock := &mockThreatService{
		createErr: domain.ErrEmptyName,
	}

	createBody := &domain.Threat{
		MD5:      "d41d8cd98f00b204e9800998ecf8427e",
		Name:     "",
		Type:     "trojan",
		Severity: "high",
	}

	body, err := json.Marshal(createBody)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}

	handler := NewThreatHandler(mock)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /info", handler.Create)

	req := httptest.NewRequest(http.MethodPost, "/info", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if mock.createCalls != 1 {
		t.Errorf("expected 1 call to Create, got %d", mock.createCalls)
	}
}

func TestUpdate_Success(t *testing.T) {
	mock := &mockThreatService{}

	md5 := "d41d8cd98f00b204e9800998ecf8427e"

	updateBody := `{
        "name": "Updated Threat",
        "type": "ransomware",
        "severity": "critical"
    }`

	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /info/{md5}", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/info/"+md5,
		strings.NewReader(updateBody),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d",
			http.StatusOK, rec.Code)
	}

	if mock.updateCalls != 1 {
		t.Errorf("expected 1 Update call, got %d", mock.updateCalls)
	}

	if mock.updatedThreat == nil {
		t.Fatal("expected updated threat, got nil")
	}

	if mock.updatedThreat.MD5 != md5 {
		t.Errorf("expected MD5 %s, got %s",
			md5, mock.updatedThreat.MD5)
	}

	if mock.updatedThreat.Name != "Updated Threat" {
		t.Errorf("expected name Updated Threat, got %s",
			mock.updatedThreat.Name)
	}

	if mock.updatedThreat.Type != "ransomware" {
		t.Errorf("expected type ransomware, got %s",
			mock.updatedThreat.Type)
	}

	if mock.updatedThreat.Severity != "critical" {
		t.Errorf("expected severity critical, got %s",
			mock.updatedThreat.Severity)
	}
}

func TestUpdate_InvalidJSON(t *testing.T) {
	mock := &mockThreatService{}

	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /info/{md5}", handler.Update)

	req := httptest.NewRequest(
		http.MethodPut,
		"/info/d41d8cd98f00b204e9800998ecf8427e",
		strings.NewReader(`{"name": invalid json}`),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}

	if mock.updateCalls != 0 {
		t.Errorf("expected 0 Update calls, got %d", mock.updateCalls)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	mock := &mockThreatService{
		updateErr: domain.ErrThreatNotFound,
	}

	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /info/{md5}", handler.Update)

	body := `{
        "name": "Updated Threat",
        "type": "trojan",
        "severity": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/info/d41d8cd98f00b204e9800998ecf8427e",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d",
			http.StatusNotFound, rec.Code)
	}

	if mock.updateCalls != 1 {
		t.Errorf("expected 1 Update call, got %d", mock.updateCalls)
	}
}

func TestUpdate_InvalidMD5(t *testing.T) {
	mock := &mockThreatService{
		updateErr: domain.ErrBadMD5,
	}

	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /info/{md5}", handler.Update)

	body := `{
        "name": "Updated Threat",
        "type": "trojan",
        "severity": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/info/12345",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d",
			http.StatusBadRequest, rec.Code)
	}
}

func TestUpdate_InternalError(t *testing.T) {
	mock := &mockThreatService{
		updateErr: errors.New("database connection failed"),
	}

	handler := NewThreatHandler(mock)

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /info/{md5}", handler.Update)

	body := `{
        "name": "Updated Threat",
        "type": "trojan",
        "severity": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/info/d41d8cd98f00b204e9800998ecf8427e",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d",
			http.StatusInternalServerError, rec.Code)
	}

	if strings.Contains(rec.Body.String(), "database connection failed") {
		t.Error("internal error details should not be exposed to client")
	}
}
