package service

import (
	"context"
	"testing"
	"time"

	"famgo/auth-service/internal/config"
	"famgo/auth-service/internal/model"
	"github.com/Abdex1/FamGo-platform/shared/pkg/logger"
)

// TestGenerateTokens tests JWT token generation
func TestGenerateTokens(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiry:     15,
		RefreshExpiry: 7,
	}

	log := logger.New("debug")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	tokens, err := service.GenerateTokens("user-123", "test@example.com", "rider")
	if err != nil {
		t.Fatalf("GenerateTokens failed: %v", err)
	}

	if tokens.AccessToken == "" {
		t.Error("AccessToken is empty")
	}

	if tokens.RefreshToken == "" {
		t.Error("RefreshToken is empty")
	}

	if tokens.TokenType != "Bearer" {
		t.Errorf("expected TokenType 'Bearer', got %s", tokens.TokenType)
	}
}

// TestVerifyToken tests JWT token verification
func TestVerifyToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiry:     15,
		RefreshExpiry: 7,
	}

	log := logger.New("debug")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	// Generate token
	tokens, err := service.GenerateTokens("user-123", "test@example.com", "rider")
	if err != nil {
		t.Fatalf("GenerateTokens failed: %v", err)
	}

	// Verify token
	claims, err := service.VerifyToken(tokens.AccessToken)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if claims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got %s", claims.UserID)
	}

	if claims.Email != "test@example.com" {
		t.Errorf("expected Email 'test@example.com', got %s", claims.Email)
	}

	if claims.Role != "rider" {
		t.Errorf("expected Role 'rider', got %s", claims.Role)
	}
}

// TestHashPassword tests password hashing with bcrypt
func TestHashPassword(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test"}
	log := logger.New("debug")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	password := "MySecurePassword123!"

	// Hash password
	hash, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("hash is empty")
	}

	if hash == password {
		t.Error("hash should not equal plain password")
	}
}

// TestGenerateOTP tests OTP generation
func TestGenerateOTP(t *testing.T) {
	cfg := &config.Config{JWTSecret: "test"}
	log := logger.New("debug")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	// Generate multiple OTPs to ensure randomness
	seen := make(map[string]bool)

	for i := 0; i < 10; i++ {
		otp, err := service.generateOTP()
		if err != nil {
			t.Fatalf("generateOTP failed: %v", err)
		}

		// Check OTP is 6 digits
		if len(otp) != 6 {
			t.Errorf("expected 6 digit OTP, got %s (length %d)", otp, len(otp))
		}

		// Check all characters are digits
		for _, ch := range otp {
			if ch < '0' || ch > '9' {
				t.Errorf("OTP contains non-digit: %c", ch)
			}
		}

		seen[otp] = true
	}

	// Basic check for randomness (not all should be same)
	if len(seen) == 1 {
		t.Error("OTP generation appears non-random (all same)")
	}
}

// TestRefreshToken tests token refresh
func TestRefreshToken(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiry:     15,
		RefreshExpiry: 7,
	}

	log := logger.New("debug")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	ctx := context.Background()

	// Generate initial tokens
	tokens, err := service.GenerateTokens("user-123", "test@example.com", "rider")
	if err != nil {
		t.Fatalf("GenerateTokens failed: %v", err)
	}

	// Wait briefly to ensure new timestamp
	time.Sleep(100 * time.Millisecond)

	// Refresh token
	newTokens, err := service.RefreshToken(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshToken failed: %v", err)
	}

	// Verify new access token has different signature (due to new IssuedAt)
	if newTokens.AccessToken == tokens.AccessToken {
		t.Error("new access token should be different from old one")
	}

	// Verify claims are preserved
	newClaims, err := service.VerifyToken(newTokens.AccessToken)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if newClaims.UserID != "user-123" {
		t.Errorf("expected UserID 'user-123', got %s", newClaims.UserID)
	}
}

// BenchmarkHashPassword benchmarks password hashing
func BenchmarkHashPassword(b *testing.B) {
	cfg := &config.Config{JWTSecret: "test"}
	log := logger.New("error")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	password := "MySecurePassword123!"

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		service.HashPassword(password)
	}
}

// BenchmarkGenerateTokens benchmarks token generation
func BenchmarkGenerateTokens(b *testing.B) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiry:     15,
		RefreshExpiry: 7,
	}

	log := logger.New("error")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		service.GenerateTokens("user-123", "test@example.com", "rider")
	}
}

// BenchmarkVerifyToken benchmarks token verification
func BenchmarkVerifyToken(b *testing.B) {
	cfg := &config.Config{
		JWTSecret:     "test-secret-key",
		JWTExpiry:     15,
		RefreshExpiry: 7,
	}

	log := logger.New("error")
	service := &AuthService{
		config: cfg,
		logger: log,
	}

	// Generate token once
	tokens, _ := service.GenerateTokens("user-123", "test@example.com", "rider")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		service.VerifyToken(tokens.AccessToken)
	}
}
