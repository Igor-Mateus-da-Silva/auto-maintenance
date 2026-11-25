package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"car-maintenance/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func TestCORS_WhitelistAndNoWildcard(t *testing.T) {
	cors := middleware.CORS("http://localhost:4200, https://app.example.com, *")
	handler := cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Allowed origin
	req := httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	req.Header.Set("Origin", "https://app.example.com")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
		t.Errorf("expected allowed origin, got %q", got)
	}
	if rec.Header().Get("Vary") != "Origin" {
		t.Errorf("expected Vary: Origin, got %q", rec.Header().Get("Vary"))
	}

	// Disallowed origin (must not receive ACAO header)
	reqBad := httptest.NewRequest(http.MethodGet, "/api/vehicles", nil)
	reqBad.Header.Set("Origin", "https://evil.example.com")
	recBad := httptest.NewRecorder()
	handler.ServeHTTP(recBad, reqBad)

	if got := recBad.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("expected empty ACAO for disallowed origin, got %q", got)
	}
}

func TestSecurityHeaders(t *testing.T) {
	sec := middleware.SecurityHeaders(true)
	handler := sec(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("missing X-Content-Type-Options: nosniff")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("missing X-Frame-Options: DENY")
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Error("missing Cache-Control: no-store")
	}
	if rec.Header().Get("Strict-Transport-Security") == "" {
		t.Error("expected HSTS header when isProd=true")
	}
}

func TestRequireAuth_RejectsHS512AndMissingExp(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef"
	authMw := middleware.RequireAuth(secret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Token signed with HS512 instead of HS256
	hs512Token, _ := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": "admin",
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(secret))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer "+hs512Token)
	rec := httptest.NewRecorder()
	authMw.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for HS512 token, got %d", rec.Code)
	}

	// HS256 token without exp claim
	noExpToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "admin",
	}).SignedString([]byte(secret))

	reqNoExp := httptest.NewRequest(http.MethodGet, "/api/auth/verify", nil)
	reqNoExp.Header.Set("Authorization", "Bearer "+noExpToken)
	recNoExp := httptest.NewRecorder()
	authMw.ServeHTTP(recNoExp, reqNoExp)
	if recNoExp.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 for token without exp claim, got %d", recNoExp.Code)
	}
}
