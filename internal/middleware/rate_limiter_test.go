package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"posgoexpert-rate-limiter/internal/limiter"
	"posgoexpert-rate-limiter/internal/logger"
	"testing"

	"github.com/sirupsen/logrus"
)

func init() {
	logger.Init()
	logger.Logger.SetLevel(logrus.PanicLevel) // Only log panic level during tests
}

// MockStore implements database.DatabaseStore interface for testing
type MockStore struct {
	allowAll bool
}

func (m *MockStore) Incr(ctx context.Context, key string, expirationSeconds int) (int, error) {
	if m.allowAll {
		return 1, nil
	}
	return 1000, nil // Return a high number to simulate rate limit exceeded
}

func (m *MockStore) Get(ctx context.Context, key string) (int, error) {
	if m.allowAll {
		return 100, nil // Default limit for allowed case
	}
	return 5, nil // Lower limit for blocked case
}

func (m *MockStore) TTL(ctx context.Context, key string) (int, error) {
	return 60, nil // Default TTL
}

func (m *MockStore) Set(ctx context.Context, key string, value int, expirationSeconds int) error {
	return nil
}

func BenchmarkRateLimitMiddleware(b *testing.B) {
	tests := []struct {
		name     string
		allowAll bool
	}{
		{"AllowedRequests", true},
		{"BlockedRequests", false},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			store := &MockStore{allowAll: tt.allowAll}
			rateLimiter := limiter.NewLimiter(store, 100, 60) // 100 requests per minute
			middleware := RateLimitMiddleware(rateLimiter)

			// Create a simple handler that will be wrapped
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create the test server
			handler := middleware(nextHandler)

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					req := httptest.NewRequest("GET", "http://example.com", nil)
					req.Header.Set("X-Forwarded-For", "192.168.1.1")
					req.Header.Set("API_KEY", "test-key")

					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, req)
				}
			})
		})
	}
}

// BenchmarkRateLimitMiddlewareWithDifferentIPs tests the middleware performance with different IPs
func BenchmarkRateLimitMiddlewareWithDifferentIPs(b *testing.B) {
	store := &MockStore{allowAll: true}
	rateLimiter := limiter.NewLimiter(store, 100, 60) // 100 requests per minute
	middleware := RateLimitMiddleware(rateLimiter)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware(nextHandler)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			// Generate different IPs for each request
			ip := generateIP(counter)
			counter++

			req := httptest.NewRequest("GET", "http://example.com", nil)
			req.Header.Set("X-Forwarded-For", ip)
			req.Header.Set("API_KEY", "test-key")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
		}
	})
}

// Helper function to generate different IPs
func generateIP(n int) string {
	return fmt.Sprintf("192.168.1.%d", n%255+1)
}
