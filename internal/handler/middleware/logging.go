package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(rw, r)

		duration := time.Since(now)
		method := r.Method
		path := r.URL.Path
		code := rw.statusCode
		requestID, ok := GetRequestID(r.Context())
		if code == 0 {
			code = http.StatusOK
		}

		if !ok {
			requestID = "No request ID"
		}

		logger.Info("HTTP request",
			slog.String("method", method),
			slog.String("path", path),
			slog.String("request_id", requestID),
			slog.Duration("duration", duration),
			slog.Int("status", code))
	})
}
