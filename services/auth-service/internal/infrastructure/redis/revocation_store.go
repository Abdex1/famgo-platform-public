/*
# PHASE 1 — SESSION REVOCATION

# =========================================================

# STEP 3 — REDIS REVOCATION STORE

internal/infrastructure/redis/revocation_store.go
*/
package redis

import (
	"context"
	"time"
)

func (s *SessionStore) RevokeSession(
	ctx context.Context,
	sessionID string,
) error {

	return s.client.Del(
		ctx,
		"session:"+sessionID,
	).Err()
}

func (s *SessionStore) BlacklistToken(
	ctx context.Context,
	jti string,
	ttl time.Duration,
) error {

	return s.client.Set(
		ctx,
		"blacklist:"+jti,
		"revoked",
		ttl,
	).Err()
}
