package handler

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	accessKey string
	jwtSecret string
}

func NewAuthHandler(accessKey, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		accessKey: accessKey,
		jwtSecret: jwtSecret,
	}
}

type LoginRequest struct {
	AccessKey string `json:"access_key"`
}

type LoginResponse struct {
	Token string `json:"token"`
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "corpo da requisição inválido"})
		return
	}

	// Validate access key in constant time
	keyMatch := subtle.ConstantTimeCompare([]byte(req.AccessKey), []byte(h.accessKey)) == 1
	if !keyMatch {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "chave de acesso inválida"})
		return
	}

	// Generate JWT signed with HS256, 7-day expiration
	claims := jwt.MapClaims{
		"sub": "admin",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "falha ao gerar token de autenticação"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}

// Verify handles GET /api/auth/verify (requires RequireAuth middleware)
func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(VerifyResponse{Authenticated: true})
}
