//go:build ignore

package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// Global validator instance
var validate = validator.New()

// ============================================================================
// REQUEST MODELS WITH VALIDATION TAGS
// ============================================================================

// SignupRequest represents user registration request
type SignupRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=128,password"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" validate:"required,e164"`
	Role      string `json:"role" validate:"omitempty,oneof=passenger driver"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// PasswordResetRequest represents password reset request
type PasswordResetRequest struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=128,password"`
	OTPCode     string `json:"otp_code" validate:"required,len=6,numeric"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// ============================================================================
// CUSTOM VALIDATORS
// ============================================================================

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length (already handled by min tag, but explicit here)
	if len(password) < 8 || len(password) > 128 {
		return false
	}

	// Check for uppercase letters
	hasUpper := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	// Check for lowercase letters
	hasLower := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	// Check for numbers
	hasNumber := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return false
	}

	// Check for special characters
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	hasSpecial := false
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return false
	}

	return true
}

// validateE164Phone validates E.164 phone format
func validateE164Phone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// E.164 format: +[country code][number]
	// Example: +251911234567
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return e164Regex.MatchString(phone)
}

// ============================================================================
// INITIALIZATION
// ============================================================================

func init() {
	// Register custom validators
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("e164", validateE164Phone)
}

// ============================================================================
// PUBLIC VALIDATION FUNCTIONS
// ============================================================================

// ValidateSignup validates signup request
func ValidateSignup(req SignupRequest) error {
	if err := validate.Struct(req); err != nil {
		return formatValidationError(err)
	}

	// Additional business logic validation
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}

// ValidateLogin validates login request
func ValidateLogin(req LoginRequest) error {
	if err := validate.Struct(req); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// ValidatePasswordReset validates password reset request
func ValidatePasswordReset(req PasswordResetRequest) error {
	if err := validate.Struct(req); err != nil {
		return formatValidationError(err)
	}

	// Verify OTP is numeric
	if !regexp.MustCompile(`^\d{6}$`).MatchString(req.OTPCode) {
		return fmt.Errorf("OTP must be 6 digits")
	}

	return nil
}

// ValidateUpdateProfile validates profile update request
func ValidateUpdateProfile(req UpdateProfileRequest) error {
	if err := validate.Struct(req); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// ValidateChangePassword validates password change request
func ValidateChangePassword(req ChangePasswordRequest) error {
	if err := validate.Struct(req); err != nil {
		return formatValidationError(err)
	}

	// Verify passwords match
	if req.NewPassword != req.ConfirmPassword {
		return fmt.Errorf("new password and confirmation don't match")
	}

	// Verify new password is different from current
	if req.CurrentPassword == req.NewPassword {
		return fmt.Errorf("new password must be different from current password")
	}

	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// formatValidationError formats validation errors for user-friendly output
func formatValidationError(err error) error {
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	var messages []string

	for _, fieldError := range validationErrors {
		messages = append(messages, formatFieldError(fieldError))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
}

// formatFieldError formats a single field validation error
func formatFieldError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", err.Field(), err.Param())
	case "eqfield":
		return fmt.Sprintf("%s must match %s", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", err.Field())
	case "password":
		return fmt.Sprintf("%s must contain uppercase, lowercase, numbers, and special characters", err.Field())
	case "e164":
		return fmt.Sprintf("%s must be in E.164 format (e.g., +251911234567)", err.Field())
	default:
		return fmt.Sprintf("%s failed validation: %s", err.Field(), err.Tag())
	}
}

// ValidateEmail validates email format more strictly
func ValidateEmail(email string) error {
	if len(email) == 0 {
		return fmt.Errorf("email is required")
	}

	if len(email) > 255 {
		return fmt.Errorf("email is too long")
	}

	// Use validator
	return validate.Var(email, "email")
}

// ValidatePhone validates phone format
func ValidatePhone(phone string) error {
	if len(phone) == 0 {
		return fmt.Errorf("phone is required")
	}

	return validate.Var(phone, "e164")
}

// ValidatePasswordStrength validates password meets all requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password is too long")
	}

	if err := validate.Var(password, "password"); err != nil {
		// Custom logic for better error messages
		if !hasUppercase(password) {
			return fmt.Errorf("password must contain at least one uppercase letter")
		}
		if !hasLowercase(password) {
			return fmt.Errorf("password must contain at least one lowercase letter")
		}
		if !hasDigit(password) {
			return fmt.Errorf("password must contain at least one number")
		}
		if !hasSpecialChar(password) {
			return fmt.Errorf("password must contain at least one special character (!@#$%^&*)")
		}
		return nil
	}

	return nil
}

// Helper functions for password validation
func hasUppercase(s string) bool {
	for _, c := range s {
		if unicode.IsUpper(c) {
			return true
		}
	}
	return false
}

func hasLowercase(s string) bool {
	for _, c := range s {
		if unicode.IsLower(c) {
			return true
		}
	}
	return false
}

func hasDigit(s string) bool {
	for _, c := range s {
		if unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func hasSpecialChar(s string) bool {
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
	for _, c := range s {
		if strings.ContainsRune(specialChars, c) {
			return true
		}
	}
	return false
}
