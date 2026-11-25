package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"car-maintenance/internal/middleware"
)

func TestRateLimiter_BlocksExcessiveRequests(t *testing.T) {
	rl := middleware.NewRateLimiter(3, time.Minute)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 3 requests from IP 1 should succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
		req.RemoteAddr = "192.168.1.10:12345"
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, rec.Code)
		}
	}

	// 4th request from IP 1 should be rate-limited (429)
	req := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	req.RemoteAddr = "192.168.1.10:54321"
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("4th request: expected 429, got %d", rec.Code)
	}
	if rec.Header().Get("Retry-After") == "" {
		t.Error("expected Retry-After header on 429 response")
	}

	// Request from a different IP should still succeed
	reqOther := httptest.NewRequest(http.MethodPost, "/api/login", nil)
	reqOther.RemoteAddr = "192.168.1.20:11111"
	recOther := httptest.NewRecorder()
	handler.ServeHTTP(recOther, reqOther)

	if recOther.Code != http.StatusOK {
		t.Errorf("different IP: expected 200, got %d", recOther.Code)
	}
}
