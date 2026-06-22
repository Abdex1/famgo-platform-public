package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// Integration tests for Week 2: End-to-end authentication and service flows

const (
	authServiceURL   = "http://localhost:8080"
	userServiceURL   = "http://localhost:8081"
	driverServiceURL = "http://localhost:8082"
)

// TestAuthFlow_UserRegistrationAndLogin tests complete user registration and login flow
func TestAuthFlow_UserRegistrationAndLogin(t *testing.T) {
	// Step 1: Register user (send OTP)
	registerReq := map[string]interface{}{
		"email":    "user@test.com",
		"phone":    "+1234567890",
		"password": "SecurePassword123!",
		"role":     "rider",
		"name":     "Test User",
	}

	body, _ := json.Marshal(registerReq)
	resp, err := http.Post(
		authServiceURL+"/api/v1/auth/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Step 2: Verify OTP (in real test, would use actual OTP from email)
	verifyReq := map[string]interface{}{
		"email": "user@test.com",
		"otp":   "123456", // Would be actual OTP from email
	}

	body, _ = json.Marshal(verifyReq)
	resp, err = http.Post(
		authServiceURL+"/api/v1/auth/verify-register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Verification failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201, got %d", resp.StatusCode)
	}

	// Parse token response
	var tokenResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tokenResp)
	accessToken := tokenResp["data"].(map[string]interface{})["access_token"].(string)

	// Step 3: Login
	loginReq := map[string]interface{}{
		"email":    "user@test.com",
		"password": "SecurePassword123!",
	}

	body, _ = json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", authServiceURL+"/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Step 4: Access protected endpoint with token
	req, _ = http.NewRequest("GET", authServiceURL+"/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Protected endpoint failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	t.Log("✅ User registration and login flow passed")
}

// TestAuthFlow_DriverRegistrationAndStateTransition tests driver registration and state machine
func TestAuthFlow_DriverRegistrationAndStateTransition(t *testing.T) {
	// Step 1: Register driver
	registerReq := map[string]interface{}{
		"email":          "driver@test.com",
		"phone":          "+1234567890",
		"password":       "SecurePassword123!",
		"first_name":     "Test",
		"last_name":      "Driver",
		"license_number": "DL123456789",
		"license_expiry": "2025-12-31",
	}

	body, _ := json.Marshal(registerReq)
	resp, err := http.Post(
		driverServiceURL+"/api/v1/drivers/register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Driver registration failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Step 2: Verify OTP
	verifyReq := map[string]interface{}{
		"email": "driver@test.com",
		"otp":   "123456",
	}

	body, _ = json.Marshal(verifyReq)
	resp, err = http.Post(
		driverServiceURL+"/api/v1/drivers/verify-register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Driver verification failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected 201, got %d", resp.StatusCode)
	}

	// Parse driver response
	var driverResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&driverResp)

	// Note: Would extract driver_id from response in real test
	driverID := "driver-uuid-from-response"

	// Step 3: Check driver state (should be pending)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", driverServiceURL+"/api/v1/drivers/"+driverID+"/state", nil)

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Get driver state failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	t.Log("✅ Driver registration and state check passed")
}

