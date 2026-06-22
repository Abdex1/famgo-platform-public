package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"famgo/auth-service/internal/model"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// DB returns the underlying database connection.
func (r *UserRepository) DB() *sqlx.DB {
	return r.db
}


// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users (id, email, password_hash, phone, role, status, email_verified, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, email, password_hash, phone, role, status, email_verified, created_at, updated_at
	`

	return r.db.QueryRowxContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Phone,
		user.Role,
		user.Status,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdateedAt,
	).StructScan(user)
}

// GetUserByEmail retrieves user by email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id, email, password_hash, phone, role, status, email_verified, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	user := &model.User{}
	if err := r.db.QueryRowxContext(ctx, query, email).StructScan(user); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves user by ID
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	query := `
	SELECT id, email, password_hash, phone, role, status, email_verified, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	user := &model.User{}
	if err := r.db.QueryRowxContext(ctx, query, userID).StructScan(user); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `
	UPDATE users
	SET password_hash = $1, phone = $2, status = $3, email_verified = $4, updated_at = $5
	WHERE id = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		user.PasswordHash,
		user.Phone,
		user.Status,
		user.EmailVerified,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

// OTPRepository handles database operations for OTP
type OTPRepository struct {
	db *sqlx.DB
}

// NewOTPRepository creates a new OTP repository
func NewOTPRepository(db *sqlx.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

// SaveOTP saves an OTP to the database
func (r *OTPRepository) SaveOTP(ctx context.Context, otp *model.OTPVerification) error {
	query := `
	INSERT INTO otp_verification (id, email, otp, expires_at, attempts, verified, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		otp.ID,
		otp.Email,
		otp.OTP,
		otp.ExpiresAt,
		otp.Attempts,
		otp.Verified,
		otp.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	return nil
}

// GetOTPByEmailAndOTP retrieves OTP by email and OTP code
func (r *OTPRepository) GetOTPByEmailAndOTP(ctx context.Context, email, otp string) (*model.OTPVerification, error) {
	query := `
	SELECT id, email, otp, expires_at, attempts, verified, created_at
	FROM otp_verification
	WHERE email = $1 AND otp = $2 AND verified = FALSE
	ORDER BY created_at DESC
	LIMIT 1
	`

	otpRec := &model.OTPVerification{}
	if err := r.db.QueryRowxContext(ctx, query, email, otp).StructScan(otpRec); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New("OTP not found or already verified")
		}
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// Check if OTP is expired
	if otpRec.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("OTP expired")
	}

	return otpRec, nil
}

// MarkOTPAsVerified marks OTP as verified
func (r *OTPRepository) MarkOTPAsVerified(ctx context.Context, otpID string) error {
	query := `
	UPDATE otp_verification
	SET verified = TRUE
	WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, otpID)
	if err != nil {
		return fmt.Errorf("failed to mark OTP as verified: %w", err)
	}

	return nil
}

// IncrementOTPAttempts increments OTP attempt counter
func (r *OTPRepository) IncrementOTPAttempts(ctx context.Context, otpID string) error {
	query := `
	UPDATE otp_verification
	SET attempts = attempts + 1
	WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, otpID)
	if err != nil {
		return fmt.Errorf("failed to increment OTP attempts: %w", err)
	}

	return nil
}

// DeleteExpiredOTPs deletes all expired OTPs (cleanup job)
func (r *OTPRepository) DeleteExpiredOTPs(ctx context.Context) error {
	query := `
	DELETE FROM otp_verification
	WHERE expires_at < NOW() AND verified = TRUE
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired OTPs: %w", err)
	}

	return nil
}
