package domain

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService defines authentication business logic
type AuthService interface {
	// Login authenticates a user and returns tokens
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)

	// Register creates a new user account
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)

	// VerifyToken validates and parses a JWT token
	VerifyToken(ctx context.Context, tokenString string) (*TokenClaims, error)

	// RefreshToken generates new access token from refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error)

	// RegisterDevice associates device fingerprint with user
	RegisterDevice(ctx context.Context, userID string, req *DeviceRegistrationRequest) error

	// VerifyOTP validates one-time password
	VerifyOTP(ctx context.Context, userID string, otp string) (bool, error)

	// GenerateOTP creates a new one-time password
	GenerateOTP(ctx context.Context, userID string) (string, error)

	// EnableMFA activates multi-factor authentication
	EnableMFA(ctx context.Context, userID string) (*MFAResponse, error)

	// VerifyMFA validates MFA token
	VerifyMFA(ctx context.Context, userID string, token string) (bool, error)

	// GetRBAC retrieves role-based access control for user
	GetRBAC(ctx context.Context, userID string) (*RBACPolicy, error)

	// CreateSession starts user session
	CreateSession(ctx context.Context, userID string, req *SessionRequest) (*SessionResponse, error)

	// ValidateSession checks if session is active
	ValidateSession(ctx context.Context, sessionID string) (bool, error)

	// Logout terminates user session
	Logout(ctx context.Context, sessionID string) error

	// RevokeAllSessions invalidates all user sessions
	RevokeAllSessions(ctx context.Context, userID string) error
}

// LoginRequest authentication credentials
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"device_id"`
}

// RegisterRequest new user registration
type RegisterRequest struct {
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"email"`
	Password    string `json:"password" binding:"min=8"`
	FullName    string `json:"full_name"`
	DeviceID    string `json:"device_id"`
	Fingerprint string `json:"fingerprint"`
}

// AuthResponse successful authentication
type AuthResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
	User         *UserProfile  `json:"user"`
	Session      *SessionData  `json:"session"`
}

// TokenResponse token refresh result
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// TokenClaims JWT claims structure
type TokenClaims struct {
	UserID      string
	Phone       string
	Email       string
	Roles       []string
	MFAVerified bool
	jwt.RegisteredClaims
}

// UserProfile authenticated user info
type UserProfile struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SessionRequest session creation parameters
type SessionRequest struct {
	DeviceID    string `json:"device_id"`
	DeviceName  string `json:"device_name"`
	IPAddress   string `json:"ip_address"`
	Fingerprint string `json:"fingerprint"`
	UserAgent   string `json:"user_agent"`
}

// SessionResponse active session info
type SessionResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	DeviceID  string    `json:"device_id"`
}

// SessionData session details
type SessionData struct {
	ID       string `json:"id"`
	DeviceID string `json:"device_id"`
}

// DeviceRegistrationRequest device association
type DeviceRegistrationRequest struct {
	DeviceID    string `json:"device_id" binding:"required"`
	DeviceName  string `json:"device_name"`
	DeviceType  string `json:"device_type"` // ios, android, web
	Fingerprint string `json:"fingerprint"`
	OSVersion   string `json:"os_version"`
	AppVersion  string `json:"app_version"`
}

// MFAResponse multi-factor auth setup
type MFAResponse struct {
	Secret     string   `json:"secret"`
	QRCode     string   `json:"qr_code"`
	BackupCodes []string `json:"backup_codes"`
}

// RBACPolicy role-based access control
type RBACPolicy struct {
	UserID string
	Roles  []string
	Permissions map[string]bool
	UpdatedAt time.Time
}

// Errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidOTP         = errors.New("invalid OTP")
	ErrOTPExpired         = errors.New("OTP expired")
	ErrMFARequired        = errors.New("MFA verification required")
	ErrDeviceNotTrusted   = errors.New("device not trusted")
	ErrSessionExpired     = errors.New("session expired")
	ErrSessionNotFound    = errors.New("session not found")
)
