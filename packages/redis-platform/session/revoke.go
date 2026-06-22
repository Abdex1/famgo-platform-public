
// STEP 7 — REVOKE SESSION

//packages/redis-platform/session/revoke.go

package session

import (
	"context"
)

func (s *Store) Revoke(
	ctx context.Context,
	sessionID string,
) error {

	sess, err := s.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	sess.IsRevoked = true

	return s.Save(ctx, *sess)
}
