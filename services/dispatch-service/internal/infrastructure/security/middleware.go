package security

import (
	"context"
	"net/http"
	"strings"
	"time"
)

const permissionDispatchWrite = "dispatch:write"
const permissionDispatchRead = "dispatch:read"

type RateLimiter interface {
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}

type HTTPSecurityMiddleware struct {
	validator *TokenValidator
	limiter   RateLimiter
}

func NewHTTPSecurityMiddleware(validator *TokenValidator, limiter RateLimiter) *HTTPSecurityMiddleware {
	return &HTTPSecurityMiddleware{validator: validator, limiter: limiter}
}

func (m *HTTPSecurityMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if m.limiter != nil {
			clientKey := clientRateLimitKey(r)
			allowed, err := m.limiter.Allow(r.Context(), clientKey, 120, time.Minute)
			if err != nil || !allowed {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
		}

		claims, err := m.validator.Validate(r.Context(), r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		required := permissionDispatchRead
		if r.Method != http.MethodGet {
			required = permissionDispatchWrite
		}
		if !m.validator.HasPermission(claims, required) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isPublicPath(path string) bool {
	return path == "/v1/health" || path == "/healthz" || path == "/metrics" || path == "/version"
}

func clientRateLimitKey(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if auth != "" {
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) == 2 {
			return "auth:" + parts[1][:min(16, len(parts[1]))]
		}
	}
	return "ip:" + r.RemoteAddr
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
