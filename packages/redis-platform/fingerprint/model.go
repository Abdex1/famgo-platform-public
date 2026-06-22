
# PHASE 7 — DEVICE FINGERPRINT ENGINE

---
/*
# STEP 10 — CREATE DEVICE MODEL

packages/redis-platform/fingerprint/model.go
*/
package fingerprint

import "time"

type DeviceFingerprint struct {
	UserID       string    `json:"user_id"`
	DeviceID     string    `json:"device_id"`

	IPAddress    string    `json:"ip_address"`

	UserAgent    string    `json:"user_agent"`

	OS           string    `json:"os"`
	Browser      string    `json:"browser"`

	Country      string    `json:"country"`
	City         string    `json:"city"`

	LastSeenAt   time.Time `json:"last_seen_at"`

	RiskScore    float64   `json:"risk_score"`
}
