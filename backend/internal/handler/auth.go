package handler

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const AuthCookieName = "auth_token"

type AuthHandler struct {
	accessKey string
	jwtSecret string
	isProd    bool
}

func NewAuthHandler(accessKey, jwtSecret string, isProd bool) *AuthHandler {
	return &AuthHandler{
		accessKey: accessKey,
		jwtSecret: jwtSecret,
		isProd:    isProd,
	}
}

type LoginRequest struct {
	AccessKey string `json:"access_key"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type VerifyResponse struct {
	Authenticated bool `json:"authenticated"`
}

// Login handles POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	// Validate access key in constant time using fixed-length SHA-256 digests
	reqHash := sha256.Sum256([]byte(req.AccessKey))
	expectedHash := sha256.Sum256([]byte(h.accessKey))
	keyMatch := subtle.ConstantTimeCompare(reqHash[:], expectedHash[:]) == 1
	if !keyMatch {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "chave de acesso inválida"})
		return
	}

	// Generate JWT signed with HS256, 7-day expiration
	ttl := 7 * 24 * time.Hour
	claims := jwt.MapClaims{
		"sub": "admin",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "falha ao gerar token de autenticação"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(VerifyResponse{Authenticated: true})
}

// Logout handles POST /api/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(VerifyResponse{Authenticated: false})
}

// Verify handles GET /api/auth/verify (requires RequireAuth middleware)
func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(VerifyResponse{Authenticated: true})
}
