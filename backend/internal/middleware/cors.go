package middleware

import (
	"net/http"
	"strings"
)

// CORS validates the request Origin against a comma-separated whitelist of origins.
// Wildcard "*" is intentionally ignored to prevent unsafe credentialed cross-origin access.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	allowedSet := make(map[string]struct{})
	for _, raw := range strings.Split(allowedOrigins, ",") {
		origin := strings.TrimSpace(raw)
		if origin != "" && origin != "*" {
			allowedSet[origin] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")
			if _, ok := allowedSet[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
