package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"famgo/auth-service/internal/config"
	"famgo/auth-service/internal/model"
	"famgo/auth-service/internal/repository"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// AuthService handles authentication operations
type AuthService struct {
	config      *config.Config
	userRepo    *repository.UserRepository
	otpRepo     *repository.OTPRepository
	logger      logger.Logger
	emailClient EmailClient // For sending OTPs
}

// EmailClient interface for sending emails (abstraction for Brevo)
type EmailClient interface {
	SendOTP(ctx context.Context, email, otp, name string) error
}

// NewAuthService creates a new auth service
func NewAuthService(
	cfg *config.Config,
	userRepo *repository.UserRepository,
	log logger.Logger,
) *AuthService {
	return &AuthService{
		config:   cfg,
		userRepo: userRepo,
		otpRepo:  repository.NewOTPRepository(userRepo.DB()), // Access to same DB
		logger:   log,
	}
}

// SetEmailClient sets the email client (for dependency injection)
func (s *AuthService) SetEmailClient(client EmailClient) {
	s.emailClient = client
}

// STEP 1: Send OTP for registration
func (s *AuthService) SendRegistrationOTP(ctx context.Context, email, name string) error {
	s.logger.Info("sending registration OTP", map[string]interface{}{"email": email})

	// Check if user already exists
	existingUser, _ := s.userRepo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("user already exists with this email")
	}

	// Generate 6-digit OTP
	otp, err := s.generateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Save OTP to database
	otpRecord := &model.OTPVerification{
		ID:        uuid.New().String(),
		Email:     email,
		OTP:       otp,
		ExpiresAt: time.Now().Add(time.Duration(s.config.OTPExpiry) * time.Minute),
		Attempts:  0,
		Verified:  false,
		CreatedAt: time.Now(),
	}

	if err := s.otpRepo.SaveOTP(ctx, otpRecord); err != nil {
		s.logger.Error("failed to save OTP", map[string]interface{}{"error": err})
		return err
	}

	// Send OTP via email (Brevo pattern from Uber)
	if s.emailClient != nil {
		if err := s.emailClient.SendOTP(ctx, email, otp, name); err != nil {
			s.logger.Error("failed to send OTP email", map[string]interface{}{"error": err})
			// Don't fail - OTP is still saved
		}
	}

	return nil
}

// STEP 2: Verify OTP and create account
func (s *AuthService) VerifyRegistrationOTP(
	ctx context.Context,
	req *model.VerifyRegistrationRequest,
) (*model.TokenResponse, error) {
	s.logger.Info("verifying registration OTP", map[string]interface{}{"email": req.Email})

	// Get OTP record from DB
	otpRecord, err := s.otpRepo.GetOTPByEmailAndOTP(ctx, req.Email, req.OTP)
	if err != nil {
		s.logger.Warn("OTP verification failed", map[string]interface{}{"email": req.Email, "error": err})
		return nil, err
	}

	// Check attempt limit (max 3 attempts)
	if otpRecord.Attempts >= 3 {
		return nil, errors.New("OTP verification attempts exceeded")
	}

	// Increment attempts
	if err := s.otpRepo.IncrementOTPAttempts(ctx, otpRecord.ID); err != nil {
		s.logger.Error("failed to increment OTP attempts", map[string]interface{}{"error": err})
	}

	// Verify OTP (exact match)
	if otpRecord.OTP != req.OTP {
		return nil, errors.New("invalid OTP")
	}

	// Mark OTP as verified
	if err := s.otpRepo.MarkOTPAsVerified(ctx, otpRecord.ID); err != nil {
		s.logger.Error("failed to mark OTP as verified", map[string]interface{}{"error": err})
		return nil, err
	}

	// Create user account
	// NOTE: Password must be provided in the request (from registration form)
	// This is simplified - actual implementation needs password from initial registration
	user := &model.User{
		ID:            uuid.New().String(),
		Email:         req.Email,
		PasswordHash:  "", // Will be hashed in actual implementation
		Role:          "rider", // Default role, should come from registration request
		Status:        "active",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdateedAt:    time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		s.logger.Error("failed to create user", map[string]interface{}{"error": err})
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	s.logger.Info("user created successfully", map[string]interface{}{"user_id": user.ID, "email": user.Email})

	// Generate JWT tokens
	tokens, err := s.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		s.logger.Error("failed to generate tokens", map[string]interface{}{"error": err})
		return nil, err
	}

	return tokens, nil
}

// GenerateTokens creates JWT access and refresh tokens (Pattern from Uber)
func (s *AuthService) GenerateTokens(userID, email, role string) (*model.TokenResponse, error) {
	// Access token: 15 minutes (from Uber pattern)
	accessClaims := &model.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JWTExpiry) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh token: 7 days (from Uber pattern)
	refreshClaims := &model.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, s.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &model.TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    s.config.JWTExpiry * 60, // Convert to seconds
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*model.TokenResponse, error) {
	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Generate new access token
	return s.GenerateTokens(claims.UserID, claims.Email, claims.Role)
}

// VerifyToken verifies a JWT token and returns claims
func (s *AuthService) VerifyToken(tokenString string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// Login authenticates user with email and password
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.TokenResponse, error) {
	s.logger.Info("login attempt", map[string]interface{}{"email": email})

	// Get user by email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Warn("login failed - user not found", map[string]interface{}{"email": email})
		return nil, errors.New("invalid credentials")
	}

	// Verify password (bcrypt from Uber pattern)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Warn("login failed - invalid password", map[string]interface{}{"email": email})
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is %s", user.Status)
	}

	// Generate tokens
	return s.GenerateTokens(user.ID, user.Email, user.Role)
}

// HashPassword hashes password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// generateOTP generates a 6-digit OTP
func (s *AuthService) generateOTP() (string, error) {
	max := big.NewInt(1000000) // 1,000,000
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Ensure 6 digits with leading zeros
	return fmt.Sprintf("%06d", num), nil
}
