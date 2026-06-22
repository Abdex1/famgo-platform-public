/*
# REFRESH TOKEN ENTITY

internal/domain/entities/refresh_token.go
*/
package entities

import "time"

type RefreshTokenFamily struct {
    ID          string
    UserID      string
    FamilyID    string
    TokenHash   string

    DeviceID    string
    DeviceName  string
    IPAddress   string
    UserAgent   string

    IsRevoked   bool

    ExpiresAt   time.Time
    CreatedAt   time.Time
}

