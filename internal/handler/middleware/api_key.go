package middleware

import (
	"crypto/subtle"
	"net/http"
)

func APIKeyMiddleware(apiKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providedKey := r.Header.Get("X-API-Key")

		if subtle.ConstantTimeCompare(
			[]byte(providedKey),
			[]byte(apiKey),
		) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
