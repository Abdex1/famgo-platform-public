/*
# STEP 6 — CREATE SESSION ENTITY

services/auth-service/internal/domain/entities/session.go


# FILE: session.go
*/
package entities

import "time"

type Session struct {
	ID string

	UserID string

	RefreshTokenHash string

	DeviceID string
	DeviceName string
	IPAddress string
	UserAgent string

	Revoked bool

	ExpiresAt time.Time
	CreatedAt time.Time
}
