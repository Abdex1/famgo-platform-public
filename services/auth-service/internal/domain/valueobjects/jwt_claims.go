/*
# STEP 8 — CREATE JWT CLAIMS MODEL

services/auth-service/internal/domain/valueobjects/jwt_claims.go


# FILE: jwt_claims.go
*/
package valueobjects

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID string `json:"user_id"`

	Role string `json:"role"`

	SessionID string `json:"session_id"`

	DeviceID string `json:"device_id"`

	jwt.RegisteredClaims
}
