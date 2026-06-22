/*
# PHASE 1 — REFRESH TOKEN ROTATION ENGINE

# =========================================================

# STEP 2 — CREATE ROTATION USECASE

internal/application/usecases/refresh_usecase.go
*/
package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"famgo/auth-service/internal/infrastructure/security"
)

type RefreshRepository interface {
	FindRefreshToken(
		ctx context.Context,
		hash string,
	) error

	RevokeFamily(
		ctx context.Context,
		familyID string,
	) error
}

type RefreshUseCase struct {
	repo       RefreshRepository
	jwtManager *security.JWTManager
}

func NewRefreshUseCase(
	repo RefreshRepository,
	jwtManager *security.JWTManager,
) *RefreshUseCase {

	return &RefreshUseCase{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

func (u *RefreshUseCase) Execute(
	ctx context.Context,
	refreshToken string,
	userID string,
	role string,
) (string, string, error) {

	hash := security.HashToken(refreshToken)

	err := u.repo.FindRefreshToken(ctx, hash)

	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(
		userID,
		role,
	)

	if err != nil {
		return "", "", err
	}

	newRefresh := uuid.NewString()

	return accessToken, newRefresh, nil
}
