/*
# STEP 9 — CREATE JWT SERVICE

services/auth-service/internal/domain/services/jwt_service.go


# FILE: jwt_service.go
*/
package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"famgo/auth-service/internal/domain/valueobjects"
)

type JWTService struct {
	secret string
	issuer string
	audience string
}

func NewJWTService(
	secret string,
	issuer string,
	audience string,
) *JWTService {
	return &JWTService{
		secret: secret,
		issuer: issuer,
		audience: audience,
	}
}

func (s *JWTService) GenerateAccessToken(
	userID string,
	role string,
	sessionID string,
	deviceID string,
) (string, error) {

	claims := valueobjects.Claims{
		UserID: userID,
		Role: role,
		SessionID: sessionID,
		DeviceID: deviceID,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: s.issuer,
			Audience: []string{s.audience},
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(s.secret))
}
