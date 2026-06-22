/*
# STEP 7 — CREATE DEVICE FINGERPRINT ENTITY

services/auth-service/internal/domain/entities/device.go


# FILE: device.go
*/
package entities

import "time"

type DeviceFingerprint struct {
	ID string

	UserID string

	DeviceID string

	Platform string
	OS string
	Browser string
	IPAddress string

	LastSeenAt time.Time
	CreatedAt time.Time
}
