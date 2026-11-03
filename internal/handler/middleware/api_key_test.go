package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	const validAPIKey = "test-secret-key"

	tests := []struct {
		name           string
		providedKey    string
		expectedStatus int
	}{
		{
			name:           "valid API key",
			providedKey:    validAPIKey,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid API key",
			providedKey:    "wrong-key",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing API key",
			providedKey:    "",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			handler := APIKeyMiddleware(validAPIKey, nextHandler)

			req := httptest.NewRequest(http.MethodGet, "/info/test", nil)

			if tt.providedKey != "" {
				req.Header.Set("X-API-Key", tt.providedKey)
			}

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf(
					"expected status %d, got %d",
					tt.expectedStatus,
					rec.Code,
				)
			}

			if tt.expectedStatus == http.StatusOK && !nextCalled {
				t.Error("expected next handler to be called")
			}

			if tt.expectedStatus == http.StatusUnauthorized && nextCalled {
				t.Error("next handler should not be called")
			}
		})
	}
}
