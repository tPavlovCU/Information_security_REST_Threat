package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseWriter_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	rw.WriteHeader(http.StatusNotFound)

	if rw.statusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d",
			http.StatusNotFound, rw.statusCode)
	}

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected response status %d, got %d",
			http.StatusNotFound, rec.Code)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	_, err := rw.Write([]byte("Hello"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if rw.statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d",
			http.StatusOK, rw.statusCode)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected response status %d, got %d",
			http.StatusOK, rec.Code)
	}
}

func TestResponseWriter_WriteHeaderTwice(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec}

	rw.WriteHeader(http.StatusOK)
	rw.WriteHeader(http.StatusNotFound)

	if rw.statusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d",
			http.StatusOK, rw.statusCode)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected response status %d, got %d",
			http.StatusOK, rec.Code)
	}
}
