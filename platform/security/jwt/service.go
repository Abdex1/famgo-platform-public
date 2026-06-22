
//platform/security/jwt/service.go

package jwt

import (
    "time"

    gojwt "github.com/golang-jwt/jwt/v5"

    "platform/security/jwt/claims"
    "platform/security/jwt/keys"
)

type Service struct {
    config Config
    keys   *keys.KeyManager
}

func New(config Config, keys *keys.KeyManager) *Service {
    return &Service{
        config: config,
        keys:   keys,
    }
}

func (s *Service) GenerateAccessToken(
    userID string,
    email string,
    role string,
    tenantID string,
) (string, error) {
    now := time.Now()

    c := claims.Claims{
        UserID:    userID,
        Email:     email,
        Role:      role,
        TenantID:  tenantID,
        TokenType: "access",
        RegisteredClaims: gojwt.RegisteredClaims{
            Issuer:    s.config.Issuer,
            Audience:  []string{s.config.Audience},
            Subject:   userID,
            IssuedAt:  gojwt.NewNumericDate(now),
            ExpiresAt: gojwt.NewNumericDate(now.Add(s.config.AccessTokenTTL)),
        },
    }

    token := gojwt.NewWithClaims(gojwt.SigningMethodRS256, c)
    token.Header["kid"] = s.keys.ActiveKID()

    return token.SignedString(s.keys.PrivateKey())
}
