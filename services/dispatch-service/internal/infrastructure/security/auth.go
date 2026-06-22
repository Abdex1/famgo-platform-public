package security

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingToken = errors.New("missing authorization token")
	ErrInvalidToken = errors.New("invalid authorization token")
)

type Claims struct {
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	Scopes   []string `json:"scopes"`
	jwt.RegisteredClaims
}

type TokenValidator struct {
	secret []byte
	issuer string
}

func NewTokenValidator(secret, issuer string) *TokenValidator {
	return &TokenValidator{secret: []byte(secret), issuer: issuer}
}

func (v *TokenValidator) Validate(ctx context.Context, authHeader string) (*Claims, error) {
	_ = ctx
	tokenString, err := extractBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return v.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidToken
	}
	if v.issuer != "" && claims.Issuer != v.issuer {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (v *TokenValidator) HasPermission(claims *Claims, permission string) bool {
	if claims == nil {
		return false
	}
	for _, role := range claims.Roles {
		if role == "admin" || role == "service" {
			return true
		}
	}
	for _, scope := range claims.Scopes {
		if scope == permission || scope == "dispatch:*" {
			return true
		}
	}
	return false
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", ErrMissingToken
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", ErrMissingToken
	}
	return parts[1], nil
}
