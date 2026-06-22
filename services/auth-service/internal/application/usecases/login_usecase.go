
/*
# PHASE 1 — ENTERPRISE LOGIN FLOW
# =========================================================

# STEP 1 — CREATE AUTH USECASE

services/auth-service/internal/application/usecases/login_usecase.go

go get github.com/google/uuid
*/
package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"famgo/auth-service/internal/domain/entities"
	"famgo/auth-service/internal/infrastructure/security"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	StoreRefreshToken(
		ctx context.Context,
		token *entities.RefreshTokenFamily,
	) error
}

type SessionStore interface {
	StoreSession(
		ctx context.Context,
		sessionID string,
		userID string,
		ttl time.Duration,
	) error
}

type LoginUseCase struct {
	userRepo     UserRepository
	jwtManager   *security.JWTManager
	sessionStore SessionStore
}
type EventPublisher interface {
	Publish(
		ctx context.Context,
		topic string,
		key string,
		payload any,
	) error
}
func NewLoginUseCase(
	userRepo UserRepository,
	jwtManager *security.JWTManager,
	sessionStore SessionStore,
) *LoginUseCase {

	return &LoginUseCase{
		userRepo:     userRepo,
		jwtManager:   jwtManager,
		sessionStore: sessionStore,
	}
}

func (u *LoginUseCase) Execute(
	ctx context.Context,
	email string,
	password string,
	deviceID string,
	ip string,
	userAgent string,
) (string, string, error) {

	user, err := u.userRepo.GetByEmail(ctx, email)

	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := u.jwtManager.GenerateAccessToken(
		user.ID,
		string(user.Role),
	)

	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.NewString()

	familyID := uuid.NewString()

	refreshEntity := &entities.RefreshTokenFamily{
		ID:         uuid.NewString(),
		UserID:     user.ID,
		FamilyID:   familyID,
		TokenHash:  security.HashToken(refreshToken),

		DeviceID:   deviceID,
		IPAddress:  ip,
		UserAgent:  userAgent,

		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	}

	err = u.userRepo.StoreRefreshToken(
		ctx,
		refreshEntity,
	)

	if err != nil {
		return "", "", err
	}

	err = u.sessionStore.StoreSession(
		ctx,
		familyID,
		user.ID,
		30*24*time.Hour,
	)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