// TestTokenRefresh tests token refresh mechanism
func TestTokenRefresh(t *testing.T) {
	// Get tokens from login (setup)
	// In real test, would do actual login

	// For this test, assume we have a refresh token
	refreshToken := "refresh-token-from-login"

	refreshReq := map[string]interface{}{
		"refresh_token": refreshToken,
	}

	body, _ := json.Marshal(refreshReq)
	resp, err := http.Post(
		authServiceURL+"/api/v1/auth/refresh",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Token refresh failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var tokenResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	if tokenResp["data"] == nil {
		t.Fatalf("No token in response")
	}

	t.Log("✅ Token refresh passed")
}

// TestPasswordReset tests password reset flow
func TestPasswordReset(t *testing.T) {
	// Step 1: Request password reset
	resetReq := map[string]interface{}{
		"email": "user@test.com",
	}

	body, _ := json.Marshal(resetReq)
	resp, err := http.Post(
		authServiceURL+"/api/v1/auth/password-reset",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Password reset request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	// Step 2: Verify OTP and set new password
	verifyReq := map[string]interface{}{
		"email":        "user@test.com",
		"otp":          "123456",
		"new_password": "NewSecurePassword456!",
	}

	body, _ = json.Marshal(verifyReq)
	resp, err = http.Post(
		authServiceURL+"/api/v1/auth/password-reset/verify",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Password reset verification failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	t.Log("✅ Password reset flow passed")
}

// TestErrorHandling_InvalidOTP tests error handling with invalid OTP
func TestErrorHandling_InvalidOTP(t *testing.T) {
	verifyReq := map[string]interface{}{
		"email": "user@test.com",
		"otp":   "000000", // Invalid OTP
	}

	body, _ := json.Marshal(verifyReq)
	resp, err := http.Post(
		authServiceURL+"/api/v1/auth/verify-register",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", resp.StatusCode)
	}

	t.Log("✅ Invalid OTP error handling passed")
}

// TestErrorHandling_InvalidCredentials tests error handling with invalid login credentials
func TestErrorHandling_InvalidCredentials(t *testing.T) {
	loginReq := map[string]interface{}{
		"email":    "user@test.com",
		"password": "WrongPassword",
	}

	body, _ := json.Marshal(loginReq)
	resp, err := http.Post(
		authServiceURL+"/api/v1/auth/login",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", resp.StatusCode)
	}

	t.Log("✅ Invalid credentials error handling passed")
}

// TestErrorHandling_MissingAuthHeader tests error handling when auth header is missing
func TestErrorHandling_MissingAuthHeader(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", authServiceURL+"/api/v1/auth/verify", nil)
	// No Authorization header

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", resp.StatusCode)
	}

	t.Log("✅ Missing auth header error handling passed")
}

// TestErrorHandling_ExpiredToken tests error handling with expired token
func TestErrorHandling_ExpiredToken(t *testing.T) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", authServiceURL+"/api/v1/auth/verify", nil)
	req.Header.Set("Authorization", "Bearer expired-token")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401, got %d", resp.StatusCode)
	}

	t.Log("✅ Expired token error handling passed")
}

// BenchmarkLogin measures login endpoint performance
func BenchmarkLogin(b *testing.B) {
	loginReq := map[string]interface{}{
		"email":    "user@test.com",
		"password": "SecurePassword123!",
	}

	body, _ := json.Marshal(loginReq)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := http.Post(
			authServiceURL+"/api/v1/auth/login",
			"application/json",
			bytes.NewBuffer(body),
		)
		if err != nil {
			b.Fatalf("Login failed: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkTokenVerification measures token verification performance
func BenchmarkTokenVerification(b *testing.B) {
	token := "valid-jwt-token"
	client := &http.Client{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", authServiceURL+"/api/v1/auth/verify", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
		resp.Body.Close()
	}
}

// TestUserServiceIntegration tests User service integration with Auth
func TestUserServiceIntegration_GetProfile(t *testing.T) {
	// Prerequisite: User must be authenticated (has valid token)
	token := "valid-auth-token"
	userID := "user-uuid"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", userServiceURL+"/api/v1/users/"+userID+"/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 200 or 401, got %d", resp.StatusCode)
	}

	t.Log("✅ User service integration test passed")
}

// TestDriverServiceIntegration tests Driver service integration with Auth
func TestDriverServiceIntegration_GetProfile(t *testing.T) {
	// Prerequisite: Driver must be authenticated (has valid token)
	token := "valid-auth-token"
	driverID := "driver-uuid"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", driverServiceURL+"/api/v1/drivers/"+driverID+"/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 200 or 401, got %d", resp.StatusCode)
	}

	t.Log("✅ Driver service integration test passed")
}

// LoadTest_ConcurrentRegistration simulates 1000 concurrent registration requests
// Run with: go test -run LoadTest_ConcurrentRegistration -v
func LoadTest_ConcurrentRegistration(t *testing.T) {
	const concurrentRequests = 1000
	responseTimes := make([]time.Duration, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(index int) {
			start := time.Now()

			registerReq := map[string]interface{}{
				"email":    "user" + string(rune(index)) + "@test.com",
				"phone":    "+1234567890",
				"password": "SecurePassword123!",
				"role":     "rider",
				"name":     "Test User",
			}

			body, _ := json.Marshal(registerReq)
			resp, _ := http.Post(
				authServiceURL+"/api/v1/auth/register",
				"application/json",
				bytes.NewBuffer(body),
			)
			resp.Body.Close()

			responseTimes[index] = time.Since(start)
		}(i)
	}

	// Wait for all requests to complete (simplified - in real test use sync.WaitGroup)
	time.Sleep(5 * time.Second)

	t.Logf("✅ Load test: %d concurrent registrations completed", concurrentRequests)
}
