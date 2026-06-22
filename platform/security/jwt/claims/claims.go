
//platform/security/jwt/claims/claims.go

package claims

import gojwt "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID      string `json:"user_id"`
    Email       string `json:"email"`
    Role        string `json:"role"`
    TenantID    string `json:"tenant_id"`
    DeviceID    string `json:"device_id"`
    SessionID   string `json:"session_id"`
    TokenType   string `json:"token_type"`

    gojwt.RegisteredClaims
}
