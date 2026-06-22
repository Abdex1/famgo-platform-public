
// 7. ENTERPRISE VALIDATION PLATFORM

//shared/utilities/validation/validation.go

package validation

import (
    "regexp"
    "strings"
)

var (
    emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) bool {
    email = strings.TrimSpace(email)
    return emailRegex.MatchString(email)
}

func ValidatePassword(password string) bool {
    return len(password) >= 12
}
