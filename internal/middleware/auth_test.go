package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Generate a valid token
	claims := jwt.MapClaims{
		"sub":  "test-user",
		"role": "admin",
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token, err := GenerateTestJWT(claims)
	assert.NoError(t, err)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware := AuthMiddleware(handler)

	// Create a test request
	req := httptest.NewRequest("POST", "/v1/tenants", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Record the response
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware := AuthMiddleware(handler)

	// Create a test request with an invalid token
	req := httptest.NewRequest("POST", "/v1/tenants", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	// Record the response
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Invalid token\n", rr.Body.String())
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware := AuthMiddleware(handler)

	// Create a test request with no token
	req := httptest.NewRequest("POST", "/v1/tenants", nil)

	// Record the response
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "Missing Authorization header\n", rr.Body.String())
}

func TestAuthMiddleware_RBAC_UserRole(t *testing.T) {
	// Generate a token with user role
	claims := jwt.MapClaims{
		"sub":  "test-user",
		"role": "user", // Non-admin role
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token, err := GenerateTestJWT(claims)
	assert.NoError(t, err)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Apply middleware
	middleware := AuthMiddleware(handler)

	// Create a POST request (should be restricted to admin)
	req := httptest.NewRequest("POST", "/v1/tenants", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Record the response
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Equal(t, "Admin role required\n", rr.Body.String())

	// Test a GET request (should be allowed)
	req = httptest.NewRequest("GET", "/v1/tenants", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code)
}
