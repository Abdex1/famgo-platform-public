/*
# PHASE 6 — OTP PLATFORM

# STEP 9 — CREATE OTP STORE

packages/redis-platform/otp/store.go
*/
package otp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{rdb: rdb}
}

func otpKey(target string) string {
	return fmt.Sprintf("otp:%s", target)
}

func (s *Store) Save(
	ctx context.Context,
	target string,
	code string,
	ttl time.Duration,
) error {

	return s.rdb.Set(
		ctx,
		otpKey(target),
		code,
		ttl,
	).Err()
}

func (s *Store) Verify(
	ctx context.Context,
	target string,
	code string,
) (bool, error) {

	val, err := s.rdb.Get(
		ctx,
		otpKey(target),
	).Result()

	if err != nil {
		return false, err
	}

	if val != code {
		return false, nil
	}

	_ = s.rdb.Del(ctx, otpKey(target)).Err()

	return true, nil
}
