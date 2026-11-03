package middleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, nil))
}

func TestLoggingMiddleware_StatusCode(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	handler := LoggingMiddleware(logger, next)

	req := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "status=404") {
		t.Errorf("Expected status=404 in log, got: %s", logOutput)
	}
}

func TestLoggingMiddleware_RequestID(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	requestID := "test-request-id"

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := LoggingMiddleware(logger, next)

	req := httptest.NewRequest(http.MethodGet, "/info/test", nil)

	ctx := context.WithValue(req.Context(), requestIDKey, requestID)
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "request_id="+requestID) {
		t.Errorf("Expected request_id=%s in log, got: %s", requestID, logOutput)
	}
}

func TestLoggingMiddleware_DefaultStatusCode(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(&buf)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	handler := LoggingMiddleware(logger, next)

	req := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	if !strings.Contains(logOutput, "status=200") {
		t.Errorf("Expected status=200 in log, got: %s", logOutput)
	}
}
