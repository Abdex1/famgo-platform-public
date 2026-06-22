//STEP 3 — CREATE SESSION MODEL
//packages/redis-platform/session/model.go

package session

import "time"

type Session struct {
	SessionID       string    `json:"session_id"`
	UserID          string    `json:"user_id"`
	RefreshTokenID  string    `json:"refresh_token_id"`
	TokenFamilyID   string    `json:"token_family_id"`

	DeviceID        string    `json:"device_id"`
	DeviceName      string    `json:"device_name"`
	OS              string    `json:"os"`
	Browser         string    `json:"browser"`

	IPAddress       string    `json:"ip_address"`
	Country         string    `json:"country"`
	City            string    `json:"city"`

	UserAgent       string    `json:"user_agent"`

	IsRevoked       bool      `json:"is_revoked"`
	IsSuspicious    bool      `json:"is_suspicious"`

	CreatedAt       time.Time `json:"created_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	ExpiresAt       time.Time `json:"expires_at"`
}
