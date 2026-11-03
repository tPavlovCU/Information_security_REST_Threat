package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingRecoveryMiddleware_Panic(t *testing.T) {
	var logBuffer bytes.Buffer

	logger := slog.New(slog.NewTextHandler(&logBuffer, nil))

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("unexpected error")
	})

	var handler http.Handler = next
	handler = RecoveryMiddleware(logger, handler)
	handler = LoggingMiddleware(logger, handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	logOutput := logBuffer.String()

	if !strings.Contains(logOutput, "status=500") {
		t.Errorf("expected status=500 in logs, got: %s", logOutput)
	}

	// Проверяем, что RecoveryMiddleware записал panic.
	if !strings.Contains(logOutput, "unexpected error") {
		t.Errorf("expected panic message in logs, got: %s", logOutput)
	}
}
