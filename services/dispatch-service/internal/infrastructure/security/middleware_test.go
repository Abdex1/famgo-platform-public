package security

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type stubLimiter struct {
	allowed bool
}

func (s stubLimiter) Allow(context.Context, string, int, time.Duration) (bool, error) {
	return s.allowed, nil
}

func signToken(t *testing.T, secret, issuer string, scopes []string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: "user-1",
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: issuer,
		},
	})
	signed, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return signed
}

func TestHTTPSecurityMiddlewareAllowsReadWithScope(t *testing.T) {
	secret := "test-secret-key-min-32-characters"
	validator := NewTokenValidator(secret, "famgo-platform")
	mw := NewHTTPSecurityMiddleware(validator, stubLimiter{allowed: true})

	called := false
	handler := mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/v1/dispatch/match", nil)
	req.Header.Set("Authorization", "Bearer "+signToken(t, secret, "famgo-platform", []string{"dispatch:read"}))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHTTPSecurityMiddlewareRejectsMissingToken(t *testing.T) {
	mw := NewHTTPSecurityMiddleware(NewTokenValidator("secret", ""), stubLimiter{allowed: true})
	req := httptest.NewRequest(http.MethodPost, "/v1/dispatch/match", nil)
	rec := httptest.NewRecorder()
	mw.Wrap(http.HandlerFunc(http.NotFound)).ServeHTTP(rec, req)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHTTPSecurityMiddlewareSkipsPublicPaths(t *testing.T) {
	mw := NewHTTPSecurityMiddleware(NewTokenValidator("secret", ""), stubLimiter{allowed: false})
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	rec := httptest.NewRecorder()
	mw.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
