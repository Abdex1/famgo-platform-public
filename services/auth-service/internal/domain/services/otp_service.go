/*
# PHASE 7 — OTP VERIFICATION

# CREATE OTP SERVICE

internal/domain/services/otp_service.go
*/
package services

import (
    "fmt"
    "math/rand"
)

func GenerateOTP() string {
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}
