/*
# PHASE 2 — OTP DELIVERY PLATFORM

# =========================================================

# STEP 7 — CREATE OTP STORE

internal/infrastructure/redis/otp_store.go
*/
package redis

import (
	"context"
	"time"
)

func (s *SessionStore) StoreOTP(
	ctx context.Context,
	key string,
	code string,
) error {

	return s.client.Set(
		ctx,
		"otp:"+key,
		code,
		5*time.Minute,
	).Err()
}

func (s *SessionStore) VerifyOTP(
	ctx context.Context,
	key string,
	code string,
) (bool, error) {

	stored, err := s.client.Get(
		ctx,
		"otp:"+key,
	).Result()

	if err != nil {
		return false, err
	}

	return stored == code, nil
}
