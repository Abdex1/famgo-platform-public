package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents an authenticated user
type User struct {
	ID            string    `db:"id"`
	Email         string    `db:"email"`
	PasswordHash  string    `db:"password_hash"`
	Phone         string    `db:"phone"`
	Role          string    `db:"role"`        // "rider" or "driver"
	Status        string    `db:"status"`      // "pending", "active", "suspended"
	EmailVerified bool      `db:"email_verified"`
	CreatedAt     time.Time `db:"created_at"`
	UpdateedAt    time.Time `db:"updated_at"`
}

// OTPVerification represents an OTP record
type OTPVerification struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	OTP       string    `db:"otp"`
	ExpiresAt time.Time `db:"expires_at"`
	Attempts  int       `db:"attempts"`
	Verified  bool      `db:"verified"`
	CreatedAt time.Time `db:"created_at"`
}

// TokenResponse represents JWT tokens returned to client
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Claims represents JWT claims (Pattern from Uber)
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// RegistrationRequest represents first step of registration (send OTP)
type RegistrationRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=rider driver"`
	Name     string `json:"name" validate:"required"`
}

// VerifyRegistrationRequest represents second step (verify OTP and create account)
type VerifyRegistrationRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordResetVerifyRequest represents verification and new password
type PasswordResetVerifyRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTP         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// SuccessResponse represents success response structure
type SuccessResponse struct {
	Code   string      `json:"code"`
	Data   interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
