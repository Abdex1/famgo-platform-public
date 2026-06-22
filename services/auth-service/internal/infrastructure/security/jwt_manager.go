/*
# PHASE 4 — CREATE JWT TOKEN MANAGER

internal/infrastructure/security/jwt_manager.go

go get github.com/golang-jwt/jwt/v5
*/
package security

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secret string
}

func NewJWTManager(secret string) *JWTManager {
    return &JWTManager{
        secret: secret,
    }
}

func (m *JWTManager) GenerateAccessToken(
    userID string,
    role string,
) (string, error) {

    claims := jwt.MapClaims{
        "sub": userID,
        "role": role,
        "type": "access",
        "exp": time.Now().Add(15 * time.Minute).Unix(),
        "iat": time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    return token.SignedString([]byte(m.secret))
}
