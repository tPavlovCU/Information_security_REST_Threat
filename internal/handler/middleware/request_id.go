package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randomBytes := make([]byte, 16)
		_, err := rand.Read(randomBytes)
		if err != nil {
			http.Error(w, "Failed to generate request ID", http.StatusInternalServerError)
			return
		}
		hexString := hex.EncodeToString(randomBytes)
		ctx := context.WithValue(r.Context(), requestIDKey, hexString)
		r = r.WithContext(ctx)
		w.Header().Set("X-Request-ID", hexString)
		next.ServeHTTP(w, r)
	})
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok
}
