/*
# PHASE 4 — REDIS DEVICE FINGERPRINT ENGINE

# =========================================================

# STEP 12 — CREATE DEVICE FINGERPRINT STORE

internal/infrastructure/redis/device_store.go
*/
package redis

import (
	"context"
	"encoding/json"
)

type DeviceFingerprint struct {
	DeviceID   string `json:"device_id"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
}

func (s *SessionStore) StoreDevice(
	ctx context.Context,
	userID string,
	device DeviceFingerprint,
) error {

	data, _ := json.Marshal(device)

	return s.client.Set(
		ctx,
		"device:"+userID+":"+device.DeviceID,
		data,
		0,
	).Err()
}
