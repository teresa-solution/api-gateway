package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitMiddleware limits requests per client
func RateLimitMiddleware(next http.Handler) http.Handler {
	// Map to store limiter per client IP
	limiters := make(map[string]*rate.Limiter)
	mu := sync.Mutex{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use client IP as the key (simplified; in production, consider X-Forwarded-For)
		clientIP := r.RemoteAddr

		mu.Lock()
		limiter, exists := limiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(rate.Every(time.Second*10), 100) // 100 requests per 10 seconds
			limiters[clientIP] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
