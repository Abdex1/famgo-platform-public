//go:build ignore

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// SIGNUP VALIDATION TESTS
// ============================================================================

func TestValidateSignup_Success(t *testing.T) {
	req := SignupRequest{
		Email:     "abdu.aliyi@example.com",
		Password:  "SecurePass123!@#",
		FirstName: "abdu",
		LastName:  "aliyi",
		Phone:     "+251911234567",
		Role:      "passenger",
	}

	err := ValidateSignup(req)
	assert.NoError(t, err, "valid signup request should not error")
}

func TestValidateSignup_MissingEmail(t *testing.T) {
	req := SignupRequest{
		Email:     "",
		Password:  "SecurePass123!@#",
		FirstName: "abdu",
		LastName:  "aliyi",
		Phone:     "+251911234567",
	}

	err := ValidateSignup(req)
	assert.Error(t, err, "missing email should error")
	assert.Contains(t, err.Error(), "email", "error should mention email")
}

func TestValidateSignup_InvalidEmail(t *testing.T) {
	req := SignupRequest{
		Email:     "invalid-email",
		Password:  "SecurePass123!@#",
		FirstName: "abdu",
		LastName:  "aliyi",
		Phone:     "+251911234567",
	}

	err := ValidateSignup(req)
	assert.Error(t, err, "invalid email should error")
}

func TestValidateSignup_WeakPassword(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		reason   string
	}{
		{"TooShort", "Short1!", "too short"},
		{"NoUppercase", "lowercase123!@#", "missing uppercase"},
		{"NoLowercase", "UPPERCASE123!@#", "missing lowercase"},
		{"NoNumbers", "NoNumbers!@#Abc", "missing numbers"},
		{"NoSpecialChars", "NoSpecialChars123", "missing special characters"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := SignupRequest{
				Email:     "abdu@example.com",
				Password:  tc.password,
				FirstName: "abdu",
				LastName:  "aliyi",
				Phone:     "+251911234567",
			}

			err := ValidateSignup(req)
			assert.Error(t, err, "weak password should error: "+tc.reason)
		})
	}
}

func TestValidateSignup_InvalidPhone(t *testing.T) {
	req := SignupRequest{
		Email:     "abdu@example.com",
		Password:  "SecurePass123!@#",
		FirstName: "abdu",
		LastName:  "aliyi",
		Phone:     "123456789", // Missing country code
	}

	err := ValidateSignup(req)
	assert.Error(t, err, "invalid phone should error")
}

func TestValidateSignup_ShortName(t *testing.T) {
	req := SignupRequest{
		Email:     "abdu@example.com",
		Password:  "SecurePass123!@#",
		FirstName: "J", // Too short
		LastName:  "aliyi",
		Phone:     "+251911234567",
	}

	err := ValidateSignup(req)
	assert.Error(t, err, "short first name should error")
}

// ============================================================================
// LOGIN VALIDATION TESTS
// ============================================================================

func TestValidateLogin_Success(t *testing.T) {
	req := LoginRequest{
		Email:    "abdu@example.com",
		Password: "SecurePass123!@#",
	}

	err := ValidateLogin(req)
	assert.NoError(t, err, "valid login request should not error")
}

func TestValidateLogin_MissingEmail(t *testing.T) {
	req := LoginRequest{
		Email:    "",
		Password: "SecurePass123!@#",
	}

	err := ValidateLogin(req)
	assert.Error(t, err, "missing email should error")
}

func TestValidateLogin_MissingPassword(t *testing.T) {
	req := LoginRequest{
		Email:    "abdu@example.com",
		Password: "",
	}

	err := ValidateLogin(req)
	assert.Error(t, err, "missing password should error")
}

func TestValidateLogin_InvalidEmail(t *testing.T) {
	req := LoginRequest{
		Email:    "invalid-email",
		Password: "SecurePass123!@#",
	}

	err := ValidateLogin(req)
	assert.Error(t, err, "invalid email should error")
}

// ============================================================================
// PASSWORD RESET VALIDATION TESTS
// ============================================================================

func TestValidatePasswordReset_Success(t *testing.T) {
	req := PasswordResetRequest{
		Email:       "abdu@example.com",
		NewPassword: "NewSecure123!@#",
		OTPCode:     "123456",
	}

	err := ValidatePasswordReset(req)
	assert.NoError(t, err, "valid password reset should not error")
}

func TestValidatePasswordReset_InvalidOTP(t *testing.T) {
	testCases := []struct {
		name    string
		otpCode string
	}{
		{"TooShort", "12345"},
		{"TooLong", "1234567"},
		{"NotNumeric", "12345a"},
		{"Empty", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := PasswordResetRequest{
				Email:       "abdu@example.com",
				NewPassword: "NewSecure123!@#",
				OTPCode:     tc.otpCode,
			}

			err := ValidatePasswordReset(req)
			assert.Error(t, err, "invalid OTP should error: "+tc.name)
		})
	}
}

// ============================================================================
// HELPER FUNCTION TESTS
// ============================================================================

