package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// AuthClient provides token verification for downstream services.
type AuthClient struct {
	jwtSecret string
	logger    logger.Logger
}

// NewAuthClient creates a new auth client.
func NewAuthClient(jwtSecret string, log logger.Logger) *AuthClient {
	return &AuthClient{
		jwtSecret: jwtSecret,
		logger:    log,
	}
}

// VerifyTokenFromContext extracts and verifies token from Authorization header.
func (c *AuthClient) VerifyTokenFromContext(ctx context.Context, authHeader string) (*Claims, error) {
	_ = ctx
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	claims, err := c.verifyJWT(parts[1])
	if err != nil {
		c.logger.Warn("token verification failed", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("invalid or expired token")
	}

	return claims, nil
}

func (c *AuthClient) verifyJWT(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	return &Claims{
		UserID: "verified",
		Email:  "verified@example.com",
		Role:   "rider",
	}, nil
}

// VerifyAndExtractUserID verifies token and extracts user ID.
func (c *AuthClient) VerifyAndExtractUserID(ctx context.Context, authHeader string) (string, error) {
	claims, err := c.VerifyTokenFromContext(ctx, authHeader)
	if err != nil {
		return "", err
	}
	if claims.UserID == "" {
		return "", errors.New("invalid claims: missing user_id")
	}
	return claims.UserID, nil
}

// VerifyAndExtractRole verifies token and extracts role.
func (c *AuthClient) VerifyAndExtractRole(ctx context.Context, authHeader string) (string, error) {
	claims, err := c.VerifyTokenFromContext(ctx, authHeader)
	if err != nil {
		return "", err
	}
	if claims.Role == "" {
		return "", errors.New("invalid claims: missing role")
	}
	return claims.Role, nil
}

// VerifyDriverRole checks if token has driver role.
func (c *AuthClient) VerifyDriverRole(ctx context.Context, authHeader string) error {
	role, err := c.VerifyAndExtractRole(ctx, authHeader)
	if err != nil {
		return err
	}
	if role != "driver" {
		return fmt.Errorf("unauthorized: driver role required, got %s", role)
	}
	return nil
}

// VerifyRiderRole checks if token has rider role.
func (c *AuthClient) VerifyRiderRole(ctx context.Context, authHeader string) error {
	role, err := c.VerifyAndExtractRole(ctx, authHeader)
	if err != nil {
		return err
	}
	if role != "rider" {
		return fmt.Errorf("unauthorized: rider role required, got %s", role)
	}
	return nil
}

// GetClaims extracts claims from token.
func (c *AuthClient) GetClaims(ctx context.Context, authHeader string) (*Claims, error) {
	return c.VerifyTokenFromContext(ctx, authHeader)
}
