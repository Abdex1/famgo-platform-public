package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	domain "famgo/auth-service/internal/domain"
)

// AuthRepository manages auth data persistence
type AuthRepository struct {
	db *sqlx.DB
}

// NewAuthRepository creates repository
func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// User database model
type User struct {
	ID           string    `db:"id"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// CreateUser stores new user
func (r *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO auth.users (id, phone, email, password_hash, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Phone,
		user.Email,
		user.PasswordHash,
		user.FullName,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByPhone retrieves user by phone
func (r *AuthRepository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	user := &User{}
	query := `SELECT id, phone, email, password_hash, full_name, created_at, updated_at FROM auth.users WHERE phone = $1`

	err := r.db.GetContext(ctx, user, query, phone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves user by ID
func (r *AuthRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	user := &User{}
	query := `SELECT id, phone, email, password_hash, full_name, created_at, updated_at FROM auth.users WHERE id = $1`

	err := r.db.GetContext(ctx, user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Session database model
type Session struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	DeviceID    string    `db:"device_id"`
	Fingerprint string    `db:"fingerprint"`
	IPAddress   string    `db:"ip_address"`
	UserAgent   string    `db:"user_agent"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiresAt   time.Time `db:"expires_at"`
}

// CreateSession stores new session
func (r *AuthRepository) CreateSession(ctx context.Context, session *Session) error {
	query := `
		INSERT INTO auth.sessions (id, user_id, device_id, fingerprint, ip_address, user_agent, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		session.ID,
		session.UserID,
		session.DeviceID,
		session.Fingerprint,
		session.IPAddress,
		session.UserAgent,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// GetSessionByID retrieves active session
func (r *AuthRepository) GetSessionByID(ctx context.Context, id string) (*Session, error) {
	session := &Session{}
	query := `
		SELECT id, user_id, device_id, fingerprint, ip_address, user_agent, created_at, expires_at
		FROM auth.sessions
		WHERE id = $1 AND expires_at > NOW()
	`

	err := r.db.GetContext(ctx, session, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// InvalidateSession marks session as expired
func (r *AuthRepository) InvalidateSession(ctx context.Context, sessionID string) error {
	query := `UPDATE auth.sessions SET expires_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to invalidate session: %w", err)
	}

	return nil
}

// StoreOTP saves one-time password
func (r *AuthRepository) StoreOTP(ctx context.Context, userID string, otp string, expiresAt time.Time) error {
	query := `
		INSERT INTO auth.otp_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET token = $2, expires_at = $3
	`

	_, err := r.db.ExecContext(ctx, query, userID, otp, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	return nil
}

// VerifyOTP checks OTP validity
func (r *AuthRepository) VerifyOTP(ctx context.Context, userID string, otp string) (bool, error) {
	query := `
		SELECT token FROM auth.otp_tokens
		WHERE user_id = $1 AND token = $2 AND expires_at > NOW()
	`

	var token string
	err := r.db.GetContext(ctx, &token, query, userID, otp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify OTP: %w", err)
	}

	// Delete OTP after verification
	deleteQuery := `DELETE FROM auth.otp_tokens WHERE user_id = $1`
	r.db.ExecContext(ctx, deleteQuery, userID)

	return true, nil
}

// Device database model
type Device struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	Fingerprint string    `db:"fingerprint"`
	DeviceName  string    `db:"device_name"`
	DeviceType  string    `db:"device_type"`
	OSVersion   string    `db:"os_version"`
	AppVersion  string    `db:"app_version"`
	CreatedAt   time.Time `db:"created_at"`
}

// RegisterDevice stores device fingerprint
func (r *AuthRepository) RegisterDevice(ctx context.Context, device *Device) error {
	query := `
		INSERT INTO auth.devices (id, user_id, fingerprint, device_name, device_type, os_version, app_version, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		device.ID,
		device.UserID,
		device.Fingerprint,
		device.DeviceName,
		device.DeviceType,
		device.OSVersion,
		device.AppVersion,
		device.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to register device: %w", err)
	}

	return nil
}

// RBAC database model
type RBAC struct {
	ID          string                 `db:"id"`
	UserID      string                 `db:"user_id"`
	Roles       []string               `db:"roles"`
	Permissions map[string]interface{} `db:"permissions"`
	UpdatedAt   time.Time              `db:"updated_at"`
}

// StoreRBAC persists user roles and permissions
func (r *AuthRepository) StoreRBAC(ctx context.Context, userID string, roles []string, permissions map[string]interface{}) error {
	rolesJSON, _ := json.Marshal(roles)
	permsJSON, _ := json.Marshal(permissions)

	query := `
		INSERT INTO auth.rbac_policies (id, user_id, roles, permissions, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET roles = $3, permissions = $4, updated_at = $5
	`

	_, err := r.db.ExecContext(ctx, query,
		fmt.Sprintf("rbac-%s", userID),
		userID,
		string(rolesJSON),
		string(permsJSON),
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to store RBAC: %w", err)
	}

	return nil
}

// GetRBAC retrieves user roles and permissions
func (r *AuthRepository) GetRBAC(ctx context.Context, userID string) (*domain.RBACPolicy, error) {
	var roles string
	var permissions string

	query := `
		SELECT roles, permissions FROM auth.rbac_policies
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&roles, &permissions)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return default policy
			return &domain.RBACPolicy{
				UserID:      userID,
				Roles:       []string{"user"},
				Permissions: map[string]bool{"ride:create": true},
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, fmt.Errorf("failed to get RBAC: %w", err)
	}

	var rolesList []string
	var permsMap map[string]bool

	json.Unmarshal([]byte(roles), &rolesList)
	json.Unmarshal([]byte(permissions), &permsMap)

	return &domain.RBACPolicy{
		UserID:      userID,
		Roles:       rolesList,
		Permissions: permsMap,
		UpdatedAt:   time.Now(),
	}, nil
}