func TestValidateEmail_Success(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"abdu.aliyi@famgo.co.uk",
		"test+tag@domain.com",
	}

	for _, email := range validEmails {
		err := ValidateEmail(email)
		assert.NoError(t, err, "valid email should not error: "+email)
	}
}

func TestValidateEmail_Fail(t *testing.T) {
	invalidEmails := []string{
		"",
		"invalid",
		"invalid@",
		"invalid@.",
		"@invalid.com",
	}

	for _, email := range invalidEmails {
		err := ValidateEmail(email)
		assert.Error(t, err, "invalid email should error: "+email)
	}
}

func TestValidatePhone_Success(t *testing.T) {
	validPhones := []string{
		"+251911234567",
		"+1234567890",
		"+44201234567",
	}

	for _, phone := range validPhones {
		err := ValidatePhone(phone)
		assert.NoError(t, err, "valid phone should not error: "+phone)
	}
}

func TestValidatePhone_Fail(t *testing.T) {
	invalidPhones := []string{
		"",
		"1234567890",
		"+",
		"abc123456",
	}

	for _, phone := range invalidPhones {
		err := ValidatePhone(phone)
		assert.Error(t, err, "invalid phone should error: "+phone)
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		valid    bool
	}{
		{"ValidStrong", "MyStrong123!@#", true},
		{"ValidStrong2", "P@ssw0rd123", true},
		{"TooShort", "Short1!", false},
		{"NoUpper", "lowercase123!@#", false},
		{"NoLower", "UPPERCASE123!@#", false},
		{"NoNumbers", "NoNumbers!@#Abc", false},
		{"NoSpecial", "NoSpecial123Abc", false},
		{"TooLong", "A" + string(make([]rune, 129)) + "!", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tc.password)
			if tc.valid {
				assert.NoError(t, err, "should be valid: "+tc.name)
			} else {
				assert.Error(t, err, "should be invalid: "+tc.name)
			}
		})
	}
}

func TestUpdateProfileRequest_Valid(t *testing.T) {
	req := UpdateProfileRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Phone:     "+251911234567",
	}

	err := ValidateUpdateProfile(req)
	assert.NoError(t, err, "valid profile update should not error")
}

func TestChangePasswordRequest_Success(t *testing.T) {
	req := ChangePasswordRequest{
		CurrentPassword: "OldPassword123!@#",
		NewPassword:     "NewPassword456!@#",
		ConfirmPassword: "NewPassword456!@#",
	}

	err := ValidateChangePassword(req)
	assert.NoError(t, err, "valid password change should not error")
}

func TestChangePasswordRequest_MismatchedPasswords(t *testing.T) {
	req := ChangePasswordRequest{
		CurrentPassword: "OldPassword123!@#",
		NewPassword:     "NewPassword456!@#",
		ConfirmPassword: "DifferentPassword!@#",
	}

	err := ValidateChangePassword(req)
	assert.Error(t, err, "mismatched passwords should error")
	assert.Contains(t, err.Error(), "don't match")
}

func TestChangePasswordRequest_SameAsOld(t *testing.T) {
	req := ChangePasswordRequest{
		CurrentPassword: "OldPassword123!@#",
		NewPassword:     "OldPassword123!@#",
		ConfirmPassword: "OldPassword123!@#",
	}

	err := ValidateChangePassword(req)
	assert.Error(t, err, "same password should error")
	assert.Contains(t, err.Error(), "different")
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkValidateSignup(b *testing.B) {
	req := SignupRequest{
		Email:     "abdu.aliyi@example.com",
		Password:  "SecurePass123!@#",
		FirstName: "abdu",
		LastName:  "aliyi",
		Phone:     "+251911234567",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateSignup(req)
	}
}

func BenchmarkValidateLogin(b *testing.B) {
	req := LoginRequest{
		Email:    "abdu@example.com",
		Password: "SecurePass123!@#",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateLogin(req)
	}
}

// ============================================================================
// TABLE-DRIVEN TESTS
// ============================================================================

type ValidateTestCase struct {
	Name    string
	Request SignupRequest
	Valid   bool
	ErrMsg  string
}

func TestValidateSignup_TableDriven(t *testing.T) {
	cases := []ValidateTestCase{
		{
			Name: "Valid",
			Request: SignupRequest{
				Email:     "test@example.com",
				Password:  "Test1234!@#",
				FirstName: "Test",
				LastName:  "User",
				Phone:     "+251911234567",
			},
			Valid: true,
		},
		{
			Name: "InvalidEmail",
			Request: SignupRequest{
				Email:     "invalid",
				Password:  "Test1234!@#",
				FirstName: "Test",
				LastName:  "User",
				Phone:     "+251911234567",
			},
			Valid:  false,
			ErrMsg: "email",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := ValidateSignup(tc.Request)
			if tc.Valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tc.ErrMsg != "" {
					assert.Contains(t, err.Error(), tc.ErrMsg)
				}
			}
		})
	}
}
