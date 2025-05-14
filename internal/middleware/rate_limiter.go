package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"posgoexpert-rate-limiter/internal/limiter"
	"posgoexpert-rate-limiter/internal/logger"
	"strings"
)

func getClientIP(r *http.Request) string {
	// Prioriza IP real de headers (útil com proxies reversos)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func RateLimitMiddleware(limiter *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ip := getClientIP(r)
			token := r.Header.Get("API_KEY")
			result, err := limiter.Check(ctx, ip, token)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("Error checking rate limit: %v", err))
				http.Error(w, "internal limiter error", http.StatusInternalServerError)
				return
			}
			if !result.Allowed {
				http.Error(w, result.Message, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
