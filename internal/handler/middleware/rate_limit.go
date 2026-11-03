package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type client struct {
	requests  int
	resetTime time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*client
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	clients := make(map[string]*client)

	return &RateLimiter{
		clients: clients,
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	clientIP, ok := rl.clients[ip]
	if !ok {
		clientIP = &client{
			requests:  1,
			resetTime: now.Add(rl.window),
		}
		rl.clients[ip] = clientIP
		return true
	}

	if !now.Before(clientIP.resetTime) {
		clientIP.requests = 1
		clientIP.resetTime = (now).Add(rl.window)
		return true
	} else {
		clientIP.requests++
		return clientIP.requests <= rl.limit
	}
}

func RateLimitMiddleware(limiter *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		if !limiter.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
