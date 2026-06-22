
# 4. ENTERPRISE JWT PLATFORM

//platform/security/jwt/config.go

package jwt

import "time"

type Config struct {
    Issuer               string
    Audience             string
    AccessTokenTTL       time.Duration
    RefreshTokenTTL      time.Duration
    EnableKeyRotation    bool
    ActiveKeyID          string
}
