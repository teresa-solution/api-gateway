package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

// JWT secret (in production, use an environment variable or secret manager)
var jwtSecret = []byte("super-secret-jwt-key-2025")

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Return a standard error instead of http.ErrorHandler
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			log.Error().Err(err).Msg("Invalid JWT token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims and check role for specific endpoints
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Role not specified in token", http.StatusForbidden)
			return
		}

		// Restrict POST (create) and DELETE (delete) to admin role
		if (r.Method == http.MethodPost || r.Method == http.MethodDelete) && role != "admin" {
			http.Error(w, "Admin role required", http.StatusForbidden)
			return
		}

		log.Info().Msg("Authentication successful")
		next.ServeHTTP(w, r)
	})
}

// Utility function to generate a JWT token for testing
func GenerateTestJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
