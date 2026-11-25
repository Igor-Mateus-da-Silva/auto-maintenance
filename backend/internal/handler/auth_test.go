package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"car-maintenance/internal/handler"
	"car-maintenance/internal/middleware"
)

func TestAuth_LoginAndVerify(t *testing.T) {
	accessKey := "test-secret-key"
	jwtSecret := "test-jwt-secret-very-secure-32-chars"

	authHandler := handler.NewAuthHandler(accessKey, jwtSecret, false)
	requireAuth := middleware.RequireAuth(jwtSecret)

	// 1. Test Login with wrong key -> 401
	wrongBody, _ := json.Marshal(map[string]string{"access_key": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(wrongBody))
	rr := httptest.NewRecorder()
	authHandler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for invalid key, got %d", rr.Code)
	}

	var errResp handler.ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&errResp); err != nil || errResp.Error != "chave de acesso inválida" {
		t.Errorf("Expected error 'chave de acesso inválida', got '%s'", errResp.Error)
	}

	// 2. Test Login with correct key -> 200 with HttpOnly cookie
	validBody, _ := json.Marshal(map[string]string{"access_key": accessKey})
	req = httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewReader(validBody))
	rr = httptest.NewRecorder()
	authHandler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for valid key, got %d", rr.Code)
	}

	var loginResp handler.VerifyResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil || !loginResp.Authenticated {
		t.Fatalf("Expected authenticated: true, got %+v", loginResp)
	}

	var authCookie *http.Cookie
	for _, c := range rr.Result().Cookies() {
		if c.Name == handler.AuthCookieName {
			authCookie = c
			break
		}
	}
	if authCookie == nil || authCookie.Value == "" {
		t.Fatal("Expected auth_token cookie to be set")
	}
	if !authCookie.HttpOnly {
		t.Error("Expected auth_token cookie to have HttpOnly=true")
	}

	// 3. Test Verify route without cookie/token -> 401
	verifyHandler := requireAuth(http.HandlerFunc(authHandler.Verify))
	req = httptest.NewRequest(http.MethodGet, "/api/auth/verify", nil)
	rr = httptest.NewRecorder()
	verifyHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for request without token, got %d", rr.Code)
	}

	// 4. Test Verify route with invalid cookie -> 401
	req = httptest.NewRequest(http.MethodGet, "/api/auth/verify", nil)
	req.AddCookie(&http.Cookie{Name: handler.AuthCookieName, Value: "invalid-token"})
	rr = httptest.NewRecorder()
	verifyHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for request with invalid cookie, got %d", rr.Code)
	}

	// 5. Test Verify route with valid HttpOnly cookie -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/auth/verify", nil)
	req.AddCookie(authCookie)
	rr = httptest.NewRecorder()
	verifyHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for valid cookie, got %d", rr.Code)
	}

	var verifyResp handler.VerifyResponse
	if err := json.NewDecoder(rr.Body).Decode(&verifyResp); err != nil || !verifyResp.Authenticated {
		t.Errorf("Expected authenticated: true, got %+v", verifyResp)
	}

	// 6. Test Logout clears cookie
	req = httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	rr = httptest.NewRecorder()
	authHandler.Logout(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 on logout, got %d", rr.Code)
	}
	for _, c := range rr.Result().Cookies() {
		if c.Name == handler.AuthCookieName && c.MaxAge >= 0 {
			t.Errorf("Expected logout to expire auth_token cookie, got MaxAge=%d", c.MaxAge)
		}
	}
}
