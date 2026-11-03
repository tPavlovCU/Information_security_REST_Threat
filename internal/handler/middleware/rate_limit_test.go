package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_FirstRequest(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	if !limiter.Allow("192.168.1.1") {
		t.Error("first request should be allowed")
	}
}

func TestRateLimiter_WithinLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !limiter.Allow("192.168.1.1") {
			t.Errorf("request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_ExceedsLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		limiter.Allow("192.168.1.1")
	}

	if limiter.Allow("192.168.1.1") {
		t.Error("fourth request should be rejected")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)

	if !limiter.Allow("192.168.1.1") {
		t.Error("first IP request should be allowed")
	}

	if !limiter.Allow("192.168.1.2") {
		t.Error("second IP request should be allowed")
	}
}

func TestRateLimiter_ResetWindow(t *testing.T) {
	limiter := NewRateLimiter(1, 50*time.Millisecond)

	if !limiter.Allow("192.168.1.1") {
		t.Fatal("first request should be allowed")
	}

	if limiter.Allow("192.168.1.1") {
		t.Fatal("second request should be rejected")
	}

	time.Sleep(60 * time.Millisecond)

	if !limiter.Allow("192.168.1.1") {
		t.Error("request should be allowed after window reset")
	}
}
