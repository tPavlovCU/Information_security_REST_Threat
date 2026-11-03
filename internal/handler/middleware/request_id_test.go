package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_Success(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	header := rec.Header().Get("X-Request-ID")
	if len(header) != 32 {
		t.Errorf("Expected Request ID length 32, got %d", len(header))
	}
}

func TestRequestID_Context(t *testing.T) {
	requestId := ""
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		requestId, ok = GetRequestID(r.Context())
		if !ok {
			t.Error("Request ID not found in context")
		}

		if requestId == "" {
			t.Error("Request ID is empty")
		}
	})
	handler := RequestIDMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	header := rec.Header().Get("X-Request-ID")
	if len(header) != 32 {
		t.Errorf("Expected Request ID length 32, got %d", len(header))
	}

	if header != requestId {
		t.Errorf("Request ID in GetRequestID and in Header are different, GetRequestID: %v, Header: %v", requestId, header)
	}
}

func TestRequestID_Unique(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := RequestIDMiddleware(next)

	req1 := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec1 := httptest.NewRecorder()

	req2 := httptest.NewRequest(http.MethodGet, "/info/test", nil)
	rec2 := httptest.NewRecorder()

	handler.ServeHTTP(rec1, req1)
	header1 := rec1.Header().Get("X-Request-ID")

	handler.ServeHTTP(rec2, req2)
	header2 := rec2.Header().Get("X-Request-ID")

	if header1 == header2 {
		t.Errorf("Headers must be unique")
	}
}
