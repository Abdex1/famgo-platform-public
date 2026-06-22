/*
# PHASE 9 — CREATE RBAC MIDDLEWARE

internal/interfaces/rest/middleware/rbac.go
*/
package middleware

import (
    "net/http"

    "github.com/golang-jwt/jwt/v5"
)

func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

            claims, ok := r.Context().Value("claims").(jwt.MapClaims)

            if !ok {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            userRole := claims["role"]

            if userRole != role {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
