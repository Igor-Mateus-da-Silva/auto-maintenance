package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rateLimitEntry struct {
	count   int
	resetAt time.Time
}

type RateLimiter struct {
	mu          sync.Mutex
	entries     map[string]*rateLimitEntry
	maxRequests int
	window      time.Duration
}

// NewRateLimiter creates an in-memory fixed-window rate limiter keyed by client IP.
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		entries:     make(map[string]*rateLimitEntry),
		maxRequests: maxRequests,
		window:      window,
	}
}

func (rl *RateLimiter) allow(ip string, now time.Time) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Lazy cleanup of expired entries when map grows
	if len(rl.entries) > 1000 {
		for k, v := range rl.entries {
			if now.After(v.resetAt) {
				delete(rl.entries, k)
			}
		}
	}

	entry, exists := rl.entries[ip]
	if !exists || now.After(entry.resetAt) {
		rl.entries[ip] = &rateLimitEntry{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return true, 0
	}

	if entry.count >= rl.maxRequests {
		retryAfter := entry.resetAt.Sub(now)
		if retryAfter < time.Second {
			retryAfter = time.Second
		}
		return false, retryAfter
	}

	entry.count++
	return true, 0
}

// Middleware wraps an http.Handler with IP-based rate limiting.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractClientIP(r)
		allowed, retryAfter := rl.allow(ip, time.Now())
		if !allowed {
			seconds := int(retryAfter.Seconds())
			if seconds < 1 {
				seconds = 1
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", strconv.Itoa(seconds))
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(authErrorResponse{
				Error: "muitas tentativas de login, tente novamente mais tarde",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func extractClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
