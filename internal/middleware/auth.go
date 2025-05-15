package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// AuthMiddleware checks for a valid API key in the Authorization header
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Expect "Bearer <api-key>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		apiKey := parts[1]
		// In production, validate apiKey against a secure store (e.g., environment variable or database)
		validAPIKey := "secret-api-key-123" // Replace with secure validation logic
		if apiKey != validAPIKey {
			http.Error(w, "Invalid API key", http.StatusUnauthorized)
			return
		}

		log.Info().Msg("Authentication successful")
		next.ServeHTTP(w, r)
	})
}
