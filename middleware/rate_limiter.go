package middleware

import (
	"luxsuv-v4/models"
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     time.Duration // Time between requests
	burst    int           // Maximum burst size
}

// Visitor represents a client with rate limiting info
type Visitor struct {
	limiter  chan struct{}
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: minimum time between requests (e.g., 100ms for 10 req/sec)
// burst: maximum number of requests that can be made in a burst
func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
	}

	// Clean up old visitors every minute
	go rl.cleanupVisitors()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		// Create new visitor with full bucket
		visitor = &Visitor{
			limiter:  make(chan struct{}, rl.burst),
			lastSeen: time.Now(),
		}
		// Fill the bucket
		for i := 0; i < rl.burst; i++ {
			visitor.limiter <- struct{}{}
		}
		rl.visitors[ip] = visitor
	}

	visitor.lastSeen = time.Now()

	// Try to consume a token
	select {
	case <-visitor.limiter:
		// Token available, allow request
		// Refill token after rate duration
		go func() {
			time.Sleep(rl.rate)
			select {
			case visitor.limiter <- struct{}{}:
			default:
				// Bucket is full, ignore
			}
		}()
		return true
	default:
		// No tokens available, deny request
		return false
	}
}

// cleanupVisitors removes old visitors to prevent memory leaks
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			if time.Since(visitor.lastSeen) > time.Hour {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware applies rate limiting to requests
func (rl *RateLimiter) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get client IP
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			ip = realIP
		}

		if !rl.Allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			response := models.ErrorResponse{Error: "too many requests"}
			w.Write([]byte(`{"error":"too many requests"}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}