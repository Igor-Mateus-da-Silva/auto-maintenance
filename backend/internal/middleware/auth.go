package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type authErrorResponse struct {
	Error string `json:"error"`
}

// RequireAuth returns a middleware that validates the HS256 JWT from the HttpOnly cookie
// (or Authorization: Bearer header as fallback), requiring an exp claim.
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""

			if cookie, err := r.Cookie("auth_token"); err == nil && strings.TrimSpace(cookie.Value) != "" {
				tokenStr = strings.TrimSpace(cookie.Value)
			} else if authHeader := r.Header.Get("Authorization"); authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
					respondUnauthorized(w, "formato do token inválido, esperado 'Bearer <token>'")
					return
				}
				tokenStr = strings.TrimSpace(parts[1])
			}

			if tokenStr == "" {
				respondUnauthorized(w, "autenticação necessária")
				return
			}

			token, err := jwt.Parse(
				tokenStr,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				},
				jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
				jwt.WithExpirationRequired(),
			)

			if err != nil || !token.Valid {
				respondUnauthorized(w, "token inválido ou expirado")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func respondUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(authErrorResponse{Error: message})
}
